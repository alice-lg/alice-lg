package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/pools"
	"github.com/alice-lg/alice-lg/pkg/store/testdata"
)

func TestFindByNeighbors(t *testing.T) {
	ctx := context.Background()

	rs1 := testdata.LoadTestLookupRoutes("rs1", "routeserver1")
	rs2 := testdata.LoadTestLookupRoutes("rs2", "routeserver2")

	b := NewRoutesBackend()
	b.SetRoutes(ctx, "rs1", rs1)
	b.SetRoutes(ctx, "rs2", rs2)

	q := &api.NeighborQuery{
		NeighborID: pools.Neighbors.Get("ID7254_AS31334"),
		SourceID:   pools.RouteServers.Get("rs1"),
	}

	routes, err := b.FindByNeighbors(
		ctx,
		[]*api.NeighborQuery{q},
		api.NewSearchFilters())
	if err != nil {
		t.Fatal(err)
	}

	if len(routes) != 1 {
		t.Error("Route lookup returned unexpected length", len(routes))
	}

	route := routes[0]
	if *route.NeighborID != "ID7254_AS31334" {
		t.Error("Route lookup has wrong neighbor ID")
	}
}

func TestConcurrentRoutesAccess(t *testing.T) {
	ctx := context.Background()

	t0 := time.Now()
	wg := sync.WaitGroup{}

	rs1 := testdata.LoadTestLookupRoutes("rs1", "routeserver1")
	rs2 := testdata.LoadTestLookupRoutes("rs2", "routeserver2")

	b := NewRoutesBackend()
	b.SetRoutes(ctx, "rs1", rs1)
	b.SetRoutes(ctx, "rs2", rs2)

	n1 := &api.NeighborQuery{
		NeighborID: pools.Neighbors.Get("ID7254_AS31334"),
		SourceID:   pools.RouteServers.Get("rs1"),
	}
	n2 := &api.NeighborQuery{
		NeighborID: pools.Neighbors.Get("ID163_AS31078"),
		SourceID:   pools.RouteServers.Get("rs2"),
	}

	// Current: ~327 ms, With sync.Map: 80 ms... neat
	for range 200000 {
		wg.Add(1)
		go func() {
			b.FindByNeighbors(ctx, []*api.NeighborQuery{n1, n2}, api.NewSearchFilters())
			wg.Done()
		}()
	}

	wg.Wait()
	dt := time.Since(t0)
	fmt.Println("finished after:", dt)
}
