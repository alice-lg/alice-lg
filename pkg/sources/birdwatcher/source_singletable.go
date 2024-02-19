package birdwatcher

import (
	"context"
	"log"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// SingleTableBirdwatcher is an Alice Source
type SingleTableBirdwatcher struct {
	GenericBirdwatcher
}

func (src *SingleTableBirdwatcher) fetchReceivedRoutes(
	ctx context.Context,
	neighborID string,
) (*api.Meta, api.Routes, error) {
	res, err := src.client.GetEndpoint(ctx, "/routes/protocol/"+neighborID)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	meta, routes, err := parseRoutesResponseStream(res.Body, src.config)
	if err != nil {
		return nil, nil, err
	}

	return meta, routes, nil
}

func (src *SingleTableBirdwatcher) fetchFilteredRoutes(
	ctx context.Context,
	neighborID string,
) (*api.Meta, api.Routes, error) {
	res, err := src.client.GetEndpoint(ctx, "/routes/filtered/"+neighborID)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'routes_filtered' module active in birdwatcher?")
		return nil, nil, err
	}
	defer res.Body.Close()

	meta, routes, err := parseRoutesResponseStream(res.Body, src.config)
	if err != nil {
		return nil, nil, err
	}

	return meta, routes, nil
}

func (src *SingleTableBirdwatcher) fetchNotExportedRoutes(
	ctx context.Context,
	neighborID string,
) (*api.Meta, api.Routes, error) {
	res, err := src.client.GetEndpoint(ctx, "/routes/noexport/"+neighborID)
	if err != nil {
		log.Println("WARNING Could not retrieve routes not exported:", err)
		log.Println("Is the 'routes_noexport' module active in birdwatcher?")
		return nil, nil, err
	}
	defer res.Body.Close()

	meta, routes, err := parseRoutesResponseStream(res.Body, src.config)
	if err != nil {
		return nil, nil, err
	}

	return meta, routes, nil
}

// RoutesRequired is a specialized request to fetch:
//
// - RoutesExported and
// - RoutesFiltered
//
// from Birdwatcher. As the not exported routes can be very many
// these are optional and can be loaded on demand using the
// RoutesNotExported() API.
//
// A route deduplication is applied.
func (src *SingleTableBirdwatcher) fetchRequiredRoutes(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	// Allow only one concurrent request for this neighbor
	// to our backend server.
	src.routesFetchMutex.Lock(neighborID)
	defer src.routesFetchMutex.Unlock(neighborID)

	// Check if we have a cache hit
	response := src.routesRequiredCache.Get(neighborID)
	if response != nil {
		return response, nil
	}

	// First: get routes received
	apiStatus, receivedRoutes, err := src.fetchReceivedRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	// Second: get routes filtered
	_, filteredRoutes, err := src.fetchFilteredRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	// Perform route deduplication
	importedRoutes := api.Routes{}
	if len(receivedRoutes) > 0 {
		peer := receivedRoutes[0].Gateway
		learntFrom := receivedRoutes[0].LearntFrom

		filteredRoutes = src.filterRoutesByPeerOrLearntFrom(filteredRoutes, peer, learntFrom)
		importedRoutes = src.filterRoutesByDuplicates(receivedRoutes, filteredRoutes)
	}

	response = &api.RoutesResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		Imported: importedRoutes,
		Filtered: filteredRoutes,
	}

	// Cache result
	src.routesRequiredCache.Set(neighborID, response)

	return response, nil
}

// Neighbors get neighbors from protocols
func (src *SingleTableBirdwatcher) Neighbors(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	// Check if we hit the cache
	response := src.neighborsCache.Get()
	if response != nil {
		return response, nil
	}

	// Query birdwatcher
	bird, err := src.client.GetJSON(ctx, "/protocols/bgp")
	if err != nil {
		return nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, src.config)
	if err != nil {
		return nil, err
	}

	// Parse the neighbors
	neighbors, err := parseNeighbors(bird, src.config)
	if err != nil {
		return nil, err
	}

	neighbors, err = sources.FilterHiddenNeighbors(neighbors, src.config.HiddenNeighbors)
	if err != nil {
		return nil, err
	}

	response = &api.NeighborsResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		Neighbors: neighbors,
	}

	// Cache result
	src.neighborsCache.Set(response)
	return response, nil // dereference for now
}

// NeighborsSummary is for now an alias of Neighbors
func (src *SingleTableBirdwatcher) NeighborsSummary(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	return src.Neighbors(ctx)
}

// Routes gets filtered and exported routes
func (src *SingleTableBirdwatcher) Routes(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Fetch required routes first (received and filtered)
	required, err := src.fetchRequiredRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	// Optional: NoExport
	_, notExported, err := src.fetchNotExportedRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	response.Meta = required.Meta
	response.Imported = required.Imported
	response.Filtered = required.Filtered
	response.NotExported = notExported

	return response, nil
}

// RoutesReceived gets all received routes
func (src *SingleTableBirdwatcher) RoutesReceived(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Check if we hit the cache
	cachedRoutes := src.routesRequiredCache.Get(neighborID)
	if cachedRoutes != nil {
		response.Meta = cachedRoutes.Meta
		response.Imported = cachedRoutes.Imported
		return response, nil
	}

	// Fetch required routes first (received and filtered)
	// However: Store in separate cache for faster access
	routes, err := src.fetchRequiredRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	response.Meta = routes.Meta
	response.Imported = routes.Imported

	return response, nil
}

// RoutesFiltered gets all filtered routes
func (src *SingleTableBirdwatcher) RoutesFiltered(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Check if we hit the cache
	cachedRoutes := src.routesRequiredCache.Get(neighborID)
	if cachedRoutes != nil {
		response.Meta = cachedRoutes.Meta
		response.Filtered = cachedRoutes.Filtered
		return response, nil
	}

	// Fetch required routes first (received and filtered)
	// However: Store in separate cache for faster access
	routes, err := src.fetchRequiredRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	response.Meta = routes.Meta
	response.Filtered = routes.Filtered

	return response, nil
}

// RoutesNotExported get all not exported routes
func (src *SingleTableBirdwatcher) RoutesNotExported(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	// Check if we hit the cache
	response := src.routesNotExportedCache.Get(neighborID)
	if response != nil {
		return response, nil
	}

	// Fetch not exported routes
	apiStatus, routes, err := src.fetchNotExportedRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	response = &api.RoutesResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		NotExported: routes,
	}

	// Cache result
	src.routesNotExportedCache.Set(neighborID, response)

	return response, nil
}

// AllRoutes retrieves a route dump
func (src *SingleTableBirdwatcher) AllRoutes(
	ctx context.Context,
) (*api.RoutesResponse, error) {
	// First fetch all routes from the master table
	mainTable := src.GenericBirdwatcher.config.MainTable

	// Routes received
	res, err := src.client.GetEndpoint(ctx, "/routes/table/"+mainTable)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	meta, birdImported, err := parseRoutesResponseStream(res.Body, src.config)
	if err != nil {
		return nil, err
	}

	// Routes filtered
	res, err = src.client.GetEndpoint(ctx, "/routes/table/"+mainTable+"/filtered")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	_, birdFiltered, err := parseRoutesResponseStream(res.Body, src.config)
	if err != nil {
		return nil, err
	}

	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: meta,
		},
		Imported: birdImported,
		Filtered: birdFiltered,
	}

	return response, nil
}
