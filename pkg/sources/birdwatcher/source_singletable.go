package birdwatcher

import (
	"log"
	"sort"

	"github.com/alice-lg/alice-lg/pkg/api"
)

type SingleTableBirdwatcher struct {
	GenericBirdwatcher
}

func (self *SingleTableBirdwatcher) fetchReceivedRoutes(neighborId string) (*api.ApiStatus, api.Routes, error) {
	// Query birdwatcher
	bird, err := self.client.GetJson("/routes/protocol/" + neighborId)
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	received, err := parseRoutes(bird, self.config)
	if err != nil {
		log.Println("WARNING Could not retrieve received routes:", err)
		log.Println("Is the 'routes_protocol' module active in birdwatcher?")
		return &apiStatus, nil, err
	}

	return &apiStatus, received, nil
}

func (self *SingleTableBirdwatcher) fetchFilteredRoutes(neighborId string) (*api.ApiStatus, api.Routes, error) {
	// Query birdwatcher
	bird, err := self.client.GetJson("/routes/filtered/" + neighborId)
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	filtered, err := parseRoutes(bird, self.config)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'routes_filtered' module active in birdwatcher?")
		return &apiStatus, nil, err
	}

	return &apiStatus, filtered, nil
}

func (self *SingleTableBirdwatcher) fetchNotExportedRoutes(neighborId string) (*api.ApiStatus, api.Routes, error) {
	// Query birdwatcher
	bird, err := self.client.GetJson("/routes/noexport/" + neighborId)

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	notExported, err := parseRoutes(bird, self.config)
	if err != nil {
		log.Println("WARNING Could not retrieve routes not exported:", err)
		log.Println("Is the 'routes_noexport' module active in birdwatcher?")
	}

	return &apiStatus, notExported, nil
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
func (self *SingleTableBirdwatcher) fetchRequiredRoutes(neighborId string) (*api.RoutesResponse, error) {
	// Allow only one concurrent request for this neighbor
	// to our backend server.
	self.routesFetchMutex.Lock(neighborId)
	defer self.routesFetchMutex.Unlock(neighborId)

	// Check if we have a cache hit
	response := self.routesRequiredCache.Get(neighborId)
	if response != nil {
		return response, nil
	}

	// First: get routes received
	apiStatus, receivedRoutes, err := self.fetchReceivedRoutes(neighborId)
	if err != nil {
		return nil, err
	}

	// Second: get routes filtered
	_, filteredRoutes, err := self.fetchFilteredRoutes(neighborId)
	if err != nil {
		return nil, err
	}

	// Perform route deduplication
	importedRoutes := api.Routes{}
	if len(receivedRoutes) > 0 {
		peer := receivedRoutes[0].Gateway
		learntFrom := mustString(receivedRoutes[0].Details["learnt_from"], peer)

		filteredRoutes = self.filterRoutesByPeerOrLearntFrom(filteredRoutes, peer, learntFrom)
		importedRoutes = self.filterRoutesByDuplicates(receivedRoutes, filteredRoutes)
	}

	response = &api.RoutesResponse{
		Api:      *apiStatus,
		Imported: importedRoutes,
		Filtered: filteredRoutes,
	}

	// Cache result
	self.routesRequiredCache.Set(neighborId, response)

	return response, nil
}

// Get neighbors from protocols
func (self *SingleTableBirdwatcher) Neighbours() (*api.NeighboursResponse, error) {
	// Check if we hit the cache
	response := self.neighborsCache.Get()
	if response != nil {
		return response, nil
	}

	// Query birdwatcher
	bird, err := self.client.GetJson("/protocols/bgp")
	if err != nil {
		return nil, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, err
	}

	// Parse the neighbors
	neighbours, err := parseNeighbours(bird, self.config)
	if err != nil {
		return nil, err
	}

	response = &api.NeighboursResponse{
		Api:        apiStatus,
		Neighbours: neighbours,
	}

	// Cache result
	self.neighborsCache.Set(response)

	return response, nil // dereference for now
}

// Get filtered and exported routes
func (self *SingleTableBirdwatcher) Routes(neighbourId string) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Fetch required routes first (received and filtered)
	required, err := self.fetchRequiredRoutes(neighbourId)
	if err != nil {
		return nil, err
	}

	// Optional: NoExport
	_, notExported, err := self.fetchNotExportedRoutes(neighbourId)
	if err != nil {
		return nil, err
	}

	response.Api = required.Api
	response.Imported = required.Imported
	response.Filtered = required.Filtered
	response.NotExported = notExported

	return response, nil
}

// Get all received routes
func (self *SingleTableBirdwatcher) RoutesReceived(neighborId string) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Check if we hit the cache
	cachedRoutes := self.routesRequiredCache.Get(neighborId)
	if cachedRoutes != nil {
		response.Api = cachedRoutes.Api
		response.Imported = cachedRoutes.Imported
		return response, nil
	}

	// Fetch required routes first (received and filtered)
	// However: Store in separate cache for faster access
	routes, err := self.fetchRequiredRoutes(neighborId)
	if err != nil {
		return nil, err
	}

	response.Api = routes.Api
	response.Imported = routes.Imported

	return response, nil
}

// Get all filtered routes
func (self *SingleTableBirdwatcher) RoutesFiltered(neighborId string) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Check if we hit the cache
	cachedRoutes := self.routesRequiredCache.Get(neighborId)
	if cachedRoutes != nil {
		response.Api = cachedRoutes.Api
		response.Filtered = cachedRoutes.Filtered
		return response, nil
	}

	// Fetch required routes first (received and filtered)
	// However: Store in separate cache for faster access
	routes, err := self.fetchRequiredRoutes(neighborId)
	if err != nil {
		return nil, err
	}

	response.Api = routes.Api
	response.Filtered = routes.Filtered

	return response, nil
}

// Get all not exported routes
func (self *SingleTableBirdwatcher) RoutesNotExported(neighborId string) (*api.RoutesResponse, error) {
	// Check if we hit the cache
	response := self.routesNotExportedCache.Get(neighborId)
	if response != nil {
		return response, nil
	}

	// Fetch not exported routes
	apiStatus, routes, err := self.fetchNotExportedRoutes(neighborId)
	if err != nil {
		return nil, err
	}

	response = &api.RoutesResponse{
		Api:         *apiStatus,
		NotExported: routes,
	}

	// Cache result
	self.routesNotExportedCache.Set(neighborId, response)

	return response, nil
}

func (self *SingleTableBirdwatcher) AllRoutes() (*api.RoutesResponse, error) {
	// First fetch all routes from the master table
	mainTable := self.GenericBirdwatcher.config.MainTable
	birdImported, err := self.client.GetJson("/routes/table/" + mainTable)
	if err != nil {
		return nil, err
	}

	// Then fetch all filtered routes from the master table
	birdFiltered, err := self.client.GetJson("/routes/table/" + mainTable + "/filtered")
	if err != nil {
		return nil, err
	}

	// Use api status from second request
	apiStatus, err := parseApiStatus(birdFiltered, self.config)
	if err != nil {
		return nil, err
	}

	response := &api.RoutesResponse{
		Api: apiStatus,
	}

	// Parse the routes
	imported := parseRoutesData(birdImported["routes"].([]interface{}), self.config)
	// Sort routes for deterministic ordering
	sort.Sort(imported)
	response.Imported = imported

	// Parse the routes
	filtered := parseRoutesData(birdFiltered["routes"].([]interface{}), self.config)
	// Sort routes for deterministic ordering
	sort.Sort(filtered)
	response.Filtered = filtered

	return response, nil
}
