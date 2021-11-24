package memory

import (
	"context"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// A RoutesMap is a mapping between a source ID and a
// routes response.
type RoutesMap [string]*api.RoutesResponse

// RoutesBackend implements an in memory backend
// for the routes store.
type RoutesBackend struct {
	routes map[string]*api.RoutesResponse
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
	routes *api.RoutesResponse,
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
