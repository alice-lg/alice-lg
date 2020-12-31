package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/backend/api"
)

type RoutesStore struct {
	routesMap map[string]*api.RoutesResponse
	statusMap map[string]StoreStatus
	configMap map[string]*SourceConfig

	refreshInterval time.Duration
	lastRefresh     time.Time

	sync.RWMutex
}

func NewRoutesStore(config *Config) *RoutesStore {

	// Build mapping based on source instances
	routesMap := make(map[string]*api.RoutesResponse)
	statusMap := make(map[string]StoreStatus)
	configMap := make(map[string]*SourceConfig)

	for _, source := range config.Sources {
		id := source.Id

		configMap[id] = source
		routesMap[id] = &api.RoutesResponse{}
		statusMap[id] = StoreStatus{
			State: STATE_INIT,
		}
	}

	// Set refresh interval as duration, fall back to
	// five minutes if no interval is set.
	refreshInterval := time.Duration(
		config.Server.RoutesStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}

	store := &RoutesStore{
		routesMap:       routesMap,
		statusMap:       statusMap,
		configMap:       configMap,
		refreshInterval: refreshInterval,
	}
	return store
}

func (self *RoutesStore) Start() {
	log.Println("Starting local routes store")
	log.Println("Routes Store refresh interval set to:", self.refreshInterval)
	go self.init()
}

// Service initialization
func (self *RoutesStore) init() {
	// Initial refresh
	self.update()

	// Initial stats
	self.Stats().Log()

	// Periodically update store
	for {
		time.Sleep(self.refreshInterval)
		self.update()
	}
}

// Update all routes
func (self *RoutesStore) update() {
	type ResultType string
	const(
		Success = "Success"
		Failure = "Failure"
		Skipped = "Skipped"
	)

	type GetRoutesResult struct {
		sourceId string
		result  ResultType
	}

	result := make(chan GetRoutesResult)
	successCount := 0
	errorCount := 0
	t0 := time.Now()

	for sourceIdentifier, _ := range self.routesMap {
		go func(sourceId string) {
			from := time.Now()

			// Get current update state
			self.Lock()
			if self.statusMap[sourceId].State == STATE_UPDATING {
				self.Unlock()
				// Nothing to do here
				result <- GetRoutesResult{
					sourceId: sourceId,
					result: Skipped}
			} else {
				// Set update state
				self.statusMap[sourceId] = StoreStatus{
					State: STATE_UPDATING,
				}
				sourceConfig := self.configMap[sourceId]
				instance := sourceConfig.getInstance()
				self.Unlock()

				routes, err := instance.AllRoutes()
				if err != nil {
					log.Println(
						"Refreshing the routes store failed for:", sourceConfig.Name,
						"(", sourceConfig.Id, ")",
						"with:", err,
						"- NEXT STATE: ERROR",
					)

					self.Lock()
					self.statusMap[sourceId] = StoreStatus{
						State:       STATE_ERROR,
						LastError:   err,
						LastRefresh: time.Now(),
					}
					self.Unlock()

					result <- GetRoutesResult{
						sourceId: sourceId,
						result: Failure}
				} else {
					self.Lock()
					// Update data
					self.routesMap[sourceId] = routes
					// Update state
					self.statusMap[sourceId] = StoreStatus{
						LastRefresh: time.Now(),
						State:       STATE_READY,
					}
					self.lastRefresh = time.Now().UTC()
					self.Unlock()
					log.Println("Refreshed: ", sourceId, " filtered: ", len(routes.Filtered), " imported: ", len(routes.Imported), " not exported: ", len(routes.NotExported),
						" from: ", from, " to: ", time.Now())

					result <- GetRoutesResult{
						sourceId: sourceId,
						result: Success}
				}
			}
		}(sourceIdentifier)
	}

	for i := 0; i < len(self.routesMap); i++ {
		switch (<- result).result {
		case Success:
			successCount++
		case Failure:
			errorCount++
		}
	}
	close(result)

	refreshDuration := time.Since(t0)
	log.Println(
		"Refreshed routes store for", successCount, "of", successCount+errorCount,
		"sources with", errorCount, "error(s) in", refreshDuration,
	)

}

// Calculate store insights
func (self *RoutesStore) Stats() RoutesStoreStats {
	totalImported := 0
	totalFiltered := 0

	rsStats := []RouteServerRoutesStats{}

	self.RLock()
	for sourceId, routes := range self.routesMap {
		status := self.statusMap[sourceId]

		totalImported += len(routes.Imported)
		totalFiltered += len(routes.Filtered)

		serverStats := RouteServerRoutesStats{
			Name: self.configMap[sourceId].Name,

			Routes: RoutesStats{
				Filtered: len(routes.Filtered),
				Imported: len(routes.Imported),
			},

			State:     stateToString(status.State),
			UpdatedAt: status.LastRefresh,
		}

		rsStats = append(rsStats, serverStats)
	}
	self.RUnlock()

	// Make stats
	storeStats := RoutesStoreStats{
		TotalRoutes: RoutesStats{
			Imported: totalImported,
			Filtered: totalFiltered,
		},
		RouteServers: rsStats,
	}
	return storeStats
}

// Provide cache status
func (self *RoutesStore) CachedAt() time.Time {
	return self.lastRefresh
}

func (self *RoutesStore) CacheTtl() time.Time {
	return self.lastRefresh.Add(self.refreshInterval)
}

// Lookup routes transform
func routeToLookupRoute(
	source *SourceConfig,
	state string,
	route *api.Route,
) *api.LookupRoute {

	// Get neighbour
	neighbour := AliceNeighboursStore.GetNeighbourAt(source.Id, route.NeighbourId)

	// Make route
	lookup := &api.LookupRoute{
		Id: route.Id,

		NeighbourId: route.NeighbourId,
		Neighbour:   neighbour,

		Routeserver: api.Routeserver{
			Id:   source.Id,
			Name: source.Name,
		},

		State: state,

		Network:   route.Network,
		Interface: route.Interface,
		Gateway:   route.Gateway,
		Metric:    route.Metric,
		Bgp:       route.Bgp,
		Age:       route.Age,
		Type:      route.Type,
		Primary:   route.Primary,
	}

	return lookup
}

// Routes filter
func filterRoutesByPrefix(
	source *SourceConfig,
	routes api.Routes,
	prefix string,
	state string,
) api.LookupRoutes {
	results := api.LookupRoutes{}
	for _, route := range routes {
		// Naiive filtering:
		if strings.HasPrefix(strings.ToLower(route.Network), prefix) {
			lookup := routeToLookupRoute(source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}

func filterRoutesByNeighbourIds(
	source *SourceConfig,
	routes api.Routes,
	neighbourIds []string,
	state string,
) api.LookupRoutes {

	results := api.LookupRoutes{}
	for _, route := range routes {
		// Filtering:
		if MemberOf(neighbourIds, route.NeighbourId) == true {
			lookup := routeToLookupRoute(source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}

// Single RS lookup by neighbour id
func (self *RoutesStore) LookupNeighboursPrefixesAt(
	sourceId string,
	neighbourIds []string,
) chan api.LookupRoutes {
	response := make(chan api.LookupRoutes)

	go func() {
		self.RLock()
		source := self.configMap[sourceId]
		routes := self.routesMap[sourceId]
		self.RUnlock()

		filtered := filterRoutesByNeighbourIds(
			source,
			routes.Filtered,
			neighbourIds,
			"filtered")
		imported := filterRoutesByNeighbourIds(
			source,
			routes.Imported,
			neighbourIds,
			"imported")

		var result api.LookupRoutes
		result = append(filtered, imported...)

		response <- result
	}()

	return response
}

// Single RS lookup
func (self *RoutesStore) LookupPrefixAt(
	sourceId string,
	prefix string,
) chan api.LookupRoutes {

	response := make(chan api.LookupRoutes)

	go func() {
		self.RLock()
		config := self.configMap[sourceId]
		routes := self.routesMap[sourceId]
		self.RUnlock()

		filtered := filterRoutesByPrefix(
			config,
			routes.Filtered,
			prefix,
			"filtered")
		imported := filterRoutesByPrefix(
			config,
			routes.Imported,
			prefix,
			"imported")

		var result api.LookupRoutes
		result = append(filtered, imported...)

		response <- result
	}()

	return response
}

func (self *RoutesStore) LookupPrefix(prefix string) api.LookupRoutes {
	result := api.LookupRoutes{}
	responses := []chan api.LookupRoutes{}

	// Normalize prefix to lower case
	prefix = strings.ToLower(prefix)

	// Dispatch
	self.RLock()
	for sourceId, _ := range self.routesMap {
		res := self.LookupPrefixAt(sourceId, prefix)
		responses = append(responses, res)
	}
	self.RUnlock()

	// Collect
	for _, response := range responses {
		routes := <-response
		result = append(result, routes...)
		close(response)
	}

	return result
}

func (self *RoutesStore) LookupPrefixForNeighbours(
	neighbours api.NeighboursLookupResults,
) api.LookupRoutes {

	result := api.LookupRoutes{}
	responses := []chan api.LookupRoutes{}

	// Dispatch
	for sourceId, locals := range neighbours {
		lookupNeighbourIds := []string{}
		for _, n := range locals {
			lookupNeighbourIds = append(lookupNeighbourIds, n.Id)
		}

		res := self.LookupNeighboursPrefixesAt(sourceId, lookupNeighbourIds)
		responses = append(responses, res)
	}

	// Collect
	for _, response := range responses {
		routes := <-response
		result = append(result, routes...)
		close(response)
	}

	return result
}
