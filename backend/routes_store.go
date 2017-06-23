package main

import (
	"github.com/ecix/alice-lg/backend/api"
	"github.com/ecix/alice-lg/backend/sources"

	"log"
	"time"
)

const (
	STATE_INIT = iota
	STATE_READY
	STATE_UPDATING
	STATE_ERROR
)

type RoutesStats struct {
	Filtered int `json:"filtered"`
	Imported int `json:"imported"`
}

type RouteServerStats struct {
	Name   string      `json:"name"`
	Routes RoutesStats `json:"routes"`

	State     string
	UpdatedAt time.Time `json:"updated_at"`
}

type StoreStats struct {
	TotalRoutes  RoutesStats        `json:"total_routes"`
	RouteServers []RouteServerStats `json:"route_servers"`
}

type StoreStatus struct {
	LastRefresh time.Time
	LastError   error
	State       int
}

type RoutesStore struct {
	routesMap map[sources.Source]api.RoutesResponse
	statusMap map[sources.Source]StoreStatus
}

func NewRoutesStore(config *Config) *RoutesStore {

	// Build mapping based on source instances
	routesMap := make(map[sources.Source]api.RoutesResponse)
	statusMap := make(map[sources.Source]StoreStatus)
	for _, source := range config.Sources {
		instance := source.getInstance()
		routesMap[instance] = api.RoutesResponse{}
		statusMap[instance] = StoreStatus{
			State: STATE_INIT,
		}
	}

	store := &RoutesStore{
		routesMap: routesMap,
		statusMap: statusMap,
	}
	return store
}

func (self *RoutesStore) Start() {
	log.Println("Starting local routes store")

	// Initial refresh
	go self.update()
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

	log.Println("All caches refreshed")
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
	return "UNKNOWN"
}

// Calculate store insights
func (self *RoutesStore) Stats() StoreStats {
	totalImported := 0
	totalFiltered := 0

	rsStats := []RouteServerStats{}

	for source, routes := range self.routesMap {
		status := self.statusMap[source]

		totalImported += len(routes.Imported)
		totalFiltered += len(routes.Filtered)

		serverStats := RouteServerStats{
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
	storeStats := StoreStats{
		TotalRoutes: RoutesStats{
			Imported: totalImported,
			Filtered: totalFiltered,
		},
		RouteServers: rsStats,
	}
	return storeStats

}

func (self *RoutesStore) Lookup(prefix string) []api.LookupRoute {
	result := []api.LookupRoute{}

	return result
}
