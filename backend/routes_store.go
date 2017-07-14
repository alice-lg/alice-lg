package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ecix/alice-lg/backend/api"
)

type RoutesStore struct {
	routesMap map[int]api.RoutesResponse
	statusMap map[int]StoreStatus
	configMap map[int]SourceConfig

	rwlock *sync.RWMutex
}

func NewRoutesStore(config *Config) *RoutesStore {

	// Build mapping based on source instances
	routesMap := make(map[int]api.RoutesResponse)
	statusMap := make(map[int]StoreStatus)
	configMap := make(map[int]SourceConfig)

	for _, source := range config.Sources {
		id := source.Id

		configMap[id] = source
		routesMap[id] = api.RoutesResponse{}
		statusMap[id] = StoreStatus{
			State: STATE_INIT,
		}
	}

	store := &RoutesStore{
		routesMap: routesMap,
		statusMap: statusMap,
		configMap: configMap,

		rwlock: &sync.RWMutex{},
	}
	return store
}

func (self *RoutesStore) Start() {
	log.Println("Starting local routes store")
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
		// TODO: Add config option
		time.Sleep(5 * time.Minute)
		self.update()
	}
}

// Update all routes
func (self *RoutesStore) update() {
	for sourceId, _ := range self.routesMap {
		source := self.configMap[sourceId].getInstance()

		// Get current update state
		if self.statusMap[sourceId].State == STATE_UPDATING {
			continue // nothing to do here
		}

		// Set update state
		self.rwlock.Lock()
		self.statusMap[sourceId] = StoreStatus{
			State: STATE_UPDATING,
		}
		self.rwlock.Unlock()

		routes, err := source.AllRoutes()
		if err != nil {
			self.rwlock.Lock()
			self.statusMap[sourceId] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			self.rwlock.Unlock()

			continue
		}

		self.rwlock.Lock()
		// Update data
		self.routesMap[sourceId] = routes
		// Update state
		self.statusMap[sourceId] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
		self.rwlock.Unlock()
	}
}

// Calculate store insights
func (self *RoutesStore) Stats() RoutesStoreStats {
	totalImported := 0
	totalFiltered := 0

	rsStats := []RouteServerRoutesStats{}

	self.rwlock.RLock()
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
	self.rwlock.RUnlock()

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

// Lookup routes transform
func routeToLookupRoute(source SourceConfig, state string, route api.Route) api.LookupRoute {

	// Get neighbour
	neighbour := AliceNeighboursStore.GetNeighbourAt(source.Id, route.NeighbourId)

	// Make route
	lookup := api.LookupRoute{
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
	}

	return lookup
}

// Routes filter
func filterRoutesByPrefix(
	source SourceConfig,
	routes []api.Route,
	prefix string,
	state string,
) []api.LookupRoute {

	results := []api.LookupRoute{}
	for _, route := range routes {
		// Naiive filtering:
		if strings.HasPrefix(route.Network, prefix) {
			lookup := routeToLookupRoute(source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}

func filterRoutesByNeighbourIds(
	source SourceConfig,
	routes []api.Route,
	neighbourIds []string,
	state string,
) []api.LookupRoute {

	results := []api.LookupRoute{}
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
	sourceId int,
	neighbourIds []string,
) chan []api.LookupRoute {
	response := make(chan []api.LookupRoute)

	go func() {
		self.rwlock.RLock()
		source := self.configMap[sourceId]
		routes := self.routesMap[sourceId]
		self.rwlock.RUnlock()

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

		var result []api.LookupRoute
		result = append(filtered, imported...)

		response <- result
	}()

	return response
}

// Single RS lookup
func (self *RoutesStore) LookupPrefixAt(
	sourceId int,
	prefix string,
) chan []api.LookupRoute {

	response := make(chan []api.LookupRoute)

	go func() {
		self.rwlock.RLock()
		config := self.configMap[sourceId]
		routes := self.routesMap[sourceId]
		self.rwlock.RUnlock()

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

		var result []api.LookupRoute
		result = append(filtered, imported...)

		response <- result
	}()

	return response
}

func (self *RoutesStore) LookupPrefix(prefix string) []api.LookupRoute {
	result := []api.LookupRoute{}
	responses := []chan []api.LookupRoute{}

	// Dispatch
	self.rwlock.RLock()
	for sourceId, _ := range self.routesMap {
		res := self.LookupPrefixAt(sourceId, prefix)
		responses = append(responses, res)
	}
	self.rwlock.RUnlock()

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
) []api.LookupRoute {

	result := []api.LookupRoute{}
	responses := []chan []api.LookupRoute{}

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
