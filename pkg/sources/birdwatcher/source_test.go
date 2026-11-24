package birdwatcher

import (
	"fmt"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/caches"
)

func TestNewBirdwatcherUsesRoutesCacheSize(t *testing.T) {
	const size = 2
	const neighbors = 10 * size

	src := NewBirdwatcher(Config{Type: "multi_table", RoutesCacheSize: size})
	bw := src.(*MultiTableBirdwatcher)

	routesCaches := map[string]*caches.RoutesCache{
		"routesNotExportedCache": bw.routesNotExportedCache,
		"routesRequiredCache":    bw.routesRequiredCache,
	}
	for name, cache := range routesCaches {
		for i := 0; i < neighbors; i++ {
			cache.Set(fmt.Sprintf("n%d", i), &api.RoutesResponse{
				Response: api.Response{
					Meta: &api.Meta{TTL: time.Now().UTC().Add(time.Hour)},
				},
			})
		}

		held := 0
		for i := 0; i < neighbors; i++ {
			if cache.Get(fmt.Sprintf("n%d", i)) != nil {
				held++
			}
		}
		// At most one entry over the configured size.
		if held > size+1 {
			t.Errorf("%s holds %d neighbors, want at most %d for routes_cache_size %d", name, held, size+1, size)
		}
	}
}
