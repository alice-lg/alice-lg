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

	resp, err := gobgp.client.ListPeer(ctx, &api.ListPeerRequest{})
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
	    neigh.State = _resp.Peer.State.SessionState.String()
	    neigh.Description = _resp.Peer.State.Description
	    //neigh.LastError = _resp.Peer.State.Messages.String()

	    response.Neighbours = append(response.Neighbours, &neigh)
	    log.Printf("%+v\n", neigh)
	}

	log.Printf("%+v\n", response)

/*	*/

/*type Neighbour struct {
	Id string `json:"id"`

	// Mandatory fields
	Address         string        `json:"address"`
	Asn             int           `json:"asn"`
	State           string        `json:"state"`
	Description     string        `json:"description"`
	RoutesReceived  int           `json:"routes_received"`
	RoutesFiltered  int           `json:"routes_filtered"`
	RoutesExported  int           `json:"routes_exported"`
	RoutesPreferred int           `json:"routes_preferred"`
	RoutesAccepted  int           `json:"routes_accepted"`
	Uptime          time.Duration `json:"uptime"`
	LastError       string        `json:"last_error"`

	// Original response
	Details map[string]interface{} `json:"details"`
}
*/
	return &response, nil
}

// Get neighbors from neighbors summary
func (gobgp *GoBGP) summaryNeighbors() (*aliceapi.NeighboursResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}

// Get neighbors from protocols
func (gobgp *GoBGP) bgpProtocolsNeighbors() (*aliceapi.NeighboursResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}

// Get filtered and exported routes
func (gobgp *GoBGP) Routes(neighbourId string) (*aliceapi.RoutesResponse, error) {
	return nil,fmt.Errorf("Not implemented")
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

func (gobgp *GoBGP) RoutesRequired(neighborId string,) (*aliceapi.RoutesResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}


// Get all received routes
func (gobgp *GoBGP) RoutesReceived(neighborId string,) (*aliceapi.RoutesResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}


// Get all filtered routes
func (gobgp *GoBGP) RoutesFiltered(neighborId string,) (*aliceapi.RoutesResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}

// Get all not exported routes
func (gobgp *GoBGP) RoutesNotExported(neighborId string,) (*aliceapi.RoutesResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}

// Make routes lookup
func (gobgp *GoBGP) LookupPrefix(prefix string) (*aliceapi.RoutesLookupResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}

func (gobgp *GoBGP) AllRoutes() (*aliceapi.RoutesResponse, error) {
	return nil,fmt.Errorf("Not implemented")
}

