package birdwatcher

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"log"
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

	log.Println("NEW NN CACHE!!")

	return cache
}

func (self *NeighborsCache) Get() *api.NeighboursResponse {
	if self.disabled {
		return nil
	}

	if self.response == nil {
		return nil
	}

	log.Println("Response present, check ttl:")
	log.Println(self.response.CacheTtl())

	if self.response.CacheTtl() < 0 {
		return nil
	}

	return self.response
}

func (self *NeighborsCache) Set(response *api.NeighboursResponse) {
	self.response = response
}
