package caches

import (
	"github.com/alice-lg/alice-lg/pkg/api"

	"fmt"
	"testing"
	"time"
)

func TestRoutesCacheSetGet(t *testing.T) {
	cache := NewRoutesCache(false, 2)

	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: &api.Meta{
				TTL: time.Now().UTC().Add(23 * time.Millisecond),
			},
		},
	}

	nID := "neighbor_42"

	if cache.Get(nID) != nil {
		t.Error("There should not be anything cached yet!")
	}

	cache.Set(nID, response)

	fromCache := cache.Get(nID)
	if fromCache != response {
		t.Error("Expected", response, "got", fromCache)
	}

	time.Sleep(33 * time.Millisecond)

	fromCache = cache.Get(nID)
	if fromCache != nil {
		t.Error("Expected empty cache result, got:", fromCache)
	}
}

func TestRoutesCacheLru(t *testing.T) {
	cache := NewRoutesCache(false, 2)

	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: &api.Meta{
				TTL: time.Now().UTC().Add(time.Hour),
			},
		},
	}

	cache.Set("n1", response)
	cache.Set("n2", response)
	cache.Set("n3", response)
	cache.Set("n2", response)

	// n1 should be removed as last used
	if len(cache.responses) != 2 {
		t.Error("There should not be more than 2 responses. Got:",
			len(cache.responses),
		)
	}

	_, ok := cache.responses["n1"]
	if ok {
		t.Error("n1 should not be part of the key set")
	}

	// MRU is now n2, LRU: n3, let's access n3 and set n1 again
	if cache.accessedAt.LRU() != "n3" {
		t.Log("Expected n3 to be LRU")
	}
	cache.Get("n3")
	cache.Set("n1", response)

	// n2 should not be part of the key set
	_, ok = cache.responses["n1"]
	if !ok {
		t.Error("n1 should be part of the key set")
	}

	_, ok = cache.responses["n3"]
	if !ok {
		t.Error("n3 should be part of the key set")
	}

	_, ok = cache.responses["n2"]
	if ok {
		t.Error("n2 should NOT be part of the key set")
	}
}

func routesResponseWithTTL(ttl time.Duration) *api.RoutesResponse {
	return &api.RoutesResponse{
		Response: api.Response{
			Meta: &api.Meta{
				TTL: time.Now().UTC().Add(ttl),
			},
		},
	}
}

func TestRoutesCacheSetHoldsSize(t *testing.T) {
	const size = 4
	cache := NewRoutesCache(false, size)

	for i := 0; i < 3*size; i++ {
		cache.Set(fmt.Sprintf("n%d", i), routesResponseWithTTL(time.Hour))
	}

	if len(cache.responses) != size {
		t.Errorf("responses: got %d, want %d", len(cache.responses), size)
	}
	if len(cache.accessedAt) != size {
		t.Errorf("accessedAt: got %d, want %d", len(cache.accessedAt), size)
	}
}

func TestRoutesCacheGetDropsExpiredEntry(t *testing.T) {
	cache := NewRoutesCache(false, 4)
	cache.Set("expired", routesResponseWithTTL(-time.Minute))

	if cache.Get("expired") != nil {
		t.Error("Get should return nil for an expired entry")
	}
	if _, ok := cache.responses["expired"]; ok {
		t.Error("expired entry should be gone from responses")
	}
	if _, ok := cache.accessedAt["expired"]; ok {
		t.Error("expired entry should be gone from accessedAt")
	}
}

func TestRoutesCacheExpireThenSetHoldsSize(t *testing.T) {
	const size = 4
	cache := NewRoutesCache(false, size)

	for i := 0; i < size; i++ {
		cache.Set(fmt.Sprintf("expired-%d", i), routesResponseWithTTL(-time.Minute))
	}
	if expired := cache.Expire(); expired != size {
		t.Errorf("Expire: got %d, want %d", expired, size)
	}
	if len(cache.accessedAt) != 0 {
		t.Errorf("accessedAt after Expire: got %d, want 0", len(cache.accessedAt))
	}

	for i := 0; i < 3*size; i++ {
		cache.Set(fmt.Sprintf("fresh-%d", i), routesResponseWithTTL(time.Hour))
	}
	if len(cache.responses) != size {
		t.Errorf("responses: got %d, want %d", len(cache.responses), size)
	}
	if len(cache.accessedAt) != size {
		t.Errorf("accessedAt: got %d, want %d", len(cache.accessedAt), size)
	}
}

func TestRoutesCacheSetExistingKeyKeepsOthers(t *testing.T) {
	const size = 3
	cache := NewRoutesCache(false, size)

	for i := 0; i < size; i++ {
		cache.Set(fmt.Sprintf("n%d", i), routesResponseWithTTL(time.Hour))
	}
	cache.Set("n1", routesResponseWithTTL(time.Hour))

	if len(cache.responses) != size {
		t.Errorf("responses: got %d, want %d", len(cache.responses), size)
	}
	if _, ok := cache.responses["n0"]; !ok {
		t.Error("n0 should survive a re-Set of n1")
	}
}
