package caches

import (
	"github.com/alice-lg/alice-lg/pkg/api"
)

/*
The birdwatcher already caches the responses from
bird and provides the API consumers with information
on how long the information is valid.

However, to avoid unnecessary network requests to the
birdwatcher, we keep a local cache. (This comes in handy
when we are paginating the results for better client performance.)
*/

// NeighborsCache implements a cache to store neighbors
type NeighborsCache struct {
	response *api.NeighborsResponse
	disabled bool
}

// NewNeighborsCache initializes a cache for neighbor responses.
func NewNeighborsCache(disabled bool) *NeighborsCache {
	cache := &NeighborsCache{
		response: nil,
		disabled: disabled,
	}

	return cache
}

// Get retrieves the neighbors response from the cache, if present,
// and makes sure the information is still up to date.
func (cache *NeighborsCache) Get() *api.NeighborsResponse {
	if cache.disabled {
		return nil
	}

	if cache.response == nil {
		return nil
	}

	if cache.response.CacheTTL() < 0 {
		return nil
	}

	return cache.response
}

// Set updates the neighbors cache with a new response retrieved
// from a backend source.
func (cache *NeighborsCache) Set(response *api.NeighborsResponse) {
	if cache.disabled {
		return
	}
	cache.response = response
}
