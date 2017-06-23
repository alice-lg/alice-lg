package main

import (
	"github.com/ecix/alice-lg/backend/api"
	"github.com/ecix/alice-lg/backend/sources"

	"log"
	"strings"
	"time"
)

type StoreStatus struct {
	LastRefresh time.Time
	LastError   error
	State       int
}

type RoutesStore struct {
	routesMap map[sources.Source]api.RoutesResponse
	statusMap map[sources.Source]StoreStatus
	configMap map[sources.Source]SourceConfig
}

func NewRoutesStore(config *Config) *RoutesStore {

	// Build mapping based on source instances
	routesMap := make(map[sources.Source]api.RoutesResponse)
	statusMap := make(map[sources.Source]StoreStatus)
	configMap := make(map[sources.Source]SourceConfig)

	for _, source := range config.Sources {
		instance := source.getInstance()
		configMap[instance] = source
		routesMap[instance] = api.RoutesResponse{}
		statusMap[instance] = StoreStatus{
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
	for source, _ := range self.routesMap {
		// Get current update state
		if self.statusMap[source].State == STATE_UPDATING {
			continue // nothing to do here
		}

		// Set update state
		self.statusMap[source] = StoreStatus{
			State: STATE_UPDATING,
		}

		routes, err := source.AllRoutes()
		if err != nil {
			self.statusMap[source] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}

			continue
		}

		// Update data
		self.routesMap[source] = routes
		// Update state
		self.statusMap[source] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
	}
}

// Helper: stateToString
func stateToString(state int) string {
	switch state {
	case STATE_INIT:
		return "INIT"
	case STATE_READY:
		return "READY"
	case STATE_UPDATING:
		return "UPDATING"
	case STATE_ERROR:
		return "ERROR"
	}
	return "INVALID"
}

// Calculate store insights
func (self *RoutesStore) Stats() RoutesStoreStats {
	totalImported := 0
	totalFiltered := 0

	rsStats := []RouteServerStats{}

	for source, routes := range self.routesMap {
		status := self.statusMap[source]

		totalImported += len(routes.Imported)
		totalFiltered += len(routes.Filtered)

		serverStats := RouteServerStats{
			Name: self.configMap[source].Name,

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

// Single RS lookup
func (self *RoutesStore) lookupRs(
	source sources.Source,
	prefix string,
) chan []api.LookupRoute {

	response := make(chan []api.LookupRoute)
	config := self.configMap[source]
	routes := self.routesMap[source]

	go func() {
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

		result := append(filtered, imported...)

		response <- result
	}()

	return response
}

func (self *RoutesStore) Lookup(prefix string) []api.LookupRoute {
	result := []api.LookupRoute{}
	responses := []chan []api.LookupRoute{}

	// Dispatch
	for source, _ := range self.routesMap {
		res := self.lookupRs(source, prefix)
		responses = append(responses, res)
	}

	// Collect
	for _, response := range responses {
		routes := <-response
		result = append(result, routes...)
	}

	return result
}
