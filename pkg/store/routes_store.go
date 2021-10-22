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

// The RoutesStore holds a mapping of routes,
// status and cfgs and will be queried instead
// of a backend by the API
type RoutesStore struct {
	routesMap map[string]*api.RoutesResponse
	statusMap map[string]StoreStatus
	cfgMap    map[string]*config.SourceConfig

	refreshInterval time.Duration
	lastRefresh     time.Time

	neighborsStore *NeighborsStore

	sync.RWMutex
}

// NewRoutesStore makes a new store instance
// with a cfg.
func NewRoutesStore(
	neighborsStore *NeighborsStore,
	cfg *config.Config,
) *RoutesStore {
	// Build mapping based on source instances
	routesMap := make(map[string]*api.RoutesResponse)
	statusMap := make(map[string]StoreStatus)
	cfgMap := make(map[string]*config.SourceConfig)

	for _, source := range cfg.Sources {
		id := source.ID

		cfgMap[id] = source
		routesMap[id] = &api.RoutesResponse{}
		statusMap[id] = StoreStatus{
			State: STATE_INIT,
		}
	}

	// Set refresh interval as duration, fall back to
	// five minutes if no interval is set.
	refreshInterval := time.Duration(
		cfg.Server.RoutesStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}
	store := &RoutesStore{
		routesMap:       routesMap,
		statusMap:       statusMap,
		cfgMap:          cfgMap,
		refreshInterval: refreshInterval,
		neighborsStore:  neighborsStore,
	}
	return store
}

// Start starts the routes store
func (rs *RoutesStore) Start(ctx context.Context) {
	log.Println("Starting local routes store")
	log.Println("Routes Store refresh interval set to:", rs.refreshInterval)
	if err := rs.init(); err != nil {
		log.Fatal(err)
	}
}

// Service initialization
func (rs *RoutesStore) init() error {
	// Initial refresh
	rs.update()

	// Initial stats
	rs.Stats().Log()

	// Periodically update store
	for {
		time.Sleep(rs.refreshInterval)
		rs.update()
	}
}

// Update all routes
func (rs *RoutesStore) update() {
	successCount := 0
	errorCount := 0
	t0 := time.Now()

	for sourceID := range rs.routesMap {
		sourceConfig := rs.cfgMap[sourceID]
		source := sourceConfig.GetInstance()

		// Get current update state
		if rs.statusMap[sourceID].State == STATE_UPDATING {
			continue // nothing to do here
		}

		// Set update state
		rs.Lock()
		rs.statusMap[sourceID] = StoreStatus{
			State: STATE_UPDATING,
		}
		rs.Unlock()

		routes, err := source.AllRoutes()
		if err != nil {
			log.Println(
				"Refreshing the routes store failed for:", sourceConfig.Name,
				"(", sourceConfig.ID, ")",
				"with:", err,
				"- NEXT STATE: ERROR",
			)

			rs.Lock()
			rs.statusMap[sourceID] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			rs.Unlock()

			errorCount++
			continue
		}

		rs.Lock()
		// Update data
		rs.routesMap[sourceID] = routes
		// Update state
		rs.statusMap[sourceID] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
		rs.lastRefresh = time.Now().UTC()
		rs.Unlock()

		successCount++
	}

	refreshDuration := time.Since(t0)
	log.Println(
		"Refreshed routes store for", successCount, "of", successCount+errorCount,
		"sources with", errorCount, "error(s) in", refreshDuration,
	)

}

// Stats calculates some store insights
func (rs *RoutesStore) Stats() *api.RoutesStoreStats {
	totalImported := 0
	totalFiltered := 0

	rsStats := []RouteServerRoutesStats{}

	rs.RLock()
	for sourceID, routes := range rs.routesMap {
		status := rs.statusMap[sourceID]

		totalImported += len(routes.Imported)
		totalFiltered += len(routes.Filtered)

		serverStats := RouteServerRoutesStats{
			Name: rs.cfgMap[sourceID].Name,

			Routes: RoutesStats{
				Filtered: len(routes.Filtered),
				Imported: len(routes.Imported),
			},

			State:     stateToString(status.State),
			UpdatedAt: status.LastRefresh,
		}

		rsStats = append(rsStats, serverStats)
	}
	rs.RUnlock()

	// Make stats
	storeStats := &RoutesStoreStats{
		TotalRoutes: RoutesStats{
			Imported: totalImported,
			Filtered: totalFiltered,
		},
		RouteServers: rsStats,
	}
	return storeStats
}

// CachedAt provides a cache status
func (rs *RoutesStore) CachedAt() time.Time {
	return rs.lastRefresh
}

// CacheTTL returns the TTL time
func (rs *RoutesStore) CacheTTL() time.Time {
	return rs.lastRefresh.Add(rs.refreshInterval)
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
		if MemberOf(neighborIDs, route.NeighborID) == true {
			lookup := routeToLookupRoute(nStore, source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}

// LookupNeighborsPrefixesAt performs a single route server
// routes lookup by neighbor id
func (rs *RoutesStore) LookupNeighborsPrefixesAt(
	sourceID string,
	neighborIDs []string,
) chan api.LookupRoutes {
	response := make(chan api.LookupRoutes)

	go func() {
		rs.RLock()
		source := rs.cfgMap[sourceID]
		routes := rs.routesMap[sourceID]
		rs.RUnlock()

		filtered := filterRoutesByNeighborIDs(
			rs.neighborsStore,
			source,
			routes.Filtered,
			neighborIDs,
			"filtered")
		imported := filterRoutesByNeighborIDs(
			rs.neighborsStore,
			source,
			routes.Imported,
			neighborIDs,
			"imported")

		var result api.LookupRoutes
		result = append(filtered, imported...)

		response <- result
	}()

	return response
}

// LookupPrefixAt performs a single RS lookup
func (rs *RoutesStore) LookupPrefixAt(
	sourceID string,
	prefix string,
) chan api.LookupRoutes {

	response := make(chan api.LookupRoutes)

	go func() {
		rs.RLock()
		cfg := rs.cfgMap[sourceID]
		routes := rs.routesMap[sourceID]
		rs.RUnlock()

		filtered := filterRoutesByPrefix(
			rs.neighborsStore,
			cfg,
			routes.Filtered,
			prefix,
			"filtered")
		imported := filterRoutesByPrefix(
			rs.neighborsStore,
			cfg,
			routes.Imported,
			prefix,
			"imported")

		var result api.LookupRoutes
		result = append(filtered, imported...)

		response <- result
	}()

	return response
}

// LookupPrefix performs a lookup over all route servers
func (rs *RoutesStore) LookupPrefix(prefix string) api.LookupRoutes {
	result := api.LookupRoutes{}
	responses := []chan api.LookupRoutes{}

	// Normalize prefix to lower case
	prefix = strings.ToLower(prefix)

	// Dispatch
	rs.RLock()
	for sourceID := range rs.routesMap {
		res := rs.LookupPrefixAt(sourceID, prefix)
		responses = append(responses, res)
	}
	rs.RUnlock()

	// Collect
	for _, response := range responses {
		routes := <-response
		result = append(result, routes...)
		close(response)
	}

	return result
}

// LookupPrefixForNeighbors returns all routes for
// a set of neighbors.
func (rs *RoutesStore) LookupPrefixForNeighbors(
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
		res := rs.LookupNeighborsPrefixesAt(sourceID, lookupNeighborIDs)
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
