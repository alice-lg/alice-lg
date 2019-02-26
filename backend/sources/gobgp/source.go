package gobgp

import (
	aliceapi "github.com/alice-lg/alice-lg/backend/api"
	"github.com/alice-lg/alice-lg/backend/caches"
	api "github.com/osrg/gobgp/api"

	"google.golang.org/grpc"

	"log"
	"fmt"
	"context"
	"time"
	"io"
	_ "sort"
)

type GoBGP struct {
	config Config
	client api.GobgpApiClient

	// Caches: Neighbors
	neighborsCache *caches.NeighborsCache

	// Caches: Routes
	routesRequiredCache    *caches.RoutesCache
	routesReceivedCache    *caches.RoutesCache
	routesFilteredCache    *caches.RoutesCache
	routesNotExportedCache *caches.RoutesCache
}

func NewGoBGP(config Config) *GoBGP {

	dialOpts := make([]grpc.DialOption,0)

	if config.Insecure {
		dialOpts  = append(dialOpts,grpc.WithInsecure())
	} else {
		//TODO: We need credentials...
	}

	conn, err := grpc.Dial(config.Host, dialOpts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := api.NewGobgpApiClient(conn)


	// Cache settings:
	// TODO: Maybe read from config file
	neighborsCacheDisable := false

	routesCacheDisabled := false
	routesCacheMaxSize := 128

	// Initialize caches
	neighborsCache := caches.NewNeighborsCache(neighborsCacheDisable)
	routesRequiredCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)
	routesReceivedCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)
	routesFilteredCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)
	routesNotExportedCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)

	return &GoBGP{
		config: config,
		client: client,

		neighborsCache: neighborsCache,

		routesRequiredCache:    routesRequiredCache,
		routesReceivedCache:    routesReceivedCache,
		routesFilteredCache:    routesFilteredCache,
		routesNotExportedCache: routesNotExportedCache,
	}
}

func (gobgp *GoBGP) Status() (*aliceapi.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()


	resp, err := gobgp.client.GetBgp(ctx, &api.GetBgpRequest{})
	if err != nil {
		return nil, err
	}

	response := aliceapi.StatusResponse{}
	response.Status.RouterId = resp.Global.RouterId
	response.Status.Backend = "gobgp"
	return &response,nil
}

// Get bird BGP protocols
func (gobgp *GoBGP) Neighbours() (*aliceapi.NeighboursResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response := aliceapi.NeighboursResponse{}
	response.Neighbours = make(aliceapi.Neighbours,0)

	resp, err := gobgp.client.ListPeer(ctx, &api.ListPeerRequest{EnableAdvertised: true})
	if err != nil {
		return nil, err
	}
	for {
	    _resp, err := resp.Recv()
	    if err == io.EOF {
	        break
	    }

	    neigh := aliceapi.Neighbour{}

	    neigh.Address = _resp.Peer.State.NeighborAddress
	    neigh.Asn = int(_resp.Peer.State.PeerAs)
	    switch _resp.Peer.State.SessionState {
	    case api.PeerState_ESTABLISHED:
	    	neigh.State = "up"
	    default:
	    	neigh.State = "down"
	    }
	    neigh.Description = _resp.Peer.Conf.Description
	    
	    neigh.Id = PeerHash(_resp.Peer)


	    response.Neighbours = append(response.Neighbours, &neigh)
	    for _, afiSafi := range _resp.Peer.AfiSafis {
	    	neigh.RoutesReceived += int(afiSafi.State.Received)
	    	neigh.RoutesExported += int(afiSafi.State.Advertised)
	    	neigh.RoutesAccepted += int(afiSafi.State.Accepted)
	    	neigh.RoutesFiltered += (neigh.RoutesReceived-neigh.RoutesAccepted)
	    }


		if _resp.Peer.Timers.State.Uptime != nil {
			neigh.Uptime = time.Now().Sub(time.Unix(_resp.Peer.Timers.State.Uptime.Seconds,int64(_resp.Peer.Timers.State.Uptime.Nanos)))
		}

	}

	return &response, nil
}

// Get neighbors from neighbors summary
func (gobgp *GoBGP) summaryNeighbors() (*aliceapi.NeighboursResponse, error) {

	return nil,fmt.Errorf("Not implemented summaryNeighbors")
}

// Get neighbors from protocols
func (gobgp *GoBGP) bgpProtocolsNeighbors() (*aliceapi.NeighboursResponse, error) {
	return nil,fmt.Errorf("Not implemented protocols")
}

// Get filtered and exported routes
func (gobgp *GoBGP) Routes(neighbourId string) (*aliceapi.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse();
	err = gobgp.GetRoutes(neigh,api.TableType_ADJ_IN,&routes)
	if err != nil {
		return nil, err
	}
	return &routes,nil
}

/*
RoutesRequired is a specialized request to fetch:

 - RoutesExported and
 - RoutesFiltered

from Birdwatcher. As the not exported routes can be very many
these are optional and can be loaded on demand using the
RoutesNotExported() API.

A route deduplication is applied.
*/

func (gobgp *GoBGP) getRoutes(neighbourId string,) (*aliceapi.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse();
	err = gobgp.GetRoutes(neigh,api.TableType_ADJ_IN,&routes)
	if err != nil {
		return nil, err
	}
	return &routes,nil
}

func (gobgp *GoBGP) RoutesRequired(neighbourId string,) (*aliceapi.RoutesResponse, error) {
	return gobgp.getRoutes(neighbourId)
}


// Get all received routes
func (gobgp *GoBGP) RoutesReceived(neighbourId string,) (*aliceapi.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}

	routes := NewRoutesResponse();
	err = gobgp.GetRoutes(neigh,api.TableType_ADJ_IN,&routes)
	if err != nil {
		return nil, err
	}
	routes.Filtered = nil
	return &routes,nil
}


// Get all filtered routes
func (gobgp *GoBGP) RoutesFiltered(neighbourId string,) (*aliceapi.RoutesResponse, error) {
	routes, err := gobgp.getRoutes(neighbourId)
	routes.Imported = nil
	return routes, err
}

// Get all not exported routes
func (gobgp *GoBGP) RoutesNotExported(neighbourId string,) (*aliceapi.RoutesResponse, error) {
	neigh, err := gobgp.lookupNeighbour(neighbourId)
	if err != nil {
		return nil, err
	}
	routes := NewRoutesResponse();
	err = gobgp.GetRoutes(neigh,api.TableType_LOCAL,&routes)
	if err != nil {
		return nil, err
	}
	routes.Imported = nil
	return &routes,nil
}

// Make routes lookup
func (gobgp *GoBGP) LookupPrefix(prefix string) (*aliceapi.RoutesLookupResponse, error) {
	return nil,fmt.Errorf("Not implemented LookupPrefix")
}


/*
AllRoutes: 
	Here a routes dump (filtered, received) is returned, which is used to learn all prefixes to build up a local store for searching.
*/
func (gobgp *GoBGP) AllRoutes() (*aliceapi.RoutesResponse, error) {
	routes := NewRoutesResponse();
	peers, err := gobgp.GetNeighbours()
	if err != nil {
		return nil, err
	}

	for _, peer := range peers {
		err = gobgp.GetRoutes(peer,api.TableType_ADJ_IN,&routes)
		if err != nil {
			log.Print(err)
		}
	}

	return &routes, nil
}

