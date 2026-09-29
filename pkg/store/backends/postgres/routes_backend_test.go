//go:build pgstore

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/pools"
)

func TestRoutesTable(t *testing.T) {
	b := &RoutesBackend{}
	tbl := b.routesTable("rs0-example!/;")
	if tbl != "routes_rs0_example___" {
		t.Error("unexpected table:", tbl)
	}
}

func TestCountRoutesAt(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	pool := ConnectTest()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	b := &RoutesBackend{pool: pool}
	r := &api.LookupRoute{
		State: "filtered",
		Neighbor: &api.Neighbor{
			ID: "n23",
		},
		Route: &api.Route{
			Network: "1.2.3.0/24",
		},
	}
	b.initTable(ctx, tx, "rs1")
	if err := b.persist(ctx, tx, "rs1", r, now); err != nil {
		t.Fatal(err)
	}

	r.Route.Network = "1.2.6.1/24"
	if err := b.persist(ctx, tx, "rs1", r, now); err != nil {
		t.Fatal(err)
	}

	r.State = "imported"
	r.Route.Network = "1.2.5.5/24"
	if err := b.persist(ctx, tx, "rs1", r, now); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	imported, filtered, err := b.CountRoutesAt(ctx, "rs1")
	if err != nil {
		t.Fatal(err)
	}
	if imported != 1 {
		t.Error("unexpected imported:", imported)
	}
	if filtered != 2 {
		t.Error("unexpected filtered:", imported)
	}
}

func TestFindByNeighbors(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	pool := ConnectTest()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	b := &RoutesBackend{
		pool: pool,
		sources: []*config.SourceConfig{
			{ID: "rs1"},
			{ID: "rs2"},
		},
	}
	r := &api.LookupRoute{
		State: "filtered",
		Neighbor: &api.Neighbor{
			ID: "n23",
		},
		Route: &api.Route{
			Network:    "1.2.3.0/24",
			NeighborID: pools.Neighbors.Acquire("n23"),
		},
	}
	b.initTable(ctx, tx, "rs1")
	b.initTable(ctx, tx, "rs2")
	b.persist(ctx, tx, "rs1", r, now)

	r.Network = "1.4.5.0/24"
	b.persist(ctx, tx, "rs1", r, now)

	r.Neighbor.ID = "n24"
	b.persist(ctx, tx, "rs1", r, now)

	r.Neighbor.ID = "n25"
	b.persist(ctx, tx, "rs2", r, now)

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	nq1 := &api.NeighborQuery{
		NeighborID: pools.Neighbors.Acquire("n24"),
		SourceID:   pools.RouteServers.Acquire("rs1"),
	}
	nq2 := &api.NeighborQuery{
		NeighborID: pools.Neighbors.Acquire("n25"),
		SourceID:   pools.RouteServers.Acquire("rs2"),
	}

	routes, err := b.FindByNeighbors(
		ctx,
		[]*api.NeighborQuery{nq1, nq2},
		api.NewSearchFilters())
	if err != nil {
		t.Fatal(err)
	}

	if len(routes) != 2 {
		t.Error("unexpected routes:", routes)
	}
	t.Log(routes)
}

func TestFindByPrefix(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	pool := ConnectTest()
	tx, err := ConnectTest().Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	b := &RoutesBackend{
		pool: pool,
		sources: []*config.SourceConfig{
			{ID: "rs1"},
			{ID: "rs2"},
		},
	}
	r := &api.LookupRoute{
		State: "filtered",
		Neighbor: &api.Neighbor{
			ID: "n23",
		},
		Route: &api.Route{
			Network: "1.2.3.0/24",
		},
	}

	b.initTable(ctx, tx, "rs1")
	b.initTable(ctx, tx, "rs2")
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.Network = "1.2.4.0/24"
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.Network = "1.2.5.0/24"
	r.Neighbor.ID = "n24"
	b.persist(ctx, tx, "rs2", r, now)

	r.Route.Network = "5.5.5.0/24"
	r.Neighbor.ID = "n25"
	b.persist(ctx, tx, "rs1", r, now)

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	routes, err := b.FindByPrefix(ctx, "1.2.", api.NewSearchFilters(), 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(routes) != 3 {
		t.Error("unexpected routes:", routes)
	}

	routes, _ = b.FindByPrefix(ctx, "5.5.", api.NewSearchFilters(), 0)
	t.Log(routes)
}

// TestAddSourceRegistersOnlyAfterTableExists proves a failed table
// creation leaves the source unregistered.
func TestAddSourceRegistersOnlyAfterTableExists(t *testing.T) {
	b := NewRoutesBackend(ConnectTest(), nil)

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.AddSource(canceled, "rs_add_failed"); err == nil {
		t.Fatal("expected an error with a canceled context")
	}
	if b.hasSource("rs_add_failed") {
		t.Error("a source whose table was not created should not be registered")
	}

	ctx := context.Background()
	if err := b.AddSource(ctx, "rs_add_ok"); err != nil {
		t.Fatal(err)
	}
	defer b.RemoveSource(ctx, "rs_add_ok")
	if !b.hasSource("rs_add_ok") {
		t.Error("expected rs_add_ok to be registered")
	}
	if _, err := b.FindByPrefix(ctx, "192.0.2.0", api.NewSearchFilters(), 1); err != nil {
		t.Error("lookup over the registered source failed:", err)
	}
}

// TestSetRoutesRegistersSourceAfterFailedAddSource proves a source
// whose AddSource failed becomes visible to lookups on its next
// refresh instead of staying hidden.
func TestSetRoutesRegistersSourceAfterFailedAddSource(t *testing.T) {
	ctx := context.Background()
	b := NewRoutesBackend(ConnectTest(), nil)

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if err := b.AddSource(canceled, "rs_heal"); err == nil {
		t.Fatal("expected an error with a canceled context")
	}
	defer b.RemoveSource(ctx, "rs_heal")

	route := &api.LookupRoute{
		Neighbor: &api.Neighbor{ID: "n23"},
		Route:    &api.Route{Network: "198.51.100.0/24"},
	}
	if err := b.SetRoutes(ctx, "rs_heal", api.LookupRoutes{route}); err != nil {
		t.Fatal(err)
	}
	found, err := b.FindByPrefix(ctx, "198.51.100.", api.NewSearchFilters(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 {
		t.Error("expected the refreshed source in lookups, got", len(found))
	}
}

// TestLookupsDoNotRaceRemoveSource runs lookups while a source is
// added and removed; a lookup must never query a dropped table.
func TestLookupsDoNotRaceRemoveSource(t *testing.T) {
	ctx := context.Background()
	b := NewRoutesBackend(ConnectTest(), nil)
	if err := b.AddSource(ctx, "rs_stable"); err != nil {
		t.Fatal(err)
	}
	defer b.RemoveSource(ctx, "rs_stable")

	neighbor, source := "n23", "rs_churn"
	query := []*api.NeighborQuery{{NeighborID: &neighbor, SourceID: &source}}

	stop := make(chan struct{})
	errs := make(chan error, 2)
	lookup := func(find func() error) {
		for {
			select {
			case <-stop:
				errs <- nil
				return
			default:
			}
			if err := find(); err != nil {
				errs <- err
				return
			}
		}
	}
	go lookup(func() error {
		_, err := b.FindByPrefix(ctx, "198.51.100.", api.NewSearchFilters(), 0)
		return err
	})
	go lookup(func() error {
		_, err := b.FindByNeighbors(ctx, query, api.NewSearchFilters())
		return err
	})

	for i := 0; i < 50; i++ {
		if err := b.AddSource(ctx, "rs_churn"); err != nil {
			t.Fatal(err)
		}
		if err := b.RemoveSource(ctx, "rs_churn"); err != nil {
			t.Fatal(err)
		}
	}
	close(stop)
	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			t.Error("a lookup failed while a source was removed:", err)
		}
	}
}
