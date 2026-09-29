package postgres

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
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
	// sourcesMu guards sources, which is mutated by
	// AddSource and RemoveSource during a reload while
	// FindByPrefix and Init read it.
	sourcesMu sync.RWMutex
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

// sourceIDs returns a snapshot of the registered source IDs.
func (b *RoutesBackend) sourceIDs() []string {
	b.sourcesMu.RLock()
	defer b.sourcesMu.RUnlock()
	return b.sourceIDsLocked()
}

// sourceIDsLocked lists the registered source IDs; the caller
// holds sourcesMu.
func (b *RoutesBackend) sourceIDsLocked() []string {
	ids := make([]string, 0, len(b.sources))
	for _, src := range b.sources {
		ids = append(ids, src.ID)
	}
	return ids
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

	for _, id := range b.sourceIDs() {
		if err := b.initTable(ctx, tx, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// AddSource creates the (empty) routes table for a source and
// registers it once the table exists, so lookups succeed before
// the first refresh and never query a missing table.
func (b *RoutesBackend) AddSource(
	ctx context.Context,
	sourceID string,
) error {
	if b.hasSource(sourceID) {
		return nil // already registered
	}

	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := b.initTable(ctx, tx, sourceID); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	b.register(sourceID)
	return nil
}

// register adds a source whose table exists to the lookups.
func (b *RoutesBackend) register(sourceID string) {
	b.sourcesMu.Lock()
	defer b.sourcesMu.Unlock()
	if b.hasSourceLocked(sourceID) {
		return
	}
	// FindByPrefix and Init only read the ID from this slice.
	b.sources = append(b.sources, &config.SourceConfig{ID: sourceID})
}

// hasSource reports whether a source is registered.
func (b *RoutesBackend) hasSource(sourceID string) bool {
	b.sourcesMu.RLock()
	defer b.sourcesMu.RUnlock()
	return b.hasSourceLocked(sourceID)
}

// hasSourceLocked reports whether a source is registered; the
// caller holds sourcesMu.
func (b *RoutesBackend) hasSourceLocked(sourceID string) bool {
	for _, src := range b.sources {
		if src.ID == sourceID {
			return true
		}
	}
	return false
}

// RemoveSource deregisters a source and drops its routes table.
func (b *RoutesBackend) RemoveSource(
	ctx context.Context,
	sourceID string,
) error {
	// Held across the drop, so no lookup is querying the table.
	b.sourcesMu.Lock()
	defer b.sourcesMu.Unlock()
	kept := make([]*config.SourceConfig, 0, len(b.sources))
	for _, src := range b.sources {
		if src.ID != sourceID {
			kept = append(kept, src)
		}
	}
	b.sources = kept

	tbl := b.routesTable(sourceID)
	_, err := b.pool.Exec(ctx, "DROP TABLE IF EXISTS "+tbl)
	return err
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
	// SetRoutes recreates the table, so a source whose AddSource
	// failed becomes visible to lookups on its next refresh.
	b.register(sourceID)
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

	// Held for the query, so RemoveSource cannot drop a table
	// this query reads.
	b.sourcesMu.RLock()
	defer b.sourcesMu.RUnlock()

	vals := make([]any, 0, len(neighbors))
	vars := 0

	qrys := []string{}

	for _, neighborQuery := range neighbors {
		if !b.hasSourceLocked(*neighborQuery.SourceID) {
			continue // the source was removed, so its table is gone
		}
		tbl := b.routesTable(*neighborQuery.SourceID)
		param := fmt.Sprintf("$%d", vars+1)
		vals = append(vals, *neighborQuery.NeighborID)

		qry := `
			SELECT route FROM ` + tbl + `
			 WHERE neighbor_id = ` + param
		qrys = append(qrys, qry)

		vars++
	}

	if len(qrys) == 0 {
		return api.LookupRoutes{}, nil
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

	// Held for the query, so RemoveSource cannot drop a table
	// this query reads.
	b.sourcesMu.RLock()
	defer b.sourcesMu.RUnlock()

	// We are searching route.Network
	qrys := []string{}
	for _, id := range b.sourceIDsLocked() {
		tbl := b.routesTable(id)
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
