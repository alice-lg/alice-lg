package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// RoutesBackend implements an in memory backend
// for the routes store.
type RoutesBackend struct {
	routes map[string]api.LookupRoutes
	sync.Mutex
}

// NewRoutesBackend creates a new instance
func NewRoutesBackend() *RoutesBackend {
	return &RoutesBackend{
		routes: make(map[string]api.LookupRoutes),
	}
}

// SetRoutes implements the RoutesStoreBackend interface
// function for setting all routes of a source identified
// by ID.
func (r *RoutesBackend) SetRoutes(
	ctx context.Context,
	sourceID string,
	routes api.LookupRoutes,
) error {
	r.Lock()
	defer r.Unlock()

	// Remove details from routes
	for _, r := range routes {
		r.Route.Details = nil
	}

	r.routes[sourceID] = routes
	return nil
}

// CountRoutesAt returns the number of filtered and imported
// routes and implements the RoutesStoreBackend interface.
func (r *RoutesBackend) CountRoutesAt(
	ctx context.Context,
	sourceID string,
) (uint, uint, error) {
	r.Lock()
	defer r.Unlock()
	routes, ok := r.routes[sourceID]
	if !ok {
		return 0, 0, sources.ErrSourceNotFound
	}

	var (
		imported uint = 0
		filtered uint = 0
	)

	for _, route := range routes {
		if route.State == api.RouteStateFiltered {
			filtered++
		}
		if route.State == api.RouteStateImported {
			imported++
		}
	}

	return imported, filtered, nil
}

// FindByNeighbors will return the prefixes for a
// list of neighbors identified by ID.
func (r *RoutesBackend) FindByNeighbors(
	ctx context.Context,
	neighborIDs []string,
) (api.LookupRoutes, error) {
	r.Lock()
	defer r.Unlock()

	result := api.LookupRoutes{}
	for _, rs := range r.routes {
		for _, route := range rs {
			if isMemberOf(neighborIDs, route.NeighborID) {
				result = append(result, route)
			}
		}
	}
	return result, nil
}

// FindByPrefix will return the prefixes matching a pattern
func (r *RoutesBackend) FindByPrefix(
	ctx context.Context,
	prefix string,
) (api.LookupRoutes, error) {
	r.Lock()
	defer r.Unlock()

	// We make our compare case insensitive
	prefix = strings.ToLower(prefix)

	result := api.LookupRoutes{}
	for _, rs := range r.routes {
		for _, route := range rs {
			// Naiive string filtering:
			if strings.HasPrefix(strings.ToLower(route.Network), prefix) {
				result = append(result, route)
			}
		}
	}
	return result, nil
}

// Routes filter
/*
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
		if isMemberOf(neighborIDs, route.NeighborID) {
			lookup := routeToLookupRoute(nStore, source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}
*/

// isMemberOf checks if a key is present in
// a list of strings.
func isMemberOf(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}
