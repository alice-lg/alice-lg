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
func getNeighborByID(neighbours api.Neighbors, id string) (*api.Neighbor, error) {
	for _, n := range neighbours {
		if n.ID == id {
			return n, nil
		}
	}
	unknown := &api.Neighbor{
		ID:          "unknown",
		Description: "Unknown neighbor",
	}
	return unknown, fmt.Errorf("neighbor not found")
}

/*
LockMap uses the sync.Map to manage locks, accessed by a key.
TODO: Maybe this would be a nice generic helper
*/
type LockMap struct {
	locks *sync.Map
}

// NewLockMap creates a new LockMap
func NewLockMap() *LockMap {
	return &LockMap{
		locks: &sync.Map{},
	}
}

// Lock locks the lock.
func (m *LockMap) Lock(key string) {
	mutex, _ := m.locks.LoadOrStore(key, &sync.Mutex{})
	mutex.(*sync.Mutex).Lock()
}

// Unlock unlocks the locked LockMap-lock.
func (m *LockMap) Unlock(key string) {
	mutex, ok := m.locks.Load(key)
	if !ok {
		return // no lock
	}
	mutex.(*sync.Mutex).Unlock()
}

// Wouldn't we all like to know?
func isProtocolUp(protocol string) bool {
	protocol = strings.ToLower(protocol)
	return protocol == "up"
}
