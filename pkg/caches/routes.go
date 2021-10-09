package caches

import (
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
)

/*
RoutesCache stores routes responses from the backend.

Keep a kv map with neighborId <-> api.RoutesResponse
TTL is derived from the api.RoutesResponse.

To avoid memory issues, we only keep N responses (MRU) (per RS).
*/
type RoutesCache struct {
	responses  map[string]*api.RoutesResponse
	accessedAt LRUMap

	disabled bool
	size     int

	sync.Mutex
}

// NewRoutesCache initializes a new cache for route responses.
func NewRoutesCache(disabled bool, size int) *RoutesCache {
	cache := &RoutesCache{
		responses:  make(map[string]*api.RoutesResponse),
		accessedAt: make(map[string]time.Time),
		disabled:   disabled,
		size:       size,
	}

	return cache
}

// Get retrievs all routes for a given neighbor
func (cache *RoutesCache) Get(neighborID string) *api.RoutesResponse {
	if cache.disabled {
		return nil
	}

	cache.Lock()
	defer cache.Unlock()

	response, ok := cache.responses[neighborID]
	if !ok {
		return nil
	}

	if response.CacheTTL() < 0 {
		return nil
	}

	cache.accessedAt[neighborID] = time.Now()

	return response
}

// Set the routes response for a given neighbor
func (cache *RoutesCache) Set(neighborID string, response *api.RoutesResponse) {
	if cache.disabled {
		return
	}

	cache.Lock()
	defer cache.Unlock()

	if len(cache.responses) > cache.size {
		// delete LRU
		leastRecentNeighbor := cache.accessedAt.LRU()
		delete(cache.accessedAt, leastRecentNeighbor)
		delete(cache.responses, leastRecentNeighbor)
	}

	cache.accessedAt[neighborID] = time.Now()
	cache.responses[neighborID] = response
}

// Expire will flush expired keys. (TODO: naming could be better.)
func (cache *RoutesCache) Expire() int {
	cache.Lock()
	defer cache.Unlock()

	expiredKeys := []string{}
	for key, response := range cache.responses {
		if response.CacheTTL() < 0 {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(cache.responses, key)
	}

	return len(expiredKeys)
}
