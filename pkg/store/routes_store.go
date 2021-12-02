package store

import (
	"context"
	"log"
	"math/rand"
	"strings"
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
		routes api.LookupRoutes,
	) error

	// CountRoutesAt returns the number of imported
	// and filtered routes for a given route server.
	// Example: (imported, filtered, error)
	CountRoutesAt(
		ctx context.Context,
		sourceID string,
	) (uint, uint, error)

	// GetNeighborPrefixesAt retrieves the prefixes
	// announced by the neighbor at a given source
	GetNeighborPrefixesAt(
		ctx context.Context,
		sourceID string,
		neighborID string,
	) (api.LookupRoutes, error)
}

// The RoutesStore holds a mapping of routes,
// status and cfgs and will be queried instead
// of a backend by the API
type RoutesStore struct {
	backend   RoutesStoreBackend
	sources   *SourcesStore
	neighbors *NeighborsStore
}

// NewRoutesStore makes a new store instance
// with a cfg.
func NewRoutesStore(
	neighbors *NeighborsStore,
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

	log.Println("Routes refresh interval set to:", refreshInterval)

	// Store refresh information per store
	sources := NewSourcesStore(cfg, refreshInterval)
	store := &RoutesStore{
		backend:   backend,
		sources:   sources,
		neighbors: neighbors,
	}
	return store
}

// Start starts the routes store
func (s *RoutesStore) Start() {
	log.Println("Starting local routes store")

	// Periodically trigger updates
	for {
		s.update()
		time.Sleep(time.Second)
	}
}

// Update all routes from all sources, where the
// sources last refresh is longer ago than the configured
// refresh period. This is totally the same as the
// NeighborsStore.update and maybe these functions can be merged (TODO)
func (s *RoutesStore) update() {
	for _, id := range s.sources.GetSourceIDs() {
		go s.safeUpdateSource(id)
	}
}

// safeUpdateSource will try to update a source but
// will recover from a panic if something goes wrong.
// In that case, the LastError and State will be updated.
// Again. The similarity to the NeighborsStore is really sus.
func (s *RoutesStore) safeUpdateSource(id string) {
	ctx := context.TODO()

	if !s.sources.ShouldRefresh(id) {
		return // Nothing to do here
	}

	if err := s.sources.LockSource(id); err != nil {
		log.Println("Cloud not start routes refresh:", err)
		return
	}

	// Apply jitter so, we do not hit everything at once.
	// TODO: Make configurable
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	src := s.sources.Get(id)

	// Prepare for impact.
	defer func() {
		if err := recover(); err != nil {
			log.Println(
				"Recovering after failed routes refresh of",
				src.Name, "from:", err)
			s.sources.RefreshError(id, err)
		}
	}()

	if err := s.updateSource(ctx, src); err != nil {
		log.Println(
			"Refeshing routes of", src.Name, "failed:", err)
		s.sources.RefreshError(id, err)
	}
}

// Update all routes
func (s *RoutesStore) updateSource(
	ctx context.Context,
	src *config.SourceConfig,
) error {
	if err := s.awaitNeighborStore(ctx, src.ID); err != nil {
		return err
	}

	rs := src.GetInstance()
	res, err := rs.AllRoutes()
	if err != nil {
		return err
	}

	// Prepare imported routes for lookup
	imported := s.routesToLookupRoutes("imported", src, res.Imported)
	filtered := s.routesToLookupRoutes("filtered", src, res.Filtered)
	lookupRoutes := append(imported, filtered...)

	if err = s.backend.SetRoutes(ctx, src.ID, lookupRoutes); err != nil {
		return err
	}

	return s.sources.RefreshSuccess(src.ID)
}

func (s *RoutesStore) routesToLookupRoutes(
	state string,
	src *config.SourceConfig,
	routes api.Routes,
) api.LookupRoutes {
	lookupRoutes := make(api.LookupRoutes, 0, len(routes))
	for _, route := range routes {
		neighbor, err := s.neighbors.GetNeighborAt(src.ID, route.NeighborID)
		if err != nil {
			log.Println("prepare route, neighbor lookup failed:", err)
			continue
		}
		lr := &api.LookupRoute{
			Route:    route,
			State:    state,
			Neighbor: neighbor,
			RouteServer: &api.RouteServer{
				ID:   src.ID,
				Name: src.Name,
			},
		}
		lookupRoutes = append(lookupRoutes, lr)
	}
	return lookupRoutes
}

func (s *RoutesStore) awaitNeighborStore(
	ctx context.Context,
	srcID string,
) error {
	// Poll the neighbor store state for the sourceID
	// until the context is not longer valid
	for {
		err := ctx.Err()
		if err != nil {
			return err
		}
		if s.neighbors.IsReady(srcID) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Stats calculates some store insights
func (s *RoutesStore) Stats() *api.RoutesStoreStats {
	ctx := context.TODO()

	totalImported := uint(0)
	totalFiltered := uint(0)

	rsStats := []api.RouteServerRoutesStats{}

	for _, sourceID := range s.sources.GetSourceIDs() {
		status, err := s.sources.GetStatus(sourceID)
		if err != nil {
			log.Println("error while getting source status:", err)
			continue
		}

		src := s.sources.Get(sourceID)

		nImported, nFiltered, err := s.backend.CountRoutesAt(ctx, sourceID)
		if err != nil {
			log.Println("error during routes count:", err)
		}

		totalImported += nImported
		totalFiltered += nFiltered

		serverStats := api.RouteServerRoutesStats{
			Name: src.Name,
			Routes: api.RoutesStats{
				Imported: nImported,
				Filtered: nFiltered,
			},
			State:     status.State.String(),
			UpdatedAt: status.LastRefresh,
		}
		rsStats = append(sStats, serverStats)
	}

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

// SourceCachedAt provides a cache status (TODO: do we need this?)
func (s *RoutesStore) SourceCachedAt(id string) time.Time {
	status, err := s.sources.GetStatus(sourceID)
	if err != nil {
		log.Println("error while getting source cached at:", err)
		return time.Time{}
	}
	return status.LastRefresh
}

// CacheTTL returns the TTL time
func (s *RoutesStore) CacheTTL() time.Time {
	return s.sources.NextRefresh(sourceID)
}

// Lookup routes transform
func routeToLookupRoute(
	nStore *NeighborsStore,
	source *config.SourceConfig,
	state string,
	route *api.Route,
) *api.LookupRoute {
	// Get neighbor and make route
	neighbor, _ := nStore.GetNeighborAt(source.ID, route.NeighborID)
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

// LookupNeighborsPrefixesAt performs a single route server
// routes lookup by neighbor id
func (s *RoutesStore) LookupNeighborsPrefixesAt(
	sourceID string,
	neighborIDs []string,
) (api.LookupRoutes, error) {
	ctx := context.TODO()
	return s.backend.GetNeighborsPrefixesAt(
		ctx, sourceID, neighborIDs)
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
			s.neighbors,
			cfg,
			routes.Filtered,
			prefix,
			"filtered")
		imported := filterRoutesByPrefix(
			s.neighbors,
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
// a set of neighbors.
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
