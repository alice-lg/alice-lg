package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// A RoutesMap is a mapping between a source ID and a
// routes response.
type RoutesMap [string]*api.RoutesResponse

// RoutesBackend implements an in memory backend
// for the routes store.
type RoutesBackend struct {
	routes map[string]api.LookupRoutes
	sync.Mutex
}

// NewRoutesBackend creates a new instance
func NewRoutesBackend() *RoutesBackend {
	return &RoutesBackend{
		routes: make(RoutesMap),
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
	r.routesMap[sourceID] = routes
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

	imported := len(routes.Imported)
	filtered := len(routes.Filtered)

	return imported, filtered, nil
}

// GetNeighborsPrefixesAt retrieves the announced
// prefixes of a set of neighbor ids.
func (r *RoutesBackend) GetNeighborsPrefixesAt(
	ctx context.Context,
	sourceID string,
	neighborIDs []string,
) (api.LookupRoutes, error) {
	r.Lock()
	defer r.Unlock()

	result := api.LookupRoutes{}
	routes, ok := s.routes[sourceID]
	if !ok {
		return nil, sources.ErrSourceNotFound
	}

	for _, route := range routes {
		if isMemberOf(route.NeighborID, neighborIDs) {
			result = append(result, route)
		}
	}

	return result, nil
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
		if isMemberOf(route.NeighborID, neighborIDs) {
			lookup := routeToLookupRoute(nStore, source, state, route)
			results = append(results, lookup)
		}
	}
	return results
}

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
