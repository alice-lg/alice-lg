package postgres

import (
	"context"

	"github.com/alice-lg/alice-lg/pkg/api"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NeighborsBackend implements a neighbors store
// using a postgres database
type NeighborsBackend struct {
	pool *pgxpool.Pool
}

// NewNeighborsBackend initializes the backend
// with a pool.
func NewNeighborsBackend(pool *pgxpool.Pool) *NeighborsBackend {
	b := &NeighborsBackend{
		pool: pool,
	}
	return b
}

// SetNeighbors updates the current neighbors of
// a route server identified by sourceID
func (b *NeighborsBackend) SetNeighbors(
	ctx context.Context,
	sourceID string,
	neighbors api.Neighbors,
) error {
	return nil
}
