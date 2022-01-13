package postgres

import (
	"context"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestPersistNeighborLookup(t *testing.T) {
	pool := ConnectTest()
	b := &NeighborsBackend{pool: pool}
	n := &api.Neighbor{
		ID:      "n2342",
		Address: "test123",
	}
	if err := b.persistNeighbor(context.Background(), "rs1", n); err != nil {
		t.Fatal(err)
	}

	// make an update
	n.Address = "test234"
	if err := b.persistNeighbor(context.Background(), "rs1", n); err != nil {
		t.Fatal(err)
	}

	// Add a second
	n.ID = "foo23"
	if err := b.persistNeighbor(context.Background(), "rs1", n); err != nil {
		t.Fatal(err)
	}

	// Add to different rs
	if err := b.persistNeighbor(context.Background(), "rs2", n); err != nil {
		t.Fatal(err)
	}

	neighbors, err := b.GetNeighborsMapAt(context.Background(), "rs1")
	if err != nil {
		t.Fatal(err)
	}
	if neighbors["n2342"].Address != "test234" {
		t.Error("unexpected neighbors:", neighbors)
	}

	list, err := b.GetNeighborsAt(context.Background(), "rs2")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Error("unexpected neighbors list:", list)
	}
}
