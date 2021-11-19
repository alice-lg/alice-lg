package store

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
)

// RoutesStoreBackend interface
type RoutesStoreBackend interface {
	// SetRoutes updates the routes in the store after a refresh.
	SetRoutes(
		ctx context.Context,
		sourceID string,
		routes *api.RoutesResponse,
	) error
}

// The RoutesStore holds a mapping of routes,
// status and cfgs and will be queried instead
// of a backend by the API
type RoutesStore struct {
	backend        RoutesStoreBackend
	sources        *SourceStore
	neighborsStore *NeighborsStore

	sync.RWMutex
}

// NewRoutesStore makes a new store instance
// with a cfg.
func NewRoutesStore(
	neighborsStore *NeighborsStore,
	cfg *config.Config,
	backend RoutesStoreBackend,
) *RoutesStore {
	// Set refresh interval as duration, fall back to
	// five minutes if no interval is set.
	refreshInterval := time.Duration(
		cfg.Server.RoutesStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}

	log.Println("Neighbors Store refresh interval set to:", refreshInterval)

	// Store refresh information per store
	sources := NewSourcesStore(cfg, refreshInterval)

	store := &RoutesStore{
		backend:        backend,
		neighborsStore: neighborsStore,
	}
	return store
}

// Start starts the routes store
func (s *RoutesStore) Start() {
	log.Println("Starting local routes store")
	if err := s.init(); err != nil {
		log.Fatal(err)
	}
}

// Service initialization
func (s *RoutesStore) init() error {
	// Periodically trigger updates
	for {
		s.update()
		time.Sleep(time.Second)
	}
}

// Update all routes
func (s *RoutesStore) update() {
	successCount := 0
	errorCount := 0
	t0 := time.Now()

	for sourceID := range s.routesMap {
		sourceConfig := s.cfgMap[sourceID]
		source := sourceConfig.GetInstance()

		// Get current update state
		if s.statusMap[sourceID].State == StateUpdating {
			continue // nothing to do here
		}

		// Set update state
		s.Lock()
		s.statusMap[sourceID] = Status{
			State: StateUpdating,
		}
		s.Unlock()

		routes, err := source.AllRoutes()
		if err != nil {
			log.Println(
				"Refreshing the routes store failed for:", sourceConfig.Name,
				"(", sourceConfig.ID, ")",
				"with:", err,
				"- NEXT STATE: ERROR",
			)

			s.Lock()
			s.statusMap[sourceID] = Status{
				State:       StateError,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			s.Unlock()

			errorCount++
			continue
		}

		s.Lock()
		// Update data
		s.routesMap[sourceID] = routes
		// Update state
		s.statusMap[sourceID] = Status{
			LastRefresh: time.Now(),
			State:       StateReady,
		}
		s.lastRefresh = time.Now().UTC()
		s.Unlock()

		successCount++
	}

	refreshDuration := time.Since(t0)
	log.Println(
		"Refreshed routes store for", successCount, "of", successCount+errorCount,
		"sources with", errorCount, "error(s) in", refreshDuration,
	)

}

// Stats calculates some store insights
func (s *RoutesStore) Stats() *api.RoutesStoreStats {
	totalImported := 0
	totalFiltered := 0

	rsStats := []api.RouteServerRoutesStats{}

	s.RLock()
	for sourceID, routes := range s.routesMap {
		status := s.statusMap[sourceID]

		totalImported += len(routes.Imported)
		totalFiltered += len(routes.Filtered)

		serverStats := api.RouteServerRoutesStats{
			Name: s.cfgMap[sourceID].Name,

			Routes: api.RoutesStats{
				Filtered: len(routes.Filtered),
				Imported: len(routes.Imported),
			},

			State:     stateToString(status.State),
			UpdatedAt: status.LastRefresh,
		}

		rsStats = append(sStats, serverStats)
	}
	s.RUnlock()

	// Make stats
	storeStats := &api.RoutesStoreStats{
		TotalRoutes: api.RoutesStats{
			Imported: totalImported,
			Filtered: totalFiltered,
		},
		RouteServers: rsStats,
	}
	return storeStats
}

// CachedAt provides a cache status
func (s *RoutesStore) CachedAt() time.Time {
	return s.lastRefresh
}

// CacheTTL returns the TTL time
func (s *RoutesStore) CacheTTL() time.Time {
	return s.lastRefresh.Add(s.refreshInterval)
}

// Lookup routes transform
func routeToLookupRoute(
	nStore *NeighborsStore,
	source *config.SourceConfig,
	state string,
	route *api.Route,
) *api.LookupRoute {
	// Get neighbor and make route
	neighbor := nStore.GetNeighborAt(source.ID, route.NeighborID)
	lookup := &api.LookupRoute{
		Route:    route,
		State:    state,
		Neighbor: neighbor,
		RouteServer: &api.RouteServer{
			ID:   source.ID,
			Name: source.Name,
		},
	}
	return lookup
}

// Routes filter
func filterRoutesByPrefix(
	nStore *NeighborsStore,
	source *config.SourceConfig,
	routes api.Routes,
	prefix string,
	state string,
) api.LookupRoutes {
	results := api.LookupRoutes{}
	for _, route := range routes {
		// Naiive filtering:
		if strings.HasPrefix(strings.ToLower(route.Network), prefix) {
			lookup := routeToLookupRoute(nStore, source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}

func filterRoutesByNeighborIDs(
	nStore *NeighborsStore,
	source *config.SourceConfig,
	routes api.Routes,
	neighborIDs []string,
	state string,
) api.LookupRoutes {

	results := api.LookupRoutes{}
	for _, route := range routes {
		// Filtering:
		if MemberOf(neighborIDs, route.NeighborID) {
			lookup := routeToLookupRoute(nStore, source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}

// LookupNeighborsPrefixesAt performs a single route server
// routes lookup by neighbor id
func (s *RoutesStore) LookupNeighborsPrefixesAt(
	sourceID string,
	neighborIDs []string,
) chan api.LookupRoutes {
	response := make(chan api.LookupRoutes)

	go func() {
		s.RLock()
		source := s.cfgMap[sourceID]
		routes := s.routesMap[sourceID]
		s.RUnlock()

		filtered := filterRoutesByNeighborIDs(
			s.neighborsStore,
			source,
			routes.Filtered,
			neighborIDs,
			"filtered")
		imported := filterRoutesByNeighborIDs(
			s.neighborsStore,
			source,
			routes.Imported,
			neighborIDs,
			"imported")

		result := append(filtered, imported...)
		response <- result
	}()

	return response
}

// LookupPrefixAt performs a single RS lookup
func (s *RoutesStore) LookupPrefixAt(
	sourceID string,
	prefix string,
) chan api.LookupRoutes {

	response := make(chan api.LookupRoutes)

	go func() {
		s.RLock()
		cfg := s.cfgMap[sourceID]
		routes := s.routesMap[sourceID]
		s.RUnlock()

		filtered := filterRoutesByPrefix(
			s.neighborsStore,
			cfg,
			routes.Filtered,
			prefix,
			"filtered")
		imported := filterRoutesByPrefix(
			s.neighborsStore,
			cfg,
			routes.Imported,
			prefix,
			"imported")

		result := append(filtered, imported...)
		response <- result
	}()

	return response
}

// LookupPrefix performs a lookup over all route servers
func (s *RoutesStore) LookupPrefix(prefix string) api.LookupRoutes {
	result := api.LookupRoutes{}
	responses := []chan api.LookupRoutes{}

	// Normalize prefix to lower case
	prefix = strings.ToLower(prefix)

	// Dispatch
	s.RLock()
	for sourceID := range s.routesMap {
		res := s.LookupPrefixAt(sourceID, prefix)
		responses = append(responses, res)
	}
	s.RUnlock()

	// Collect
	for _, response := range responses {
		routes := <-response
		result = append(result, routes...)
		close(response)
	}

	return result
}

// LookupPrefixForNeighbors returns all routes for
// a set of neighbos.
func (s *RoutesStore) LookupPrefixForNeighbors(
	neighbors api.NeighborsLookupResults,
) api.LookupRoutes {

	result := api.LookupRoutes{}
	responses := []chan api.LookupRoutes{}

	// Dispatch
	for sourceID, locals := range neighbors {
		lookupNeighborIDs := []string{}
		for _, n := range locals {
			lookupNeighborIDs = append(lookupNeighborIDs, n.ID)
		}
		res := s.LookupNeighborsPrefixesAt(sourceID, lookupNeighborIDs)
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
