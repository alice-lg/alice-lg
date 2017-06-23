package main

import (
	"github.com/ecix/alice-lg/backend/api"

	"log"
	"strings"
	"time"
)

type RoutesStore struct {
	routesMap map[int]api.RoutesResponse
	statusMap map[int]StoreStatus
	configMap map[int]SourceConfig
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
		self.statusMap[sourceId] = StoreStatus{
			State: STATE_UPDATING,
		}

		routes, err := source.AllRoutes()
		if err != nil {
			self.statusMap[sourceId] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}

			continue
		}

		// Update data
		self.routesMap[sourceId] = routes
		// Update state
		self.statusMap[sourceId] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
	}
}

// Calculate store insights
func (self *RoutesStore) Stats() RoutesStoreStats {
	totalImported := 0
	totalFiltered := 0

	rsStats := []RouteServerRoutesStats{}

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

// Routes filter
func filterRoutes(
	config SourceConfig,
	routes []api.Route,
	prefix string,
	state string,
) []api.LookupRoute {

	results := []api.LookupRoute{}

	for _, route := range routes {
		// Naiive filtering:
		if strings.HasPrefix(route.Network, prefix) {
			lookup := api.LookupRoute{
				Id:          route.Id,
				NeighbourId: route.NeighbourId,

				Routeserver: api.Routeserver{
					Id:   config.Id,
					Name: config.Name,
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
			results = append(results, lookup)
		}
	}
	return results
}

func addNeighbour(
	sourceId int,
	route api.LookupRoute,
) api.LookupRoute {
	neighbour := AliceNeighboursStore.GetNeighbourAt(
		sourceId, route.NeighbourId)
	route.Neighbour = neighbour
	return route
}

// Single RS lookup
func (self *RoutesStore) lookupRs(
	sourceId int,
	prefix string,
) chan []api.LookupRoute {

	response := make(chan []api.LookupRoute)
	config := self.configMap[sourceId]
	routes := self.routesMap[sourceId]

	go func() {
		result := []api.LookupRoute{}

		filtered := filterRoutes(
			config,
			routes.Filtered,
			prefix,
			"filtered")
		imported := filterRoutes(
			config,
			routes.Imported,
			prefix,
			"imported")

		// Add Neighbours to results
		for _, route := range filtered {
			result = append(result, addNeighbour(sourceId, route))
		}

		for _, route := range imported {
			result = append(result, addNeighbour(sourceId, route))
		}

		response <- result
	}()

	return response
}

func (self *RoutesStore) Lookup(prefix string) []api.LookupRoute {
	result := []api.LookupRoute{}
	responses := []chan []api.LookupRoute{}

	// Dispatch
	for sourceId, _ := range self.routesMap {
		res := self.lookupRs(sourceId, prefix)
		responses = append(responses, res)
	}

	// Collect
	for _, response := range responses {
		routes := <-response
		result = append(result, routes...)
	}

	return result
}
