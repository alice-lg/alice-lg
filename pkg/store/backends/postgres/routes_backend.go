package postgres

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	// ReMatchNonChar will match on all characters not
	// within a to Z or 0 to 9
	ReMatchNonChar = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

// RoutesBackend implements a postgres store for routes.
type RoutesBackend struct {
	pool    *pgxpool.Pool
	sources []*config.SourceConfig
}

// NewRoutesBackend creates a new instance with a postgres
// connection pool.
func NewRoutesBackend(
	pool *pgxpool.Pool,
	sources []*config.SourceConfig,
) *RoutesBackend {
	return &RoutesBackend{
		pool:    pool,
		sources: sources,
	}
}

// Init will initialize all the route tables
func (b *RoutesBackend) Init(ctx context.Context) error {
	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, src := range b.sources {
		if err := b.initTable(ctx, tx, src.ID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
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

	// Create table from template
	if err := b.initTable(ctx, tx, sourceID); err != nil {
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

// Private routesTable returns the name of the routes table
// for a sourceID
func (b *RoutesBackend) routesTable(sourceID string) string {
	sourceID = ReMatchNonChar.ReplaceAllString(sourceID, "_")
	return "routes_" + sourceID
}

// Private initTable recreates the routes table
// for a single sourceID
func (b *RoutesBackend) initTable(
	ctx context.Context,
	tx pgx.Tx,
	sourceID string,
) error {
	tbl := b.routesTable(sourceID)
	qry := `
		DROP TABLE IF EXISTS ` + tbl + `;
		CREATE TABLE ` + tbl + ` ( LIKE routes INCLUDING ALL )
	`
	_, err := tx.Exec(ctx, qry)
	return err
}

// Private persist route in database
func (b *RoutesBackend) persist(
	ctx context.Context,
	tx pgx.Tx,
	sourceID string,
	route *api.LookupRoute,
	now time.Time,
) error {
	tbl := b.routesTable(sourceID)
	qry := `
		INSERT INTO ` + tbl + ` (
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
		route.Route.Network,
		sourceID,
		route.Neighbor.ID,
		route.Route.Network,
		route,
		now)
	return err
}

// Private clear removes all routes.
/*
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
*/

// Private queryCountByState will query routes and filter
// by state
func (b *RoutesBackend) queryCountByState(
	ctx context.Context,
	tx pgx.Tx,
	sourceID string,
	state string,
) pgx.Row {
	tbl := b.routesTable(sourceID)
	qry := `SELECT COUNT(1) FROM ` + tbl + ` 
			 WHERE route -> 'state' = $1`

	return tx.QueryRow(ctx, qry, "\""+state+"\"")
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
	neighbors []*api.NeighborQuery,
	filters *api.SearchFilters,
) (api.LookupRoutes, error) {
	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	vals := make([]interface{}, 0, len(neighbors))
	vars := 0

	qrys := []string{}

	for _, neighborQuery := range neighbors {
		tbl := b.routesTable(*neighborQuery.SourceID)
		param := fmt.Sprintf("$%d", vars+1)
		vals = append(vals, *neighborQuery.NeighborID)

		qry := `
			SELECT route FROM ` + tbl + `
			 WHERE neighbor_id = ` + param
		qrys = append(qrys, qry)

		vars++
	}

	qry := strings.Join(qrys, " UNION ")

	rows, err := tx.Query(ctx, qry, vals...)
	if err != nil {
		return nil, err
	}

	return fetchRoutes(rows, filters, 0)
}

// FindByPrefix will return the prefixes matching a pattern
func (b *RoutesBackend) FindByPrefix(
	ctx context.Context,
	prefix string,
	filters *api.SearchFilters,
	limit uint,
) (api.LookupRoutes, error) {
	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// We are searching route.Network
	qrys := []string{}
	for _, src := range b.sources {
		tbl := b.routesTable(src.ID)
		qry := `
			SELECT route FROM ` + tbl + `
			 WHERE network ILIKE $1
		`
		qrys = append(qrys, qry)
	}
	qry := strings.Join(qrys, " UNION ")
	rows, err := tx.Query(ctx, qry, prefix+"%")
	if err != nil {
		return nil, err
	}
	return fetchRoutes(rows, filters, limit)
}

// Private fetchRoutes will load the queried result set
func fetchRoutes(
	rows pgx.Rows,
	filters *api.SearchFilters,
	limit uint,
) (api.LookupRoutes, error) {
	var count uint
	cmd := rows.CommandTag()
	results := make(api.LookupRoutes, 0, cmd.RowsAffected())
	for rows.Next() {
		route := &api.LookupRoute{}
		if err := rows.Scan(&route); err != nil {
			return nil, err
		}
		if !filters.MatchRoute(route) {
			continue
		}
		results = append(results, route)
		count++
		if limit > 0 && count >= limit {
			return nil, api.ErrTooManyRoutes
		}
	}
	return results, nil
}
