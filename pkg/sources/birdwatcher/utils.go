package birdwatcher

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
)

/*
Helper functions for dealing with birdwatcher API data
*/

// Get neighbour by protocol id
func getNeighbourById(neighbours api.Neighbours, id string) (*api.Neighbour, error) {
	for _, n := range neighbours {
		if n.Id == id {
			return n, nil
		}
	}
	unknown := &api.Neighbour{
		Id:          "unknown",
		Description: "Unknown neighbour",
	}
	return unknown, fmt.Errorf("Neighbour not found")
}

/*
LockMap: Uses the sync.Map to manage locks, accessed by a key.
TODO: Maybe this would be a nice generic helper
*/
type LockMap struct {
	locks *sync.Map
}

func NewLockMap() *LockMap {
	return &LockMap{
		locks: &sync.Map{},
	}
}

func (self *LockMap) Lock(key string) {
	mutex, _ := self.locks.LoadOrStore(key, &sync.Mutex{})
	mutex.(*sync.Mutex).Lock()
}

func (self *LockMap) Unlock(key string) {
	mutex, ok := self.locks.Load(key)
	if !ok {
		return // Nothing to unlock
	}
	mutex.(*sync.Mutex).Unlock()
}

func isProtocolUp(protocol string) bool {
	protocol = strings.ToLower(protocol)
	return protocol == "up"
}
