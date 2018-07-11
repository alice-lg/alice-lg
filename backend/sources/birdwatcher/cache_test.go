package birdwatcher

import (
	"github.com/alice-lg/alice-lg/backend/api"

	"testing"
	"time"
)

/*
NeighborsCache Tests
*/

func TestNeighborsCacheSetGet(t *testing.T) {
	cache := NewNeighborsCache(false)

	response := &api.NeighboursResponse{
		Api: api.ApiStatus{
			Ttl: time.Now().UTC().Add(300 * time.Millisecond),
		},
	}

	if cache.Get() != nil {
		t.Error("There should not be anything cached yet!")
	}

	cache.Set(response)

	fromCache := cache.Get()
	if fromCache != response {
		t.Error("Expected", response, "got", fromCache)
	}

	// Wait a bit
	time.Sleep(333 * time.Millisecond)

	fromCache = cache.Get()
	if fromCache != nil {
		t.Error("Expected empty cache result, got:", fromCache)
	}
}
