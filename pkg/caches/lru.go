package caches

import (
	"time"
)

/*
LRUMap is a cache map which uses
a least recently used caching strategy:
Store last access in map, retrieve least recently
used key.
*/
type LRUMap map[string]time.Time

// LRU retrievs the least recently used key
func (lrumap LRUMap) LRU() string {
	t := time.Now()
	key := ""

	for k, v := range lrumap {
		if v.Before(t) {
			t = v
			key = k
		}
	}

	return key
}
