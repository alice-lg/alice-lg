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

func (gobgp *GoBGP) GetRoutes(neighborId string) (*aliceapi.RoutesResponse) {
	rr := aliceapi.RoutesResponse{}
	rr.Imported = make(aliceapi.Routes,0)
	rr.Filtered = make(aliceapi.Routes,0)
	rr.NotExported = make(aliceapi.Routes,0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()


	// Go over all peers

	peerStream, err := gobgp.client.ListPeer(ctx, &api.ListPeerRequest{EnableAdvertised: true})
	if err != nil {
		return nil
	}
	for {
	    peer, err := peerStream.Recv()
	    if err == io.EOF {
	        break
	    }

	    peerId := fmt.Sprintf("%d_%s",peer.Peer.State.PeerAs, peer.Peer.State.NeighborAddress)
	    if neighborId == "" || peerId == neighborId { 

	    	for _, family := range families {

				pathStream, err := gobgp.client.ListPath(ctx, &api.ListPathRequest{
					Name: peer.Peer.State.NeighborAddress,
					TableType: api.TableType_ADJ_IN,
					Family: &family,
				})

				if err != nil {
					return nil
				}

				for {
				    _path, err := pathStream.Recv()
				    if err == io.EOF {
				        break
				    }
				    

				    for _, path := range _path.Destination.Paths {
				    	r := aliceapi.Route{}
				    	r.Id = fmt.Sprintf("%d_%s", path.Identifier, _path.Destination.Prefix)
				    	r.NeighbourId = peerId
				    	r.Network = _path.Destination.Prefix
				    	r.Interface = "Unknown"
				    	r.Age = time.Now().Sub(time.Unix(path.Age.GetSeconds(),int64(path.Age.GetNanos())))
				    	r.Primary = path.Best

				    	attrs, _ := apiutil.GetNativePathAttributes(path)

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
										log.Printf("%+s\n", community)	
									}
								case *bgp.PathAttributeLargeCommunities:
									communities := attr.(*bgp.PathAttributeLargeCommunities) 
									for _, community := range communities.Values {
										log.Printf("%+s\n", community)
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
		}
	}
	log.Printf("%+v", rr)
	return &rr
}