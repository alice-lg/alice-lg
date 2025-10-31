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
	routes *sync.Map
}

// NewRoutesBackend creates a new instance
func NewRoutesBackend() *RoutesBackend {
	return &RoutesBackend{
		routes: &sync.Map{},
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
	r.routes.Store(sourceID, routes)
	return nil
}

// CountRoutesAt returns the number of filtered and imported
// routes and implements the RoutesStoreBackend interface.
func (r *RoutesBackend) CountRoutesAt(
	ctx context.Context,
	sourceID string,
) (uint, uint, error) {
	routes, ok := r.routes.Load(sourceID)
	if !ok {
		return 0, 0, sources.ErrSourceNotFound
	}

	var (
		imported uint = 0
		filtered uint = 0
	)

	for _, route := range routes.(api.LookupRoutes) {
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
	query []*api.NeighborQuery,
	filters *api.SearchFilters,
) (api.LookupRoutes, error) {
	result := api.LookupRoutes{}

	r.routes.Range(func(k, rs any) bool {
		for _, route := range rs.(api.LookupRoutes) {
			for _, q := range query {
				if !route.MatchNeighborQuery(q) {
					continue
				}
				if !filters.MatchRoute(route) {
					continue
				}

				result = append(result, route)
			}
		}
		return true
	})

	return result, nil
}

// FindByPrefix will return the prefixes matching a pattern
func (r *RoutesBackend) FindByPrefix(
	ctx context.Context,
	prefix string,
	filters *api.SearchFilters,
	limit uint,
) (api.LookupRoutes, error) {
	// We make our compare case insensitive
	var (
		count         uint
		limitExceeded bool
	)

	prefix = strings.ToLower(prefix)
	result := api.LookupRoutes{}
	hasPrefix := prefix != ""
	r.routes.Range(func(k, rs any) bool {
		if limit > 0 && count >= limit {
			limitExceeded = true
			return false
		}
		for _, route := range rs.(api.LookupRoutes) {
			// Naive string filtering:
			if hasPrefix && !strings.HasPrefix(strings.ToLower(route.Network), prefix) {
				continue
			}
			if !filters.MatchRoute(route) {
				continue
			}
			result = append(result, route)
			count++
			if limit > 0 && count >= limit {
				limitExceeded = true
				return false
			}
		}
		return true
	})
	if limitExceeded {
		return nil, api.ErrTooManyRoutes
	}
	return result, nil
}
