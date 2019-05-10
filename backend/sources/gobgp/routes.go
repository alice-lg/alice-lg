package gobgp;

import (
	"github.com/alice-lg/alice-lg/backend/sources/gobgp/apiutil"
	"github.com/osrg/gobgp/pkg/packet/bgp"

	aliceapi "github.com/alice-lg/alice-lg/backend/api"
	api "github.com/osrg/gobgp/api"

	"log"
	"fmt"
	"context"
	"time"
	"io"
)

var families []api.Family = []api.Family{api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_UNICAST,
	},api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_UNICAST,
	},
}

func NewRoutesResponse() (aliceapi.RoutesResponse) {
	routes := aliceapi.RoutesResponse{}
	routes.Imported = make(aliceapi.Routes,0)
	routes.Filtered = make(aliceapi.Routes,0)
	routes.NotExported = make(aliceapi.Routes,0)
	return routes
}

func generatePeerId(peer *api.Peer) string {
	return PeerHash(peer)
}

func (gobgp *GoBGP) lookupNeighbour(neighborId string) (*api.Peer,error) {

	peers, err := gobgp.GetNeighbours()
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
	    peerId := PeerHash(peer)
	    if neighborId == "" || peerId == neighborId { 
	    	return peer,nil
	    }
	}

	return nil, fmt.Errorf("Could not lookup neighbour")
}


func (gobgp *GoBGP) GetNeighbours() ([]*api.Peer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	peerStream, err := gobgp.client.ListPeer(ctx, &api.ListPeerRequest{EnableAdvertised: true})
	if err != nil {
		return nil,err
	}

	peers := make([]*api.Peer,0)

	for {
	    peer, err := peerStream.Recv()
	    if err == io.EOF {
	        break
	    }
	    peers = append(peers, peer.Peer)
    }
    return peers,nil
}
func (gobgp *GoBGP) GetRoutes(peer *api.Peer, tableType api.TableType, rr *aliceapi.RoutesResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, family := range families {

		pathStream, err := gobgp.client.ListPath(ctx, &api.ListPathRequest{
			Name: peer.State.NeighborAddress,
			TableType: tableType,
			Family: &family,
			EnableFiltered: true,
		})

		if err != nil {
			return nil
		}

		rib := make([]*api.Destination,0)
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
		    	r := aliceapi.Route{}
		    	r.Id = fmt.Sprintf("%d_%s", path.Identifier, d.Prefix)
		    	r.NeighbourId = PeerHash(peer)
		    	r.Network = d.Prefix
		    	r.Interface = "Unknown"
		    	r.Age = time.Now().Sub(time.Unix(path.Age.GetSeconds(),int64(path.Age.GetNanos())))
		    	r.Primary = path.Best

		    	attrs, _ := apiutil.GetNativePathAttributes(path)

		    	r.Bgp.Communities = make(aliceapi.Communities,0)
		    	r.Bgp.LargeCommunities = make(aliceapi.Communities,0)
		    	r.Bgp.ExtCommunities = make(aliceapi.ExtCommunities,0)
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
								_community := aliceapi.Community{int((0xffff0000&community)>>16),int(0xffff&community)}
								r.Bgp.Communities = append(r.Bgp.Communities, _community)
							}

						case *bgp.PathAttributeExtendedCommunities:
							communities := attr.(*bgp.PathAttributeExtendedCommunities)
							for _, community := range communities.Value {
								if _community, ok := community.(*bgp.TwoOctetAsSpecificExtended); ok {
									r.Bgp.ExtCommunities = append(r.Bgp.ExtCommunities, aliceapi.ExtCommunity{_community.AS, _community.LocalAdmin})	
								}
							}
						case *bgp.PathAttributeLargeCommunities:
							communities := attr.(*bgp.PathAttributeLargeCommunities) 
							for _, community := range communities.Values {
								r.Bgp.LargeCommunities = append(r.Bgp.LargeCommunities, aliceapi.Community{int(community.ASN), int(community.LocalData1), int(community.LocalData2)})
							}
		    		}
		    	}

		    	r.Metric = (r.Bgp.LocalPref + r.Bgp.Med)
		    	if path.Filtered {
		    		rr.Filtered = append(rr.Filtered, &r)
		    	}  else {
		    		rr.Imported = append(rr.Imported, &r)
		    	}
		    }
		}
	}
	
	return nil
}