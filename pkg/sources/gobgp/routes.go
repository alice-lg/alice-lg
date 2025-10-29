package gobgp

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	gobgpapi "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/packet/bgp"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/pools"
	"github.com/alice-lg/alice-lg/pkg/sources/gobgp/apiutil"
)

// NewRoutesResponse creates a new routes response
func NewRoutesResponse() api.RoutesResponse {
	routes := api.RoutesResponse{}
	routes.Imported = make(api.Routes, 0)
	routes.Filtered = make(api.Routes, 0)
	routes.NotExported = make(api.Routes, 0)
	return routes
}

func (gobgp *GoBGP) lookupNeighbor(
	ctx context.Context,
	neighborID string,
) (*gobgpapi.Peer, error) {

	peers, err := gobgp.GetNeighbors(ctx)
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
		peerID := PeerHash(peer)
		if neighborID == "" || peerID == neighborID {
			return peer, nil
		}
	}

	return nil, fmt.Errorf("could not lookup neighbor")
}

// GetNeighbors retrieves all neighbors and returns
// a list of peers.
func (gobgp *GoBGP) GetNeighbors(
	ctx context.Context,
) ([]*gobgpapi.Peer, error) {
	ctx, cancel := context.WithTimeout(
		ctx, time.Second*time.Duration(gobgp.config.ProcessingTimeout))
	defer cancel()

	peerStream, err := gobgp.client.ListPeer(
		ctx, &gobgpapi.ListPeerRequest{EnableAdvertised: true})
	if err != nil {
		return nil, err
	}

	peers := make([]*gobgpapi.Peer, 0)

	for {
		peer, err := peerStream.Recv()
		if err == io.EOF {
			break
		}
		peers = append(peers, peer.Peer)
	}
	return peers, nil
}

func extCommunitySubTypeName(subType bgp.ExtendedCommunityAttrSubType) string {
	switch subType {
	case bgp.EC_SUBTYPE_ROUTE_TARGET:
		return "rt"
	case bgp.EC_SUBTYPE_ROUTE_ORIGIN:
		return "ro"
	default:
		return "generic"
	}
}

func (gobgp *GoBGP) parsePathIntoRoute(
	path *gobgpapi.Path,
	prefix string,
) (*api.Route, error) {

	route := api.Route{}
	// route.ID = fmt.Sprintf("%s_%s", path.SourceId, prefix)
	route.NeighborID = pools.Neighbors.Acquire(
		PeerHashWithASAndAddress(path.SourceAsn, path.NeighborIp))
	route.Network = prefix
	route.Interface = pools.Interfaces.Acquire("unknown")
	route.Age = time.Since(time.Unix(path.Age.GetSeconds(), int64(path.Age.GetNanos())))
	route.Primary = path.Best

	// Set AddrFamily based on prefix
	if strings.Contains(prefix, ":") {
		route.AddrFamily = 2 // IPv6
	} else {
		route.AddrFamily = 1 // IPv4
	}

	attrs, err := apiutil.GetNativePathAttributes(path)
	if err != nil {
		return nil, err
	}

	route.BGP = &api.BGPInfo{}
	route.BGP.Communities = make(api.Communities, 0)
	route.BGP.LargeCommunities = make(api.Communities, 0)
	route.BGP.ExtCommunities = make(api.ExtCommunities, 0)

	for _, attr := range attrs {
		switch attr := attr.(type) {
		case *bgp.PathAttributeMultiExitDisc:
			route.BGP.Med = int(attr.Value)
		case *bgp.PathAttributeNextHop:
			route.Gateway = pools.Gateways4.Acquire(attr.Value.String())
			route.BGP.NextHop = pools.Gateways4.Acquire(attr.Value.String())
		case *bgp.PathAttributeLocalPref:
			route.BGP.LocalPref = int(attr.Value)
		case *bgp.PathAttributeOrigin:
			switch attr.Value {
			case bgp.BGP_ORIGIN_ATTR_TYPE_IGP:
				route.BGP.Origin = pools.Origins.Acquire("IGP")
			case bgp.BGP_ORIGIN_ATTR_TYPE_EGP:
				route.BGP.Origin = pools.Origins.Acquire("EGP")
			case bgp.BGP_ORIGIN_ATTR_TYPE_INCOMPLETE:
				route.BGP.Origin = pools.Origins.Acquire("Incomplete")
			}
		case *bgp.PathAttributeAsPath:
			for _, aspth := range attr.Value {
				for _, as := range aspth.GetAS() {
					route.BGP.AsPath = append(route.BGP.AsPath, int(as))
				}
			}
		case *bgp.PathAttributeCommunities:
			for _, community := range attr.Value {
				apiComm := api.Community{
					int((0xffff0000 & community) >> 16),
					int(0xffff & community)}

				route.BGP.Communities = append(route.BGP.Communities, apiComm)
			}
		case *bgp.PathAttributeMpReachNLRI:
			// We could look at the AFI/SAFI here but gobgp really has
			// already done the work and we can just examine the nexthop length.
			switch len(attr.Nexthop) {
			case 4:
				route.Gateway = pools.Gateways4.Acquire(attr.Nexthop.String())
				route.BGP.NextHop = pools.Gateways4.Acquire(attr.Nexthop.String())
			case 16:
				route.Gateway = pools.Gateways6.Acquire(attr.Nexthop.String())
				route.BGP.NextHop = pools.Gateways6.Acquire(attr.Nexthop.String())
			}
		case *bgp.PathAttributeExtendedCommunities:
			for _, community := range attr.Value {
				if apiComm, ok := community.(*bgp.TwoOctetAsSpecificExtended); ok {
					route.BGP.ExtCommunities = append(
						route.BGP.ExtCommunities,
						api.ExtCommunity{
							extCommunitySubTypeName(apiComm.SubType),
							int(apiComm.AS),
							int(apiComm.LocalAdmin),
						})
				}
			}
		case *bgp.PathAttributeLargeCommunities:
			for _, community := range attr.Values {
				route.BGP.LargeCommunities = append(
					route.BGP.LargeCommunities,
					api.Community{
						int(community.ASN),
						int(community.LocalData1),
						int(community.LocalData2)})
			}
		}
	}

	route.BGP.AsPath = pools.ASPaths.Acquire(route.BGP.AsPath)
	route.BGP.Communities = pools.CommunitiesSets.Acquire(route.BGP.Communities)
	route.BGP.ExtCommunities = pools.ExtCommunitiesSets.Acquire(route.BGP.ExtCommunities)
	route.BGP.LargeCommunities = pools.LargeCommunitiesSets.Acquire(route.BGP.LargeCommunities)

	route.Metric = (route.BGP.LocalPref + route.BGP.Med)

	return &route, nil
}

// GetRoutes retrieves all routes from a peer
// for a table type.
func (gobgp *GoBGP) GetRoutes(
	ctx context.Context,
	peer *gobgpapi.Peer,
	tableType gobgpapi.TableType,
	response *api.RoutesResponse,
) error {
	ctx, cancel := context.WithTimeout(
		ctx, time.Second*time.Duration(gobgp.config.ProcessingTimeout))
	defer cancel()

	for i := 1; i < 3; i++ {

		var family *gobgpapi.Family

		switch i {
		case 1:
			{
				family = &gobgpapi.Family{
					Afi:  gobgpapi.Family_AFI_IP,
					Safi: gobgpapi.Family_SAFI_UNICAST}
			}
		case 2:
			{
				family = &gobgpapi.Family{
					Afi:  gobgpapi.Family_AFI_IP6,
					Safi: gobgpapi.Family_SAFI_UNICAST}
			}
		}

		pathStream, err := gobgp.client.ListPath(ctx, &gobgpapi.ListPathRequest{
			Name:           peer.State.NeighborAddress,
			TableType:      tableType,
			Family:         family,
			EnableFiltered: true,
		})

		if err != nil {
			log.Print(err)
			continue
		}

		rib := make([]*gobgpapi.Destination, 0)
		for {
			_path, err := pathStream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Print(err)
				return err
			}
			rib = append(rib, _path.Destination)
		}

		for _, destination := range rib {
			for _, path := range destination.Paths {
				route, err := gobgp.parsePathIntoRoute(path, destination.Prefix)
				if err != nil {
					log.Println(err)
					continue
				}

				if path.Filtered {
					response.Filtered = append(response.Filtered, route)
				} else {
					response.Imported = append(response.Imported, route)
				}
			}
		}
	}

	return nil
}
