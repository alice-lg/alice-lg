package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
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
			ID:      "r1.2.3.4",
			Network: "1.2.3.0/24",
		},
	}
	b.initTable(ctx, tx, "rs1")
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.ID = "r4242"
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.ID = "r4243"
	r.State = "imported"
	b.persist(ctx, tx, "rs1", r, now)

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
			ID:      "r1.2.3.4",
			Network: "1.2.3.0/24",
		},
	}
	b.initTable(ctx, tx, "rs1")
	b.initTable(ctx, tx, "rs2")
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.ID = "r4242"
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.ID = "r4243"
	r.Neighbor.ID = "n24"
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.ID = "r4244"
	r.Neighbor.ID = "n25"
	b.persist(ctx, tx, "rs2", r, now)

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	routes, err := b.FindByNeighbors(ctx, []string{
		"n24", "n25",
	})
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
			ID:      "r1.2.3.4",
			Network: "1.2.3.0/24",
		},
	}

	b.initTable(ctx, tx, "rs1")
	b.initTable(ctx, tx, "rs2")
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.ID = "r4242"
	r.Route.Network = "1.2.4.0/24"
	b.persist(ctx, tx, "rs1", r, now)

	r.Route.ID = "r4243"
	r.Route.Network = "1.2.5.0/24"
	r.Neighbor.ID = "n24"
	b.persist(ctx, tx, "rs2", r, now)

	r.Route.ID = "r4244"
	r.Route.Network = "5.5.5.0/24"
	r.Neighbor.ID = "n25"
	b.persist(ctx, tx, "rs1", r, now)

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	routes, err := b.FindByPrefix(ctx, "1.2.")
	if err != nil {
		t.Fatal(err)
	}

	if len(routes) != 3 {
		t.Error("unexpected routes:", routes)
	}

	routes, _ = b.FindByPrefix(ctx, "5.5.")
	t.Log(routes)
}
