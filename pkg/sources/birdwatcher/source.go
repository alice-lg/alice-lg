package birdwatcher

import (
	"context"
	"fmt"
	"log"
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

	// Mutexes:
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

		// Notify about missing information
		if config.PeerTableOnly {
			log.Println("WARNING: This bird setup does not use a main table.")
			log.Println("The number of filtered routes an routes not exported can not be determined.")
		}
	}

	return birdwatcher
}

func (b *GenericBirdwatcher) filterProtocols(
	protocols map[string]any,
	protocol string,
) map[string]any {
	response := make(map[string]any)
	response["protocols"] = make(map[string]any)

	for protocolID, protocolData := range protocols {
		if protocolData.(map[string]any)["bird_protocol"] == protocol {
			response["protocols"].(map[string]any)[protocolID] = protocolData
		}
	}

	return response
}

func (b *GenericBirdwatcher) filterProtocolsBgp(
	bird ClientResponse,
) map[string]any {
	return b.filterProtocols(bird["protocols"].(map[string]any), "BGP")
}

func (b *GenericBirdwatcher) filterProtocolsPipe(
	bird ClientResponse,
) map[string]any {
	return b.filterProtocols(bird["protocols"].(map[string]any), "Pipe")
}

func (b *GenericBirdwatcher) filterRoutesByPeerOrLearntFrom(
	routes api.Routes,
	peerPtr *string,
	learntFromPtr *string,
) api.Routes {
	resultRoutes := make(api.Routes, 0, len(routes))

	// Choose routes with next_hop == gateway of this neighbor
	for _, route := range routes {
		if (route.Gateway == peerPtr) ||
			(route.Gateway == learntFromPtr) ||
			(route.LearntFrom == peerPtr) {
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
		routesMap[route.Network] = route
	}

	// Remove routes from "routes" that are contained within filterRoutes
	for _, filterRoute := range filterRoutes {
		delete(routesMap, filterRoute.Network)
	}

	for _, route := range routesMap {
		resultRoutes = append(resultRoutes, route)
	}

	// Sort routes for deterministic ordering
	sort.Sort(resultRoutes)
	routes = resultRoutes // TODO: Check if this even makes sense...

	return routes
}

func (b *GenericBirdwatcher) fetchProtocolsShort(ctx context.Context) (
	*api.Meta,
	map[string]any,
	error,
) {
	// Query birdwatcher with forced timeout
	timeout := 20 * time.Second
	if b.config.NeighborsRefreshTimeout > 0 {
		timeout = time.Duration(b.config.NeighborsRefreshTimeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	bird, err := b.client.GetJSON(ctx, "/protocols/short?uncached=true")
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

// Status retrieves the current backend status
func (b *GenericBirdwatcher) Status(ctx context.Context) (*api.StatusResponse, error) {
	bird, err := b.client.GetJSON(ctx, "/status")
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
func (b *GenericBirdwatcher) NeighborsStatus(ctx context.Context) (
	*api.NeighborsStatusResponse,
	error,
) {
	// Query birdwatcher
	apiStatus, birdProtocols, err := b.fetchProtocolsShort(ctx)
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
