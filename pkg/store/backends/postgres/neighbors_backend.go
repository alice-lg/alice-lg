package postgres

import (
	"context"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"

	"github.com/jackc/pgx/v4"
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
	// Clear current neighbors
	now := time.Now().UTC()
	for _, n := range neighbors {
		if err := b.persist(ctx, sourceID, n, now); err != nil {
			return err
		}
	}
	// Remove old neighbors
	if err := b.deleteStale(ctx, sourceID, now); err != nil {
		return err
	}
	return nil
}

// Private persist saves a neighbor to the database
func (b *NeighborsBackend) persist(
	ctx context.Context,
	sourceID string,
	neighbor *api.Neighbor,
	now time.Time,
) error {
	qry := `
	  INSERT INTO neighbors (
	  		id, rs_id, neighbor, updated_at
		) VALUES ( $1, $2, $3, $4 )
	  ON CONFLICT ON CONSTRAINT neighbors_pkey DO UPDATE
       SET neighbor   = EXCLUDED.neighbor,
		   updated_at = EXCLUDED.updated_at
	`
	_, err := b.pool.Exec(ctx, qry, neighbor.ID, sourceID, neighbor, now)
	return err
}

// Private deleteStale removes all neighbors not inserted or
// updated at a specific time.
func (b *NeighborsBackend) deleteStale(
	ctx context.Context,
	sourceID string,
	t time.Time,
) error {
	qry := `
	  DELETE FROM neighbors
	        WHERE rs_id = $1 
			  AND updated_at <> $2
	`
	_, err := b.pool.Exec(ctx, qry, sourceID, t)
	return err
}

// Private queryNeighborsAt selects all neighbors
// for a given sourceID
func (b *NeighborsBackend) queryNeighborsAt(
	ctx context.Context,
	sourceID string,
) (pgx.Rows, error) {
	qry := `
		SELECT neighbor
		  FROM neighbors
		 WHERE rs_id = $1
	`
	return b.pool.Query(ctx, qry, sourceID)
}

// GetNeighborsAt retrieves all neighbors associated
// with a route server (source).
func (b *NeighborsBackend) GetNeighborsAt(
	ctx context.Context,
	sourceID string,
) (api.Neighbors, error) {
	rows, err := b.queryNeighborsAt(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	cmd := rows.CommandTag()
	results := make(api.Neighbors, 0, cmd.RowsAffected())
	for rows.Next() {
		neighbor := &api.Neighbor{}
		if err := rows.Scan(&neighbor); err != nil {
			return nil, err
		}
		results = append(results, neighbor)
	}
	return results, nil
}

// GetNeighborsMapAt retrieve a neighbor map for a route server.
// The Neighbor is identified by ID.
func (b *NeighborsBackend) GetNeighborsMapAt(
	ctx context.Context,
	sourceID string,
) (map[string]*api.Neighbor, error) {
	rows, err := b.queryNeighborsAt(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	results := make(map[string]*api.Neighbor)
	for rows.Next() {
		neighbor := &api.Neighbor{}
		if err := rows.Scan(&neighbor); err != nil {
			return nil, err
		}
		results[neighbor.ID] = neighbor
	}
	return results, nil
}

// CountNeighborsAt retrieves the current number of
// stored neighbors.
func (b *NeighborsBackend) CountNeighborsAt(
	ctx context.Context,
	sourceID string,
) (int, error) {
	qry := `
		SELECT COUNT(1) FROM neighbors
		 WHERE rs_id = $1
	`
	count := 0
	err := b.pool.QueryRow(ctx, qry, sourceID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
