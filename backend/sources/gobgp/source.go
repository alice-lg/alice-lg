package gobgp

import (
	aliceapi "github.com/alice-lg/alice-lg/backend/api"
	"github.com/alice-lg/alice-lg/backend/caches"
	api "github.com/osrg/gobgp/api"
	"google.golang.org/grpc"

	"log"
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

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

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
	return nil,nil
}

// Get bird BGP protocols
func (gobgp *GoBGP) Neighbours() (*aliceapi.NeighboursResponse, error) {
	return nil,nil
}

// Get neighbors from neighbors summary
func (gobgp *GoBGP) summaryNeighbors() (*aliceapi.NeighboursResponse, error) {
	return nil,nil
}

// Get neighbors from protocols
func (gobgp *GoBGP) bgpProtocolsNeighbors() (*aliceapi.NeighboursResponse, error) {
	return nil,nil
}

// Get filtered and exported routes
func (gobgp *GoBGP) Routes(neighbourId string) (*aliceapi.RoutesResponse, error) {
	return nil,nil
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
	return nil,nil
}


// Get all received routes
func (gobgp *GoBGP) RoutesReceived(neighborId string,) (*aliceapi.RoutesResponse, error) {
	return nil,nil
}


// Get all filtered routes
func (gobgp *GoBGP) RoutesFiltered(neighborId string,) (*aliceapi.RoutesResponse, error) {
	return nil,nil
}

// Get all not exported routes
func (gobgp *GoBGP) RoutesNotExported(neighborId string,) (*aliceapi.RoutesResponse, error) {
	return nil,nil
}

// Make routes lookup
func (gobgp *GoBGP) LookupPrefix(prefix string) (*aliceapi.RoutesLookupResponse, error) {
	return nil,nil
}

func (gobgp *GoBGP) AllRoutes() (*aliceapi.RoutesResponse, error) {
	return nil,nil
}

