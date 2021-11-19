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
	sync.RWMutex
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
		index[neighbor.ID] = neighbor
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

	ret := make(api.Neighbors, 0, len(b.neighbors))
	for _, neighbor := range b.neighbors {
		ret = append(ret, neighbor)
	}
	return ret, nil
}

// GetNeighborAt retrieves all neighbors for a source
// identified by its ID.
func (b *NeighborsBackend) GetNeighborAt(
	ctx context.Context,
	sourceID string,
) (api.Neighbors, error) {
	b.Lock()
	defer b.Unlock()

	neighbors, ok := b.neighbors[sourceID]
	if !ok {
		return nil, sources.ErrSourceNotFound
	}

	return ret, nil
}
