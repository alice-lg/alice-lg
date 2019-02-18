package birdwatcher

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/alice-lg/alice-lg/backend/caches"
	"github.com/alice-lg/alice-lg/backend/sources"

	"sort"
)

type Birdwatcher interface {
	sources.Source
}

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

func NewBirdwatcher(config Config) Birdwatcher {
	client := NewClient(config.Api)

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

func (self *GenericBirdwatcher) filterProtocols(protocols map[string]interface{}, protocol string) map[string]interface{} {
	response := make(map[string]interface{})
	response["protocols"] = make(map[string]interface{})

	for protocolId, protocolData := range protocols {
		if protocolData.(map[string]interface{})["bird_protocol"] == protocol {
			response["protocols"].(map[string]interface{})[protocolId] = protocolData
		}
	}

	return response
}

func (self *GenericBirdwatcher) filterProtocolsBgp(bird ClientResponse) map[string]interface{} {
	return self.filterProtocols(bird["protocols"].(map[string]interface{}), "BGP")
}

func (self *GenericBirdwatcher) filterProtocolsPipe(bird ClientResponse) map[string]interface{} {
	return self.filterProtocols(bird["protocols"].(map[string]interface{}), "Pipe")
}

func (self *GenericBirdwatcher) filterRoutesByPeerOrLearntFrom(routes api.Routes, peer string, learntFrom string) api.Routes {
	result_routes := make(api.Routes, 0, len(routes))

	// Choose routes with next_hop == gateway of this neighbour
	for _, route := range routes {
		if (route.Gateway == peer) ||
			(route.Gateway == learntFrom) ||
			(route.Details["learnt_from"] == peer) {
			result_routes = append(result_routes, route)
		}
	}

	// Sort routes for deterministic ordering
	sort.Sort(result_routes)
	routes = result_routes

	return routes
}

func (self *GenericBirdwatcher) filterRoutesByDuplicates(routes api.Routes, filterRoutes api.Routes) api.Routes {
	result_routes := make(api.Routes, 0, len(routes))

	routesMap := make(map[string]*api.Route) // for O(1) access
	for _, route := range routes {
		routesMap[route.Id] = route
	}

	// Remove routes from "routes" that are contained within filterRoutes
	for _, filterRoute := range filterRoutes {
		if _, ok := routesMap[filterRoute.Id]; ok {
			delete(routesMap, filterRoute.Id)
		}
	}

	for _, route := range routesMap {
		result_routes = append(result_routes, route)
	}

	// Sort routes for deterministic ordering
	sort.Sort(result_routes)
	routes = result_routes

	return routes
}

func (self *GenericBirdwatcher) Status() (*api.StatusResponse, error) {
	// Query birdwatcher
	bird, err := self.client.GetJson("/status")
	if err != nil {
		return nil, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, err
	}

	// Parse the status
	birdStatus, err := parseBirdwatcherStatus(bird, self.config)
	if err != nil {
		return nil, err
	}

	response := &api.StatusResponse{
		Api:    apiStatus,
		Status: birdStatus,
	}

	return response, nil
}

// Make routes lookup
func (self *GenericBirdwatcher) LookupPrefix(prefix string) (*api.RoutesLookupResponse, error) {
	// Get RS info
	rs := api.Routeserver{
		Id:   self.config.Id,
		Name: self.config.Name,
	}

	// Query prefix on RS
	bird, err := self.client.GetJson("/routes/prefix?prefix=" + prefix)
	if err != nil {
		return nil, err
	}

	// Parse API status
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, err
	}

	// Parse routes
	routes, err := parseRoutes(bird, self.config)

	// Add corresponding neighbour and source rs to result
	results := api.LookupRoutes{}
	for _, src := range routes {
		// Okay. This is actually really hacky.
		// A less bruteforce approach would be highly appreciated
		route := &api.LookupRoute{
			Id: src.Id,

			Routeserver: rs,

			NeighbourId: src.NeighbourId,

			Network:   src.Network,
			Interface: src.Interface,
			Gateway:   src.Gateway,
			Metric:    src.Metric,
			Bgp:       src.Bgp,
			Age:       src.Age,
			Type:      src.Type,

			Details: src.Details,
		}
		results = append(results, route)
	}

	// Make result
	response := &api.RoutesLookupResponse{
		Api:    apiStatus,
		Routes: results,
	}
	return response, nil
}
