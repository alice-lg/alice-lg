package birdwatcher

import (
	"strings"
	"sync"
)

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
