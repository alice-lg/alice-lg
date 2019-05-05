package caches

import (
	"github.com/alice-lg/alice-lg/backend/api"
)

/*
The birdwatcher already caches the responses from
bird and provides the API consumers with information
on how long the information is valid.

However, to avoid unnecessary network requests to the
birdwatcher, we keep a local cache. (This comes in handy
when we are paginating the results for better client performance.)
*/

type NeighborsCache struct {
	response *api.NeighboursResponse
	disabled bool
}

func NewNeighborsCache(disabled bool) *NeighborsCache {
	cache := &NeighborsCache{
		response: nil,
		disabled: disabled,
	}

	return cache
}

func (self *NeighborsCache) Get() *api.NeighboursResponse {
	if self.disabled {
		return nil
	}

	if self.response == nil {
		return nil
	}

	if self.response.CacheTtl() < 0 {
		return nil
	}

	return self.response
}

func (self *NeighborsCache) Set(response *api.NeighboursResponse) {
	if self.disabled {
		return
	}

	self.response = response
}
