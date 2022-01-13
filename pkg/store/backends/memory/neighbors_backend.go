package memory

import (
	"context"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// NeighborIndex is a mapping of a neighborID to a Neighbor
type NeighborIndex map[string]*api.Neighbor

// NeighborsMap associate a route server source (ID) with
// a neighbor index.
type NeighborsMap map[string]NeighborIndex

// NeighborsBackend is an in-memory backend implementation for
// the neighbors store.
type NeighborsBackend struct {
	neighbors NeighborsMap
	sync.Mutex
}

// NewNeighborsBackend instanciates a new in memory
// neighbors backend.
func NewNeighborsBackend() *NeighborsBackend {
	return &NeighborsBackend{
		neighbors: make(NeighborsMap),
	}
}

// SetNeighbors replaces the current list of neighbors
// for a given routeserver source.
func (b *NeighborsBackend) SetNeighbors(
	ctx context.Context,
	sourceID string,
	neighbors api.Neighbors,
) error {
	b.Lock()
	defer b.Unlock()

	// Make index
	idx := make(NeighborIndex)
	for _, neighbor := range neighbors {
		idx[neighbor.ID] = neighbor
	}

	b.neighbors[sourceID] = idx
	return nil
}

// GetNeighborsAt retrieves all neighbors for a source
// identified by its ID.
func (b *NeighborsBackend) GetNeighborsAt(
	ctx context.Context,
	sourceID string,
) (api.Neighbors, error) {
	b.Lock()
	defer b.Unlock()

	neighbors, ok := b.neighbors[sourceID]
	if !ok {
		return nil, sources.ErrSourceNotFound
	}

	ret := make(api.Neighbors, 0, len(neighbors))
	for _, neighbor := range neighbors {
		ret = append(ret, neighbor)
	}
	return ret, nil
}

// GetNeighborsMapAt retrieves all neighbors for a source
// identified by its ID and returns a map.
func (b *NeighborsBackend) GetNeighborsMapAt(
	ctx context.Context,
	sourceID string,
) (map[string]*api.Neighbor, error) {
	b.Lock()
	defer b.Unlock()

	neighbors, ok := b.neighbors[sourceID]
	if !ok {
		return nil, sources.ErrSourceNotFound
	}

	// Copy neighbors map
	result := make(map[string]*api.Neighbor)
	for k, v := range neighbors {
		result[k] = v
	}
	return result, nil
}

// CountNeighborsAt retrievs the number of neighbors
// at this source.
func (b *NeighborsBackend) CountNeighborsAt(
	ctx context.Context,
	sourceID string,
) (int, error) {
	neighbors, err := b.GetNeighborsAt(ctx, sourceID)
	if err != nil {
		return 0, err
	}
	return len(neighbors), nil
}
