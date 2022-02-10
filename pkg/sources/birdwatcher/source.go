package birdwatcher

import (
	"fmt"
	"sort"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/caches"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// A Birdwatcher source is a variant of an alice
// source and implements different strategies for fetching
// route information from bird.
type Birdwatcher interface {
	sources.Source
}

// GenericBirdwatcher is an Alice data source.
type GenericBirdwatcher struct {
	config Config
	client *Client

	// Caches: Neighbors
	neighborsCache *caches.NeighborsCache

	// Caches: Routes
	routesRequiredCache    *caches.RoutesCache
	routesNotExportedCache *caches.RoutesCache

	// Mutices:
	routesFetchMutex *LockMap
}

// NewBirdwatcher creates a new Birdwatcher instance.
// This might be either a GenericBirdWatcher or a MultiTableBirdwatcher.
func NewBirdwatcher(config Config) Birdwatcher {
	client := NewClient(config.API)

	// Cache settings:
	// TODO: Maybe read from config file
	neighborsCacheDisable := false

	routesCacheDisabled := false
	routesCacheMaxSize := 128

	// Initialize caches
	neighborsCache := caches.NewNeighborsCache(neighborsCacheDisable)
	routesRequiredCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)
	routesNotExportedCache := caches.NewRoutesCache(
		routesCacheDisabled, routesCacheMaxSize)

	var birdwatcher Birdwatcher

	if config.Type == "single_table" {
		singleTableBirdwatcher := new(SingleTableBirdwatcher)

		singleTableBirdwatcher.config = config
		singleTableBirdwatcher.client = client

		singleTableBirdwatcher.neighborsCache = neighborsCache

		singleTableBirdwatcher.routesRequiredCache = routesRequiredCache
		singleTableBirdwatcher.routesNotExportedCache = routesNotExportedCache

		singleTableBirdwatcher.routesFetchMutex = NewLockMap()

		birdwatcher = singleTableBirdwatcher
	} else if config.Type == "multi_table" {
		multiTableBirdwatcher := new(MultiTableBirdwatcher)

		multiTableBirdwatcher.config = config
		multiTableBirdwatcher.client = client

		multiTableBirdwatcher.neighborsCache = neighborsCache

		multiTableBirdwatcher.routesRequiredCache = routesRequiredCache
		multiTableBirdwatcher.routesNotExportedCache = routesNotExportedCache

		multiTableBirdwatcher.routesFetchMutex = NewLockMap()

		birdwatcher = multiTableBirdwatcher
	}

	return birdwatcher
}

func (b *GenericBirdwatcher) filterProtocols(
	protocols map[string]interface{},
	protocol string,
) map[string]interface{} {
	response := make(map[string]interface{})
	response["protocols"] = make(map[string]interface{})

	for protocolID, protocolData := range protocols {
		if protocolData.(map[string]interface{})["bird_protocol"] == protocol {
			response["protocols"].(map[string]interface{})[protocolID] = protocolData
		}
	}

	return response
}

func (b *GenericBirdwatcher) filterProtocolsBgp(
	bird ClientResponse,
) map[string]interface{} {
	return b.filterProtocols(bird["protocols"].(map[string]interface{}), "BGP")
}

func (b *GenericBirdwatcher) filterProtocolsPipe(
	bird ClientResponse,
) map[string]interface{} {
	return b.filterProtocols(bird["protocols"].(map[string]interface{}), "Pipe")
}

func (b *GenericBirdwatcher) filterRoutesByPeerOrLearntFrom(
	routes api.Routes,
	peer string,
	learntFrom string,
) api.Routes {
	resultRoutes := make(api.Routes, 0, len(routes))

	// Choose routes with next_hop == gateway of this neighbor
	for _, route := range routes {
		if (route.Gateway == peer) ||
			(route.Gateway == learntFrom) ||
			(route.LearntFrom == peer) {
			resultRoutes = append(resultRoutes, route)
		}
	}

	// Sort routes for deterministic ordering
	sort.Sort(resultRoutes)
	routes = resultRoutes

	return routes
}

func (b *GenericBirdwatcher) filterRoutesByDuplicates(
	routes api.Routes,
	filterRoutes api.Routes,
) api.Routes {
	resultRoutes := make(api.Routes, 0, len(routes))

	routesMap := make(map[string]*api.Route) // for O(1) access
	for _, route := range routes {
		routesMap[route.ID] = route
	}

	// Remove routes from "routes" that are contained within filterRoutes
	for _, filterRoute := range filterRoutes {
		delete(routesMap, filterRoute.ID)
		// in theorey this guard is unneccessary
		//if _, ok := routesMap[filterRoute.ID]; ok {
		// }
	}

	for _, route := range routesMap {
		resultRoutes = append(resultRoutes, route)
	}

	// Sort routes for deterministic ordering
	sort.Sort(resultRoutes)
	routes = resultRoutes // TODO: Check if this even makes sense...

	return routes
}

func (b *GenericBirdwatcher) fetchProtocolsShort() (
	*api.Meta,
	map[string]interface{},
	error,
) {
	// Query birdwatcher
	timeout := 2 * time.Second
	if b.config.NeighborsRefreshTimeout > 0 {
		timeout = time.Duration(b.config.NeighborsRefreshTimeout) * time.Second
	}
	bird, err := b.client.GetJSONTimeout(timeout, "/protocols/short?uncached=true")
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, b.config)
	if err != nil {
		return nil, nil, err
	}

	if _, ok := bird["protocols"]; !ok {
		return nil, nil, fmt.Errorf("failed to fetch protocols")
	}

	return apiStatus, bird, nil
}

// ExpireCaches clears all local caches
func (b *GenericBirdwatcher) ExpireCaches() int {
	count := b.routesRequiredCache.Expire()
	count += b.routesNotExportedCache.Expire()
	return count
}

// Status retrievs the current backend status
func (b *GenericBirdwatcher) Status() (*api.StatusResponse, error) {
	bird, err := b.client.GetJSON("/status")
	if err != nil {
		return nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, b.config)
	if err != nil {
		return nil, err
	}

	// Parse the status
	birdStatus, err := parseBirdwatcherStatus(bird, b.config)
	if err != nil {
		return nil, err
	}

	response := &api.StatusResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		Status: birdStatus,
	}

	return response, nil
}

// NeighborsStatus retrieves neighbor status infos
func (b *GenericBirdwatcher) NeighborsStatus() (
	*api.NeighborsStatusResponse,
	error,
) {
	// Query birdwatcher
	apiStatus, birdProtocols, err := b.fetchProtocolsShort()
	if err != nil {
		return nil, err
	}

	// Parse the neighbors short
	neighbors, err := parseNeighborsShort(birdProtocols, b.config)
	if err != nil {
		return nil, err
	}

	response := &api.NeighborsStatusResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		Neighbors: neighbors,
	}
	return response, nil // dereference for now
}

// LookupPrefix makes a routes lookup
func (b *GenericBirdwatcher) LookupPrefix(
	prefix string,
) (*api.RoutesLookupResponse, error) {
	// Get RS info
	rs := &api.RouteServer{
		ID:   b.config.ID,
		Name: b.config.Name,
	}

	// Query prefix on RS
	bird, err := b.client.GetJSON("/routes/prefix?prefix=" + prefix)
	if err != nil {
		return nil, err
	}

	// Parse API status
	apiStatus, err := parseAPIStatus(bird, b.config)
	if err != nil {
		return nil, err
	}

	// Parse routes
	routes, _ := parseRoutes(bird, b.config)

	// Add corresponding neighbor and source rs to result
	results := api.LookupRoutes{}
	for _, src := range routes {
		route := &api.LookupRoute{
			RouteServer: rs,
			Route:       src,
		}
		results = append(results, route)
	}

	// Make result
	response := &api.RoutesLookupResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		Routes: results,
	}
	return response, nil
}
