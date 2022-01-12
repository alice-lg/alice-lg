package postgres

import (
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestPersistNeighbor(t *testing.T) {
	pool := ConnectTest()
	b := &NeighborsBackend{pool: pool}
	n := &api.Neighbor{
		ID: "n2342",
	}
	t.Log(n)
	t.Log(b)
}
