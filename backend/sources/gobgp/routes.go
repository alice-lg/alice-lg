package gobgp

import (
	"github.com/alice-lg/alice-lg/backend/sources/gobgp/apiutil"
	"github.com/osrg/gobgp/pkg/packet/bgp"

	"github.com/alice-lg/alice-lg/backend/api"
	gobgpapi "github.com/osrg/gobgp/api"

	"context"
	"fmt"
	"io"
	"log"
	"time"
)

var families []gobgpapi.Family = []gobgpapi.Family{gobgpapi.Family{
	Afi:  gobgpapi.Family_AFI_IP,
	Safi: gobgpapi.Family_SAFI_UNICAST,
}, gobgpapi.Family{
	Afi:  gobgpapi.Family_AFI_IP6,
	Safi: gobgpapi.Family_SAFI_UNICAST,
},
}

func NewRoutesResponse() api.RoutesResponse {
	routes := api.RoutesResponse{}
	routes.Imported = make(api.Routes, 0)
	routes.Filtered = make(api.Routes, 0)
	routes.NotExported = make(api.Routes, 0)
	return routes
}

func (gobgp *GoBGP) lookupNeighbour(neighborId string) (*gobgpapi.Peer, error) {

	peers, err := gobgp.GetNeighbours()
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
		peerId := PeerHash(peer)
		if neighborId == "" || peerId == neighborId {
			return peer, nil
		}
	}

	return nil, fmt.Errorf("Could not lookup neighbour")
}

func (gobgp *GoBGP) GetNeighbours() ([]*gobgpapi.Peer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	peerStream, err := gobgp.client.ListPeer(ctx, &gobgpapi.ListPeerRequest{EnableAdvertised: true})
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
func (gobgp *GoBGP) GetRoutes(peer *gobgpapi.Peer, tableType gobgpapi.TableType, rr *api.RoutesResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, family := range families {

		pathStream, err := gobgp.client.ListPath(ctx, &gobgpapi.ListPathRequest{
			Name:           peer.State.NeighborAddress,
			TableType:      tableType,
			Family:         &family,
			EnableFiltered: true,
		})

		if err != nil {
			return nil
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

		for _, d := range rib {
			for _, path := range d.Paths {
				r := api.Route{}
				r.Id = fmt.Sprintf("%d_%s", path.Identifier, d.Prefix)
				r.NeighbourId = PeerHash(peer)
				r.Network = d.Prefix
				r.Interface = "Unknown"
				r.Age = time.Now().Sub(time.Unix(path.Age.GetSeconds(), int64(path.Age.GetNanos())))
				r.Primary = path.Best

				attrs, _ := apiutil.GetNativePathAttributes(path)

				r.Bgp.Communities = make(api.Communities, 0)
				r.Bgp.LargeCommunities = make(api.Communities, 0)
				r.Bgp.ExtCommunities = make(api.ExtCommunities, 0)
				for _, attr := range attrs {
					switch attr.(type) {
					case *bgp.PathAttributeMultiExitDisc:
						med := attr.(*bgp.PathAttributeMultiExitDisc)
						r.Bgp.Med = int(med.Value)
					case *bgp.PathAttributeNextHop:
						nh := attr.(*bgp.PathAttributeNextHop)
						r.Gateway = nh.Value.String()
						r.Bgp.NextHop = nh.Value.String()
					case *bgp.PathAttributeLocalPref:
						lp := attr.(*bgp.PathAttributeLocalPref)
						r.Bgp.LocalPref = int(lp.Value)
					case *bgp.PathAttributeOrigin:
						origin := attr.(*bgp.PathAttributeOrigin)
						switch origin.Value {
						case bgp.BGP_ORIGIN_ATTR_TYPE_IGP:
							r.Bgp.Origin = "IGP"
						case bgp.BGP_ORIGIN_ATTR_TYPE_EGP:
							r.Bgp.Origin = "EGP"
						case bgp.BGP_ORIGIN_ATTR_TYPE_INCOMPLETE:
							r.Bgp.Origin = "Incomplete"
						}
					case *bgp.PathAttributeAsPath:
						aspath := attr.(*bgp.PathAttributeAsPath)
						for _, aspth := range aspath.Value {
							for _, as := range aspth.GetAS() {
								r.Bgp.AsPath = append(r.Bgp.AsPath, int(as))
							}
						}
					case *bgp.PathAttributeCommunities:
						communities := attr.(*bgp.PathAttributeCommunities)
						for _, community := range communities.Value {
							_community := api.Community{int((0xffff0000 & community) >> 16), int(0xffff & community)}
							r.Bgp.Communities = append(r.Bgp.Communities, _community)
						}

					case *bgp.PathAttributeExtendedCommunities:
						communities := attr.(*bgp.PathAttributeExtendedCommunities)
						for _, community := range communities.Value {
							if _community, ok := community.(*bgp.TwoOctetAsSpecificExtended); ok {
								r.Bgp.ExtCommunities = append(r.Bgp.ExtCommunities, api.ExtCommunity{_community.AS, _community.LocalAdmin})
							}
						}
					case *bgp.PathAttributeLargeCommunities:
						communities := attr.(*bgp.PathAttributeLargeCommunities)
						for _, community := range communities.Values {
							r.Bgp.LargeCommunities = append(r.Bgp.LargeCommunities, api.Community{int(community.ASN), int(community.LocalData1), int(community.LocalData2)})
						}
					}
				}

				r.Metric = (r.Bgp.LocalPref + r.Bgp.Med)
				if path.Filtered {
					rr.Filtered = append(rr.Filtered, &r)
				} else {
					rr.Imported = append(rr.Imported, &r)
				}
			}
		}
	}

	return nil
}
