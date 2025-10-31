//go:build pgstore

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestPersistNeighborLookup(t *testing.T) {
	ctx := context.Background()
	pool := ConnectTest()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	b := &NeighborsBackend{pool: pool}
	n := &api.Neighbor{
		ID:      "n2342",
		Address: "test123",
	}

	now := time.Now().UTC()
	if err := b.clear(ctx, tx, "rs1"); err != nil {
		t.Fatal(err)
	}

	if err := b.persist(ctx, tx, "rs1", n, now); err != nil {
		t.Fatal(err)
	}

	// Add a second
	n.ID = "foo23"
	if err := b.persist(ctx, tx, "rs1", n, now); err != nil {
		t.Fatal(err)
	}

	// Add to different rs
	if err := b.persist(ctx, tx, "rs2", n, now); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	neighbors, err := b.GetNeighborsMapAt(ctx, "rs1")
	if err != nil {
		t.Fatal(err)
	}
	if neighbors["n2342"].Address != "test123" {
		t.Error("unexpected neighbors:", neighbors)
	}

	list, err := b.GetNeighborsAt(ctx, "rs2")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Error("unexpected neighbors list:", list)
	}
}

func TestSetNeighbors(t *testing.T) {
	ctx := context.Background()
	pool := ConnectTest()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	b := &NeighborsBackend{pool: pool}

	// Persist an old neighbor, should be gone because stale
	n := &api.Neighbor{
		ID:      "n1",
		Address: "foo",
	}
	b.persist(ctx, tx, "rs1", n, time.Time{})

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	result, _ := b.GetNeighborsAt(ctx, "rs1")
	if len(result) != 1 {
		t.Fatal("unexpected neighbors:", result)
	}

	neighbors := api.Neighbors{
		{
			ID:      "n2342",
			Address: "test123",
		},
		{
			ID:      "n2343",
			Address: "test124",
		},
		{
			ID:      "n2345",
			Address: "test125",
		},
	}
	if err := b.SetNeighbors(ctx, "rs1", neighbors); err != nil {
		t.Fatal(err)
	}

	result, err = b.GetNeighborsAt(ctx, "rs1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(neighbors) {
		t.Error("unexpected neighbors:", result)
	}
}
