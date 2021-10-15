package caches

import (
	"github.com/alice-lg/alice-lg/pkg/api"

	"testing"
	"time"
)

func TestRoutesCacheSetGet(t *testing.T) {
	cache := NewRoutesCache(false, 2)

	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: api.Meta{
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
			Meta: api.Meta{
				TTL: time.Now().UTC().Add(23 * time.Millisecond),
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
	if !ok {
		t.Error("n2 should NOT be part of the key set")
	}
}
