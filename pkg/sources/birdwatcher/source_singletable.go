package birdwatcher

import (
	"log"
	"sort"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
)

// SingleTableBirdwatcher is an Alice Source
type SingleTableBirdwatcher struct {
	GenericBirdwatcher
}

func (src *SingleTableBirdwatcher) fetchReceivedRoutes(
	neighborID string,
) (*api.Meta, api.Routes, error) {
	// Query birdwatcher
	bird, err := src.client.GetJSON("/routes/protocol/" + neighborID)
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, src.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	received, err := parseRoutes(bird, src.config)
	if err != nil {
		log.Println("WARNING Could not retrieve received routes:", err)
		log.Println("Is the 'routes_protocol' module active in birdwatcher?")
		return apiStatus, nil, err
	}

	return apiStatus, received, nil
}

func (src *SingleTableBirdwatcher) fetchFilteredRoutes(
	neighborID string,
) (*api.Meta, api.Routes, error) {
	// Query birdwatcher
	bird, err := src.client.GetJSON("/routes/filtered/" + neighborID)
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, src.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	filtered, err := parseRoutes(bird, src.config)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'routes_filtered' module active in birdwatcher?")
		return apiStatus, nil, err
	}

	return apiStatus, filtered, nil
}

func (src *SingleTableBirdwatcher) fetchNotExportedRoutes(
	neighborID string,
) (*api.Meta, api.Routes, error) {
	// Query birdwatcher
	bird, _ := src.client.GetJSON("/routes/noexport/" + neighborID)

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, src.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	notExported, err := parseRoutes(bird, src.config)
	if err != nil {
		log.Println("WARNING Could not retrieve routes not exported:", err)
		log.Println("Is the 'routes_noexport' module active in birdwatcher?")
	}

	return apiStatus, notExported, nil
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
	apiStatus, receivedRoutes, err := src.fetchReceivedRoutes(neighborID)
	if err != nil {
		return nil, err
	}

	// Second: get routes filtered
	_, filteredRoutes, err := src.fetchFilteredRoutes(neighborID)
	if err != nil {
		return nil, err
	}

	// Perform route deduplication
	importedRoutes := api.Routes{}
	if len(receivedRoutes) > 0 {
		peer := receivedRoutes[0].Gateway
		learntFrom := decoders.String(receivedRoutes[0].Details["learnt_from"], peer)

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
func (src *SingleTableBirdwatcher) Neighbors() (*api.NeighborsResponse, error) {
	// Check if we hit the cache
	response := src.neighborsCache.Get()
	if response != nil {
		return response, nil
	}

	// Query birdwatcher
	bird, err := src.client.GetJSON("/protocols/bgp")
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

// Routes gets filtered and exported routes
func (src *SingleTableBirdwatcher) Routes(
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Fetch required routes first (received and filtered)
	required, err := src.fetchRequiredRoutes(neighborID)
	if err != nil {
		return nil, err
	}

	// Optional: NoExport
	_, notExported, err := src.fetchNotExportedRoutes(neighborID)
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
	routes, err := src.fetchRequiredRoutes(neighborID)
	if err != nil {
		return nil, err
	}

	response.Meta = routes.Meta
	response.Imported = routes.Imported

	return response, nil
}

// RoutesFiltered gets all filtered routes
func (src *SingleTableBirdwatcher) RoutesFiltered(
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
	routes, err := src.fetchRequiredRoutes(neighborID)
	if err != nil {
		return nil, err
	}

	response.Meta = routes.Meta
	response.Filtered = routes.Filtered

	return response, nil
}

// RoutesNotExported get all not exported routes
func (src *SingleTableBirdwatcher) RoutesNotExported(
	neighborID string,
) (*api.RoutesResponse, error) {
	// Check if we hit the cache
	response := src.routesNotExportedCache.Get(neighborID)
	if response != nil {
		return response, nil
	}

	// Fetch not exported routes
	apiStatus, routes, err := src.fetchNotExportedRoutes(neighborID)
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
func (src *SingleTableBirdwatcher) AllRoutes() (*api.RoutesResponse, error) {
	// First fetch all routes from the master table
	mainTable := src.GenericBirdwatcher.config.MainTable
	birdImported, err := src.client.GetJSON("/routes/table/" + mainTable)
	if err != nil {
		return nil, err
	}

	// Then fetch all filtered routes from the master table
	birdFiltered, err := src.client.GetJSON("/routes/table/" + mainTable + "/filtered")
	if err != nil {
		return nil, err
	}

	// Use api status from second request
	apiStatus, err := parseAPIStatus(birdFiltered, src.config)
	if err != nil {
		return nil, err
	}

	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
	}

	// Parse the routes
	imported := parseRoutesData(birdImported["routes"].([]interface{}), src.config)
	// Sort routes for deterministic ordering
	sort.Sort(imported)
	response.Imported = imported

	// Parse the routes
	filtered := parseRoutesData(birdFiltered["routes"].([]interface{}), src.config)
	// Sort routes for deterministic ordering
	sort.Sort(filtered)
	response.Filtered = filtered

	return response, nil
}
