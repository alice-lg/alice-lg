package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/store/testdata"
)

func TestConcurrentRoutesAccess(t *testing.T) {
	ctx := context.Background()

	t0 := time.Now()
	wg := sync.WaitGroup{}

	rs1 := testdata.LoadTestLookupRoutes("rs1", "routeserver1")
	rs2 := testdata.LoadTestLookupRoutes("rs2", "routeserver2")

	b := NewRoutesBackend()
	b.SetRoutes(ctx, "rs1", rs1)
	b.SetRoutes(ctx, "rs2", rs2)

	// Current: ~327 ms, With sync.Map: 80 ms... neat
	for i := 0; i < 200000; i++ {
		wg.Add(1)
		go func() {
			b.FindByNeighbors(ctx, []string{"ID7254_AS31334", "ID163_AS31078"})
			wg.Done()
		}()
	}

	wg.Wait()
	dt := time.Since(t0)
	fmt.Println("finished after:", dt)
}
