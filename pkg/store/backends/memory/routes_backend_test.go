package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/store"
)

func TestConcurrentRoutesAccess(t *testing.T) {
	t0 := time.Now()

	routes := store.LoadTestLookupRoutes()

	dt := time.Since(t0)
	fmt.Println("finished after:", dt)
}
