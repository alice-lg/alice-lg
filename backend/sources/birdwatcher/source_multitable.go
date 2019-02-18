package birdwatcher

import (
	"github.com/alice-lg/alice-lg/backend/api"

	"strings"

	"fmt"
	"sort"
	"log"
)


type MultiTableBirdwatcher struct {
	GenericBirdwatcher
}


func (self *MultiTableBirdwatcher) getMasterPipeName(table string) string {
	if strings.HasPrefix(table, self.config.PeerTablePrefix) {
		return self.config.PipeProtocolPrefix + table[1:]
	} else {
		return ""
	}
}

func (self *MultiTableBirdwatcher) parseProtocolToTableTree(bird ClientResponse) map[string]interface{} {
	protocols := bird["protocols"].(map[string]interface{})

	response := make(map[string]interface{})

	for _, protocolData := range protocols {
		protocol := protocolData.(map[string]interface{})

		if protocol["bird_protocol"] == "BGP" {
			table := protocol["table"].(string)
			neighborAddress := protocol["neighbor_address"].(string)

			if _, ok := response[table]; !ok {
				response[table] = make(map[string]interface{})
			}

			if _, ok := response[table].(map[string]interface{})[neighborAddress]; !ok {
				response[table].(map[string]interface{})[neighborAddress] = make(map[string]interface{})
			}

			response[table].(map[string]interface{})[neighborAddress] = protocol
		}
	}

	return response
}


func (self *MultiTableBirdwatcher) fetchProtocols() (*api.ApiStatus, map[string]interface{}, error) {
	// Query birdwatcher
	bird, err := self.client.GetJson("/protocols")
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, nil, err
	}

	if _, ok := bird["protocols"]; !ok {
		return nil, nil, fmt.Errorf("Failed to fetch protocols")
	}

	return &apiStatus, bird, nil
}

func (self *MultiTableBirdwatcher) fetchReceivedRoutes(neighborId string) (*api.ApiStatus, api.Routes, error) {
	// Query birdwatcher
	_, birdProtocols, err := self.fetchProtocols()
	if err != nil {
		return nil, nil, err
	}

	protocols := birdProtocols["protocols"].(map[string]interface{})

	if _, ok := protocols[neighborId]; !ok {
		return nil, nil, fmt.Errorf("Invalid Neighbor")
	}

	peer := protocols[neighborId].(map[string]interface{})["neighbor_address"].(string)

	// Query birdwatcher
	bird, err := self.client.GetJson("/routes/peer/" + peer)
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
		log.Println("Is the 'routes_peer' module active in birdwatcher?")
		return &apiStatus, nil, err
	}

	return &apiStatus, received, nil
}

func (self *MultiTableBirdwatcher) fetchFilteredRoutes(neighborId string) (*api.ApiStatus, api.Routes, error) {
	// Query birdwatcher
	_, birdProtocols, err := self.fetchProtocols()
	if err != nil {
		return nil, nil, err
	}

	protocols := birdProtocols["protocols"].(map[string]interface{})

	if _, ok := protocols[neighborId]; !ok {
		return nil, nil, fmt.Errorf("Invalid Neighbor")
	}

	// Stage 1 filters
	birdFiltered, err := self.client.GetJson("/routes/filtered/" + neighborId)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'routes_filtered' module active in birdwatcher?")
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(birdFiltered, self.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	filtered := parseRoutesData(birdFiltered["routes"].([]interface{}), self.config)

	// Stage 2 filters
	table := protocols[neighborId].(map[string]interface{})["table"].(string)
	pipeName := self.getMasterPipeName(table)

	// If there is no pipe to master, there is nothing left to do
	if pipeName == "" {
		return &apiStatus, filtered, nil
	}

	// Query birdwatcher
	birdPipeFiltered, err := self.client.GetJson("/routes/pipe/filtered/?table=" + table + "&pipe=" + pipeName)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'pipe_filtered' module active in birdwatcher?")
		return &apiStatus, nil, err
	}

	// Parse the routes
	pipeFiltered := parseRoutesData(birdPipeFiltered["routes"].([]interface{}), self.config)

	// Sort routes for deterministic ordering
	filtered = append(filtered, pipeFiltered...)
	sort.Sort(filtered)

	return &apiStatus, filtered, nil
}

func (self *MultiTableBirdwatcher) fetchNotExportedRoutes(neighborId string) (*api.ApiStatus, api.Routes, error) {
	// Query birdwatcher
	_, birdProtocols, err := self.fetchProtocols()
	if err != nil {
		return nil, nil, err
	}

	protocols := birdProtocols["protocols"].(map[string]interface{})

	if _, ok := protocols[neighborId]; !ok {
		return nil, nil, fmt.Errorf("Invalid Neighbor")
	}

	table := protocols[neighborId].(map[string]interface{})["table"].(string)
	pipeName := self.getMasterPipeName(table)

	// Query birdwatcher
	bird, err := self.client.GetJson("/routes/noexport/" + pipeName)

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return nil, nil, err
	}

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
func (self *MultiTableBirdwatcher) fetchRequiredRoutes(neighborId string) (*api.RoutesResponse, error) {
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
func (self *MultiTableBirdwatcher) Neighbours() (*api.NeighboursResponse, error) {
	// Query birdwatcher
	apiStatus, birdProtocols, err := self.fetchProtocols()
	if err != nil {
		return nil, err
	}

	// Parse the neighbors
	neighbours, err := parseNeighbours(self.filterProtocolsBgp(birdProtocols), self.config)
	if err != nil {
		return nil, err
	}

	pipes := self.filterProtocolsPipe(birdProtocols)["protocols"].(map[string]interface{})
	tree := self.parseProtocolToTableTree(birdProtocols)

	// Now determine the session count for each neighbor and check if the pipe
	// did filter anything
	filtered := make(map[string]int)
	for table, _ := range tree {
		allRoutesImported := int64(0)
		pipeRoutesImported := int64(0)

		// Sum up all routes from all peers for a table
		for _, protocol := range tree[table].(map[string]interface{}) {
			// Skip peers that are not up (start/down)
			if protocol.(map[string]interface{})["state"].(string) != "up" {
				continue
			}
			allRoutesImported += int64(protocol.(map[string]interface{})["routes"].(map[string]interface{})["imported"].(float64))

			pipeName := self.getMasterPipeName(table)

			if _, ok := pipes[pipeName]; ok {
				if _, ok := pipes[pipeName].(map[string]interface{})["routes"].(map[string]interface{})["imported"]; ok {
					pipeRoutesImported = int64(pipes[pipeName].(map[string]interface{})["routes"].(map[string]interface{})["imported"].(float64))
				} else {
					continue
				}
			} else {
				continue
			}
		}

		// If no routes were imported, there is nothing left to filter
		if allRoutesImported == 0 {
			continue
		}

		// If the pipe did not filter anything, there is nothing left to do
		if pipeRoutesImported == allRoutesImported {
			continue
		}

		if len(tree[table].(map[string]interface{})) == 1 {
			// Single router
			for _, protocol := range tree[table].(map[string]interface{}) {
				filtered[protocol.(map[string]interface{})["protocol"].(string)] = int(allRoutesImported-pipeRoutesImported)
			}
		} else {
			// Multiple routers
			if pipeRoutesImported == 0 {
				// 0 is a special condition, which means that the pipe did filter ALL routes of
				// all peers. Therefore we already know the amount of filtered routes and don't have
				// to query birdwatcher again.
				for _, protocol := range tree[table].(map[string]interface{}) {
					// Skip peers that are not up (start/down)
					if protocol.(map[string]interface{})["state"].(string) != "up" {
						continue
					}
					filtered[protocol.(map[string]interface{})["protocol"].(string)] = int(protocol.(map[string]interface{})["routes"].(map[string]interface{})["imported"].(float64))
				}
			} else {
				// Otherwise the pipe did import at least some routes which means that
				// we have to query birdwatcher to get the count for each peer.
				for neighborAddress, protocol := range tree[table].(map[string]interface{}) {
					table := protocol.(map[string]interface{})["table"].(string)
					pipe := self.getMasterPipeName(table)

					count, err := self.client.GetJson("/routes/pipe/filtered/count?table=" + table + "&pipe=" + pipe + "&address=" + neighborAddress)
					if err != nil {
						log.Println("WARNING Could not retrieve filtered routes count:", err)
						log.Println("Is the 'pipe_filtered_count' module active in birdwatcher?")
						return nil, err
					}

					if _, ok := count["routes"]; ok {
						filtered[protocol.(map[string]interface{})["protocol"].(string)] = int(count["routes"].(float64))
					}
				}
			}
		}
	}

	// Update the results with the information about filtered routes from the pipe
	for _, neighbor := range neighbours {
		if pipeRoutesFiltered, ok := filtered[neighbor.Id]; ok {
			neighbor.RoutesAccepted -= pipeRoutesFiltered
			neighbor.RoutesFiltered += pipeRoutesFiltered
		}
	}

	response := &api.NeighboursResponse{
		Api:        *apiStatus,
		Neighbours: neighbours,
	}

	return response, nil // dereference for now
}

// Get filtered and exported routes
func (self *MultiTableBirdwatcher) Routes(neighbourId string) (*api.RoutesResponse, error) {
	// Fetch required routes first (received and filtered)
	// However: Store in separate cache for faster access
	response, err := self.fetchRequiredRoutes(neighbourId)
	if err != nil {
		return nil, err
	}

	// Optional: NoExport
	_, notExported, err := self.fetchNotExportedRoutes(neighbourId)
	if err != nil {
		return nil, err
	}

	response.NotExported = notExported

	return response, nil
}

// Get all received routes
func (self *MultiTableBirdwatcher) RoutesReceived(neighborId string) (*api.RoutesResponse, error) {
	// Check if we have a cache hit
	response := self.routesRequiredCache.Get(neighborId)
	if response != nil {
		response.Filtered = nil
		return response, nil
	}

	// Fetch required routes first (received and filtered)
	routes, err := self.fetchRequiredRoutes(neighborId)
	if err != nil {
		return nil, err
	}

	response = &api.RoutesResponse{
		Api:      routes.Api,
		Imported: routes.Imported,
	}

	return response, nil
}

// Get all filtered routes
func (self *MultiTableBirdwatcher) RoutesFiltered(neighborId string) (*api.RoutesResponse, error) {
	// Check if we have a cache hit
	response := self.routesRequiredCache.Get(neighborId)
	if response != nil {
		response.Imported = nil
		return response, nil
	}

	// Fetch required routes first (received and filtered)
	routes, err := self.fetchRequiredRoutes(neighborId)
	if err != nil {
		return nil, err
	}

	response = &api.RoutesResponse{
		Api:      routes.Api,
		Filtered: routes.Filtered,
	}

	return response, nil
}

// Get all not exported routes
func (self *MultiTableBirdwatcher) RoutesNotExported(neighborId string) (*api.RoutesResponse, error) {
	// Check if we have a cache hit
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

func (self *MultiTableBirdwatcher) AllRoutes() (*api.RoutesResponse, error) {
	// Query birdwatcher
	_, birdProtocols, err := self.fetchProtocols()
	if err != nil {
		return nil, err
	}

	// Fetch received routes first
	birdImported, err := self.client.GetJson("/routes/table/master")
	if err != nil {
		return nil, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(birdImported, self.config)
	if err != nil {
		return nil, err
	}

	response := &api.RoutesResponse{
		Api:    apiStatus,
	}

	// Parse the routes
	imported := parseRoutesData(birdImported["routes"].([]interface{}), self.config)
	// Sort routes for deterministic ordering
	sort.Sort(imported)
	response.Imported = imported

	// Iterate over all the protocols and fetch the filtered routes for everyone
	protocolsBgp := self.filterProtocolsBgp(birdProtocols)
	for protocolId, protocolsData := range protocolsBgp["protocols"].(map[string]interface{}) {
		peer := protocolsData.(map[string]interface{})["neighbor_address"].(string)
		learntFrom := mustString(protocolsData.(map[string]interface{})["learnt_from"], peer)

		// Fetch filtered routes
		_, filtered, err := self.fetchFilteredRoutes(protocolId)
		if err != nil {
			continue
		}

		// Perform route deduplication
		filtered = self.filterRoutesByPeerOrLearntFrom(filtered, peer, learntFrom)
		response.Filtered = append(response.Filtered, filtered...)
	}

	return response, nil
}
