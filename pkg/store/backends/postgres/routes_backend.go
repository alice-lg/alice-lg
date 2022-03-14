package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// RoutesBackend implements a postgres store for routes.
type RoutesBackend struct {
	pool *pgxpool.Pool
}

// NewRoutesBackend creates a new instance with a postgres
// connection pool.
func NewRoutesBackend(pool *pgxpool.Pool) *RoutesBackend {
	return &RoutesBackend{
		pool: pool,
	}
}

// SetRoutes implements the RoutesStoreBackend interface
// function for setting all routes of a source identified
// by ID.
func (b *RoutesBackend) SetRoutes(
	ctx context.Context,
	sourceID string,
	routes api.LookupRoutes,
) error {
	now := time.Now().UTC()

	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := b.clear(ctx, tx, sourceID); err != nil {
		return err
	}
	// persist all routes
	for _, r := range routes {
		if err := b.persist(ctx, tx, sourceID, r, now); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// Private persist route in database
func (b *RoutesBackend) persist(
	ctx context.Context,
	tx pgx.Tx,
	sourceID string,
	route *api.LookupRoute,
	now time.Time,
) error {
	qry := `
		INSERT INTO routes (
				id,
				rs_id,
				neighbor_id,
				network,
				route,
				updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6
			)
	`
	_, err := tx.Exec(
		ctx,
		qry,
		route.Route.ID,
		sourceID,
		route.Neighbor.ID,
		route.Route.Network,
		route,
		now)
	return err
}

// Private clear removes all routes.
func (b *RoutesBackend) clear(
	ctx context.Context,
	tx pgx.Tx,
	sourceID string,
) error {
	qry := `
	  DELETE FROM routes WHERE rs_id = $1 
	`
	_, err := tx.Exec(ctx, qry, sourceID)
	return err
}

// Private queryCountByState will query routes and filter
// by state
func (b *RoutesBackend) queryCountByState(
	ctx context.Context,
	tx pgx.Tx,
	sourceID string,
	state string,
) pgx.Row {
	qry := `SELECT COUNT(1) FROM routes
			 WHERE rs_id = $1 AND route -> 'state' = $2`

	return tx.QueryRow(ctx, qry, sourceID, "\""+state+"\"")
}

// CountRoutesAt returns the number of filtered and imported
// routes and implements the RoutesStoreBackend interface.
func (b *RoutesBackend) CountRoutesAt(
	ctx context.Context,
	sourceID string,
) (uint, uint, error) {
	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	var (
		imported uint
		filtered uint
	)
	err = b.queryCountByState(ctx, tx, sourceID, api.RouteStateFiltered).
		Scan(&filtered)
	if err != nil {
		return 0, 0, err
	}
	err = b.queryCountByState(ctx, tx, sourceID, api.RouteStateImported).
		Scan(&imported)
	if err != nil {
		return 0, 0, err
	}
	return imported, filtered, nil
}

// FindByNeighbors will return the prefixes for a
// list of neighbors identified by ID.
func (b *RoutesBackend) FindByNeighbors(
	ctx context.Context,
	neighborIDs []string,
) (api.LookupRoutes, error) {
	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	vals := make([]interface{}, len(neighborIDs))
	for i := range neighborIDs {
		vals[i] = neighborIDs[i]
	}
	vars := make([]string, len(neighborIDs))
	for i := range neighborIDs {
		vars[i] = fmt.Sprintf("$%d", i+1)
	}
	listQry := strings.Join(vars, ",")
	qry := `
		SELECT route FROM routes
		 WHERE neighbor_id IN (` + listQry + `)`

	rows, err := tx.Query(ctx, qry, vals...)
	if err != nil {
		return nil, err
	}

	return fetchRoutes(rows)
}

// FindByPrefix will return the prefixes matching a pattern
func (b *RoutesBackend) FindByPrefix(
	ctx context.Context,
	prefix string,
) (api.LookupRoutes, error) {
	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// We are searching route.Network
	qry := `
		SELECT route FROM routes
		 WHERE network ILIKE $1
	`
	rows, err := tx.Query(ctx, qry, prefix+"%")
	if err != nil {
		return nil, err
	}
	return fetchRoutes(rows)
}

// Private fetchRoutes will load the queried result set
func fetchRoutes(rows pgx.Rows) (api.LookupRoutes, error) {
	cmd := rows.CommandTag()
	results := make(api.LookupRoutes, 0, cmd.RowsAffected())
	for rows.Next() {
		route := &api.LookupRoute{}
		if err := rows.Scan(&route); err != nil {
			return nil, err
		}
		results = append(results, route)
	}
	return results, nil
}
