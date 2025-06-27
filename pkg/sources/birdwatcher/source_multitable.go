package birdwatcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
	"github.com/alice-lg/alice-lg/pkg/pools"
)

// MultiTableBirdwatcher implements a birdwatcher with
// a multitable bird as a datasource.
type MultiTableBirdwatcher struct {
	GenericBirdwatcher
}

func (src *MultiTableBirdwatcher) getMasterPipeName(table string) string {
	ptPrefix := src.config.PeerTablePrefix
	if strings.HasPrefix(table, ptPrefix) {
		return src.config.PipeProtocolPrefix + table[len(ptPrefix):]
	}
	return ""
}

// isAltSession checks if the pipe ends in a
// known suffix, e.g. "_lg". If the alt_pipe_suffix is
// not configured, this will always be false.
func (src *MultiTableBirdwatcher) isAltSession(pipe string) bool {
	suffix := src.config.AltPipeProtocolSuffix
	if suffix == "" {
		return false
	}
	return strings.HasSuffix(pipe, suffix)
}

func (src *MultiTableBirdwatcher) getAltPipeName(pipe string) string {
	prefix := src.config.PipeProtocolPrefix
	return src.config.AltPipeProtocolPrefix + pipe[len(prefix):]
}

func (src *MultiTableBirdwatcher) parseProtocolToTableTree(
	bird ClientResponse,
) map[string]interface{} {
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
				response[table].(map[string]interface{})[neighborAddress] = make(
					map[string]interface{})
			}

			response[table].(map[string]interface{})[neighborAddress] = protocol
		}
	}

	return response
}

func (src *MultiTableBirdwatcher) fetchProtocols(ctx context.Context) (
	*api.Meta,
	map[string]interface{},
	error,
) {
	// Query birdwatcher
	bird, err := src.client.GetJSON(ctx, "/protocols")
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, src.config)
	if err != nil {
		return nil, nil, err
	}

	if _, ok := bird["protocols"]; !ok {
		return nil, nil, fmt.Errorf("failed to fetch protocols")
	}

	return apiStatus, bird, nil
}

func (src *MultiTableBirdwatcher) fetchReceivedRoutes(
	ctx context.Context,
	neighborID string,
) (*api.Meta, api.Routes, error) {
	// Query birdwatcher
	_, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, nil, err
	}

	protocols := birdProtocols["protocols"].(map[string]interface{})

	if _, ok := protocols[neighborID]; !ok {
		return nil, nil, fmt.Errorf("invalid Neighbor")
	}

	peer := protocols[neighborID].(map[string]interface{})["neighbor_address"].(string)
	table := protocols[neighborID].(map[string]interface{})["table"].(string)
	pipe := src.getMasterPipeName(table)

	qryURL := "/routes/peer/" + peer
	if src.isAltSession(pipe) {
		qryURL = "/routes/table/" + table + "/peer/" + peer
	}
	if src.config.PipeProtocolLookup == "table" {
		qryURL = "/routes/table/" + table + "/peer/" + peer
	}

	// Query birdwatcher
	bird, err := src.client.GetJSON(ctx, qryURL)
	if err != nil {
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(bird, src.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	received, err := parseRoutes(bird, src.config, true)
	if err != nil {
		log.Println("WARNING Could not retrieve received routes:", err)
		log.Println("Is the 'routes_peer' module active in birdwatcher?")
		return apiStatus, nil, err
	}

	for k := range bird {
		delete(bird, k)
	}
	bird = nil

	return apiStatus, received, nil
}

func (src *MultiTableBirdwatcher) fetchPipeFilteredRoutes(
	ctx context.Context,
	protocols map[string]interface{},
	neighborID string,
	keepDetails bool,
) (*api.Meta, api.Routes, error) {
	neighborProto := protocols[neighborID].(map[string]interface{})
	table := neighborProto["table"].(string)

	pipes := []string{}

	for pid, p := range protocols {
		prot := p.(map[string]interface{})
		if prot["bird_protocol"] == "Pipe" && prot["table"] == table {
			pipes = append(pipes, pid)
		}
	}

	filtered := api.Routes{}
	var meta *api.Meta

	for _, pipe := range pipes {
		res, err := src.client.GetEndpoint(
			ctx,
			"/routes/pipe/filtered?table="+table+"&pipe="+pipe+"&protocol="+neighborID)
		if err != nil {
			log.Println("WARNING Could not retrieve filtered routes:", err)
			log.Println("Is the 'pipe_filtered' module active in birdwatcher?")
			return nil, nil, err
		}
		defer res.Body.Close()

		// Parse the routes
		m, pipeFiltered, err := parseRoutesResponseStream(res.Body, src.config, keepDetails)
		if err != nil {
			return nil, nil, err
		}
		filtered = append(filtered, pipeFiltered...)
		meta = m
	}

	/*

		// Query birdwatcher
		res, err = src.client.GetEndpoint(
			ctx,
			"/routes/pipe/filtered/protocol?table="+table+"&pipe="+pipeName+"&protocol="+protocolID)
		if err != nil {
			log.Println("WARNING Could not retrieve filtered routes:", err)
			log.Println("Is the 'pipe_filtered' module active in birdwatcher?")
			return meta, nil, err
		}
		defer res.Body.Close()

		// Parse the routes
		_, pipeFiltered, err := parseRoutesResponseStream(res.Body, src.config, keepDetails)
		if err != nil {
			return nil, nil, err
		}
	*/

	return meta, filtered, nil
}

func (src *MultiTableBirdwatcher) fetchFilteredRoutesStream(
	ctx context.Context,
	protocols map[string]interface{},
	neighborID string,
	keepDetails bool,
) (*api.Meta, api.Routes, error) {
	neighborProto, ok := protocols[neighborID].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("invalid Neighbor")
	}
	table := neighborProto["table"].(string)

	if src.config.PipeProtocolLookup == "table" {
		return src.fetchPipeFilteredRoutes(ctx, protocols, neighborID, keepDetails)
	}

	// Stage 1 filters
	res, err := src.client.GetEndpoint(ctx, "/routes/filtered/"+neighborID)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'routes_filtered' module active in birdwatcher?")
		return nil, nil, err
	}
	defer res.Body.Close()

	meta, filtered, err := parseRoutesResponseStream(res.Body, src.config, keepDetails)
	if err != nil {
		return nil, nil, err
	}

	// Stage 2 filters
	pipeName := src.getMasterPipeName(table)

	// If there is no pipe to master, there is nothing left to do
	if pipeName == "" {
		return meta, filtered, nil
	}

	// Check if this is an alternative session and query the alt pipe instead
	if src.isAltSession(pipeName) {
		pipeName = src.getAltPipeName(pipeName)
	}

	// Query birdwatcher
	res, err = src.client.GetEndpoint(
		ctx,
		"/routes/pipe/filtered?table="+table+"&pipe="+pipeName)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'pipe_filtered' module active in birdwatcher?")
		return meta, nil, err
	}
	defer res.Body.Close()

	// Parse the routes
	_, pipeFiltered, err := parseRoutesResponseStream(res.Body, src.config, keepDetails)
	if err != nil {
		return nil, nil, err
	}

	filtered = append(filtered, pipeFiltered...)

	return meta, filtered, nil
}

/*
func (src *MultiTableBirdwatcher) fetchFilteredRoutes(
	ctx context.Context,
	neighborID string,
	keepDetails bool,
) (*api.Meta, api.Routes, error) {
	// Query birdwatcher
	_, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, nil, err
	}

	protocols := birdProtocols["protocols"].(map[string]interface{})

	if _, ok := protocols[neighborID]; !ok {
		return nil, nil, fmt.Errorf("invalid Neighbor")
	}

	// Stage 1 filters
	birdFiltered, err := src.client.GetJSON(ctx, "/routes/filtered/"+neighborID)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'routes_filtered' module active in birdwatcher?")
		return nil, nil, err
	}

	// Use api status from first request
	apiStatus, err := parseAPIStatus(birdFiltered, src.config)
	if err != nil {
		return nil, nil, err
	}

	// Parse the routes
	filtered := parseRoutesData(
		birdFiltered["routes"].([]interface{}), src.config, keepDetails)

	// Stage 2 filters
	table := protocols[neighborID].(map[string]interface{})["table"].(string)
	pipeName := src.getMasterPipeName(table)

	// If there is no pipe to master, there is nothing left to do
	if pipeName == "" {
		return apiStatus, filtered, nil
	}

	// Check if this is an alternative session and query the alt pipe instead
	if src.isAltSession(pipeName) {
		pipeName = src.getAltPipeName(pipeName)
	}

	// Query birdwatcher
	birdPipeFiltered, err := src.client.GetJSON(
		ctx, "/routes/pipe/filtered?table="+table+"&pipe="+pipeName)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'pipe_filtered' module active in birdwatcher?")
		return apiStatus, nil, err
	}

	// Parse the routes
	pipeFiltered := parseRoutesData(
		birdPipeFiltered["routes"].([]interface{}), src.config, keepDetails)

	// Sort routes for deterministic ordering
	filtered = append(filtered, pipeFiltered...)

	if !keepDetails {
		// Yes this is not the right variable name to convey this...
		sort.Sort(filtered)
	}

	return apiStatus, filtered, nil
}
*/

func (src *MultiTableBirdwatcher) fetchNotExportedRoutes(
	ctx context.Context,
	neighborID string,
) (*api.Meta, api.Routes, error) {
	// Query birdwatcher
	apiStatus, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, nil, err
	}

	protocols := birdProtocols["protocols"].(map[string]interface{})

	if _, ok := protocols[neighborID]; !ok {
		return nil, nil, fmt.Errorf("invalid neighbor")
	}

	table := protocols[neighborID].(map[string]interface{})["table"].(string)
	pipeName := src.getMasterPipeName(table)

	// Check if this is a monitoring session, if so return no routes
	// as a monitoring session never export routes. We reuse the apiStatus
	// from the fetchProtocols.
	if src.isAltSession(pipeName) {
		return apiStatus, api.Routes{}, nil
	}

	// Query birdwatcher
	bird, _ := src.client.GetJSON(ctx, "/routes/noexport/"+pipeName)

	// Use api status from first request
	apiStatus, err = parseAPIStatus(bird, src.config)
	if err != nil {
		return nil, nil, err
	}

	notExported, err := parseRoutes(bird, src.config, true)
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
func (src *MultiTableBirdwatcher) fetchRequiredRoutes(
	ctx context.Context,
	protocols map[string]interface{},
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
	_, filteredRoutes, err := src.fetchFilteredRoutesStream(ctx, protocols, neighborID, true)
	if err != nil {
		return nil, err
	}

	// Perform route deduplication
	importedRoutes := api.Routes{}
	if len(receivedRoutes) > 0 {
		// TODO: maybe we can utilize the ptr here
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

// Fetch Neighbors and map the corresponding pipe protocols
// by the 'table' attribute.
func (src *MultiTableBirdwatcher) fetchNeighborsPipeTable(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	// Query birdwatcher for all protocols (incl. pipes)
	apiStatus, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, err
	}

	bgpProtos := src.filterProtocolsBgp(birdProtocols)
	bgpProtocols := bgpProtos["protocols"].(map[string]interface{})
	pipes := src.filterProtocolsPipe(
		birdProtocols)["protocols"].(map[string]interface{})

	// Get imported routes count on protocols. Mapped by Table
	bgpCount := make(map[string]float64)
	for _, p := range bgpProtocols {
		proto := p.(map[string]interface{})
		table := proto["table"].(string)
		routes := proto["routes"].(map[string]interface{})
		imported, ok := routes["imported"].(float64)
		if ok {
			bgpCount[table] = imported
		}
	}

	// To lookup the filtered routes, we need to check the number
	// of routes exported on the pipes.
	//  Rx -> n[Px | table: Tx]
	pipeFilterCount := make(map[string]float64)
	for _, p := range pipes {
		pipe := p.(map[string]interface{})

		table := pipe["table"].(string)
		routes := pipe["routes"].(map[string]interface{})

		count := routes["exported"].(float64) // sigh
		recCount := bgpCount[table]

		filtered := recCount - count

		pf, ok := pipeFilterCount[table]
		if ok {
			pf = pf + filtered
		} else {
			pf = filtered
		}
		pipeFilterCount[table] = pf
	}

	// Update protocols
	for _, p := range bgpProtocols {
		proto := p.(map[string]interface{})
		routes := proto["routes"].(map[string]interface{})
		table := proto["table"].(string)
		filt, ok := pipeFilterCount[table]
		if ok {
			routes["filtered"] = filt
		}
	}

	// Make neighbors
	neighbors, err := parseNeighbors(bgpProtos, src.config)
	if err != nil {
		return nil, err
	}

	response := &api.NeighborsResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		Neighbors: neighbors,
	}

	return response, nil
}

// Fetch neighbors (classic)
func (src *MultiTableBirdwatcher) fetchNeighborsPipeMaster(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	// Query birdwatcher
	apiStatus, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, err
	}

	// Parse the neighbors
	neighbors, err := parseNeighbors(
		src.filterProtocolsBgp(birdProtocols), src.config)
	if err != nil {
		return nil, err
	}

	pipes := src.filterProtocolsPipe(
		birdProtocols)["protocols"].(map[string]interface{})
	tree := src.parseProtocolToTableTree(birdProtocols)

	// Now determine the session count for each neighbor and check if the pipe
	// did filter anything
	filtered := make(map[string]int)
	for table := range tree {
		allRoutesImported := int64(0)
		pipeRoutesImported := int64(0)

		// Sum up all routes from all peers for a table
		for _, protocol := range tree[table].(map[string]interface{}) {
			// Skip peers that are not up (start/down)
			if !isProtocolUp(protocol.(map[string]interface{})["state"].(string)) {
				continue
			}
			allRoutesImported += int64(protocol.(map[string]interface{})["routes"].(map[string]interface{})["imported"].(float64))

			table := protocol.(map[string]interface{})["table"].(string)
			pipeName := src.getMasterPipeName(table)

			// Check if this is an alternative session and query the alt pipe instead
			if src.isAltSession(pipeName) {
				pipeName = src.getAltPipeName(pipeName)
			}

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
				filtered[protocol.(map[string]interface{})["protocol"].(string)] = int(allRoutesImported - pipeRoutesImported)
			}
		} else {
			// Multiple routers
			if pipeRoutesImported == 0 {
				// 0 is a special condition, which means that the pipe did filter ALL routes of
				// all peers. Therefore we already know the amount of filtered routes and don't have
				// to query birdwatcher again.
				for _, protocol := range tree[table].(map[string]interface{}) {
					// Skip peers that are not up (start/down)
					if !isProtocolUp(protocol.(map[string]interface{})["state"].(string)) {
						continue
					}
					filtered[protocol.(map[string]interface{})["protocol"].(string)] = int(protocol.(map[string]interface{})["routes"].(map[string]interface{})["imported"].(float64))
				}
			} else {
				// Otherwise the pipe did import at least some routes which means that
				// we have to query birdwatcher to get the count for each peer.
				for neighborAddress, protocol := range tree[table].(map[string]interface{}) {
					table := protocol.(map[string]interface{})["table"].(string)
					pipe := src.getMasterPipeName(table)
					// Check if this is an alternative session and query the alt pipe instead
					if src.isAltSession(pipe) {
						pipe = src.getAltPipeName(pipe)
					}

					count, err := src.client.GetJSON(ctx, "/routes/pipe/filtered/count?table="+table+"&pipe="+pipe+"&address="+neighborAddress)
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
	for _, neighbor := range neighbors {
		if pipeRoutesFiltered, ok := filtered[neighbor.ID]; ok {
			neighbor.RoutesAccepted -= pipeRoutesFiltered
			neighbor.RoutesFiltered += pipeRoutesFiltered
		}
	}

	response := &api.NeighborsResponse{
		Response: api.Response{
			Meta: apiStatus,
		},
		Neighbors: neighbors,
	}
	return response, nil
}

// Neighbors get neighbors from protocols.
// TODO: this. needs. refactoring.
func (src *MultiTableBirdwatcher) Neighbors(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	// Check if we hit the cache
	response := src.neighborsCache.Get()
	if response != nil {
		return response, nil
	}

	// We can use the table to map related pipe protocols or try
	// to derive by prefix.
	if src.config.PipeProtocolLookup == "table" {
		res, err := src.fetchNeighborsPipeTable(ctx)
		if err != nil {
			return nil, err
		}
		response = res
	} else {
		res, err := src.fetchNeighborsPipeMaster(ctx)
		if err != nil {
			return nil, err
		}
		response = res
	}

	// Cache result
	src.neighborsCache.Set(response)

	return response, nil // dereference for now
}

// NeighborsSummary is for now using Neighbors
func (src *MultiTableBirdwatcher) NeighborsSummary(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	return src.Neighbors(ctx)
}

// Routes gets filtered and exported route
// from the birdwatcher backend.
func (src *MultiTableBirdwatcher) Routes(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Query birdwatcher
	_, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, err
	}
	protocols := birdProtocols["protocols"].(map[string]interface{})

	// Fetch required routes first (received and filtered)
	// However: Store in separate cache for faster access
	required, err := src.fetchRequiredRoutes(ctx, protocols, neighborID)
	if err != nil {
		return nil, err
	}

	// Optional: NoExport
	_, notExported, err := src.fetchNotExportedRoutes(ctx, neighborID)
	if err != nil {
		return nil, err
	}

	response.Response.Meta = required.Meta
	response.Imported = required.Imported
	response.Filtered = required.Filtered
	response.NotExported = notExported

	return response, nil
}

// RoutesReceived returns all received routes
func (src *MultiTableBirdwatcher) RoutesReceived(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Check if we have a cache hit
	cachedRoutes := src.routesRequiredCache.Get(neighborID)
	if cachedRoutes != nil {
		response.Response.Meta = cachedRoutes.Response.Meta
		response.Imported = cachedRoutes.Imported
		return response, nil
	}

	// Query birdwatcher
	_, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, err
	}
	protocols := birdProtocols["protocols"].(map[string]interface{})

	// Fetch required routes first (received and filtered)
	routes, err := src.fetchRequiredRoutes(ctx, protocols, neighborID)
	if err != nil {
		return nil, err
	}

	response.Meta = routes.Meta
	response.Imported = routes.Imported

	return response, nil
}

// RoutesFiltered gets all filtered routes from the backend
func (src *MultiTableBirdwatcher) RoutesFiltered(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{}

	// Check if we have a cache hit
	cachedRoutes := src.routesRequiredCache.Get(neighborID)
	if cachedRoutes != nil {
		response.Meta = cachedRoutes.Meta
		response.Filtered = cachedRoutes.Filtered
		return response, nil
	}

	// Query birdwatcher
	_, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, err
	}
	protocols := birdProtocols["protocols"].(map[string]interface{})

	// Fetch required routes first (received and filtered)
	routes, err := src.fetchRequiredRoutes(ctx, protocols, neighborID)
	if err != nil {
		return nil, err
	}

	response.Meta = routes.Meta
	response.Filtered = routes.Filtered

	return response, nil
}

// RoutesNotExported gets all not exported routes
func (src *MultiTableBirdwatcher) RoutesNotExported(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	// Check if we have a cache hit
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

// AllRoutes retrieves a routes dump from the server
func (src *MultiTableBirdwatcher) AllRoutes(
	ctx context.Context,
) (*api.RoutesResponse, error) {
	// Query birdwatcher
	_, birdProtocols, err := src.fetchProtocols(ctx)
	if err != nil {
		return nil, err
	}
	protocols := birdProtocols["protocols"].(map[string]interface{})
	mainTable := src.GenericBirdwatcher.config.MainTable

	// Fetch received routes first
	res, err := src.client.GetEndpoint(ctx, "/routes/table/"+mainTable)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	meta, imported, err := parseRoutesResponseStream(res.Body, src.config, false)
	if err != nil {
		return nil, err
	}

	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: meta,
		},
	}

	// Sort routes for deterministic ordering
	// sort.Sort(imported)
	response.Imported = imported

	// Iterate over all the protocols and fetch the filtered routes for everyone
	protocolsBgp := src.filterProtocolsBgp(birdProtocols)

	// We load the filtered routes asynchronously with workers.
	type fetchFilteredReq struct {
		protocolID string
		peer       *string
		learntFrom *string
	}
	reqQ := make(chan fetchFilteredReq, 1000)
	resQ := make(chan api.Routes, 1000)
	wg := &sync.WaitGroup{}

	// Start workers
	for i := 0; i < 42; i++ {
		wg.Add(1)
		go func() {
			// This is a worker for fetching filtered routes
			defer wg.Done()
			for req := range reqQ {
				// Fetch filtered routes
				_, filtered, err := src.fetchFilteredRoutesStream(ctx, protocols, req.protocolID, false)
				if err != nil {
					log.Println("error while fetching filtered routes:", err)
				}
				// Perform route deduplication
				filtered = src.filterRoutesByPeerOrLearntFrom(
					filtered, req.peer, req.learntFrom)

				resQ <- filtered
			}
		}()
	}

	gwpool := pools.Gateways4

	// Fill request queue
	go func() {
		for protocolID, protocolsData := range protocolsBgp["protocols"].(map[string]interface{}) {
			peer := gwpool.Acquire(
				protocolsData.(map[string]interface{})["neighbor_address"].(string))
			learntFrom := gwpool.Acquire(
				decoders.String(protocolsData.(map[string]interface{})["learnt_from"], *peer))

			reqQ <- fetchFilteredReq{
				protocolID: protocolID,
				peer:       peer,
				learntFrom: learntFrom,
			}
		}
		close(reqQ)
	}()

	// Await all workers done and close channel
	go func() {
		wg.Wait()
		close(resQ)
	}()

	// Collect results
	for filtered := range resQ {
		response.Filtered = append(response.Filtered, filtered...)
	}

	return response, nil
}
