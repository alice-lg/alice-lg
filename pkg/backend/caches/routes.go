package caches

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"sync"
	"time"
)

/*
Routes Cache:
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

func NewRoutesCache(disabled bool, size int) *RoutesCache {
	cache := &RoutesCache{
		responses:  make(map[string]*api.RoutesResponse),
		accessedAt: make(map[string]time.Time),
		disabled:   disabled,
		size:       size,
	}

	return cache
}

func (self *RoutesCache) Get(neighborId string) *api.RoutesResponse {
	if self.disabled {
		return nil
	}

	self.Lock()
	defer self.Unlock()

	response, ok := self.responses[neighborId]
	if !ok {
		return nil
	}

	if response.CacheTtl() < 0 {
		return nil
	}

	self.accessedAt[neighborId] = time.Now()

	return response
}

func (self *RoutesCache) Set(neighborId string, response *api.RoutesResponse) {
	if self.disabled {
		return
	}

	self.Lock()
	defer self.Unlock()

	if len(self.responses) > self.size {
		// delete LRU
		lru := self.accessedAt.LRU()
		delete(self.accessedAt, lru)
		delete(self.responses, lru)
	}

	self.accessedAt[neighborId] = time.Now()
	self.responses[neighborId] = response
}

func (self *RoutesCache) Expire() int {
	self.Lock()
	defer self.Unlock()

	expiredKeys := []string{}
	for key, response := range self.responses {
		if response.CacheTtl() < 0 {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(self.responses, key)
	}

	return len(expiredKeys)
}
