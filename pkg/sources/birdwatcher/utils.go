package birdwatcher

import (
	"strings"
	"sync"
)

/*
LockMap hands out one lock per key. A key stays in the map
only while a caller holds or waits for its lock.
TODO: Maybe this would be a nice generic helper
*/
type LockMap struct {
	mu    sync.Mutex
	locks map[string]*keyLock
}

// keyLock is the lock for one key and the number of
// callers holding or waiting for it.
type keyLock struct {
	sync.Mutex
	refs int
}

// NewLockMap creates a new LockMap
func NewLockMap() *LockMap {
	return &LockMap{
		locks: make(map[string]*keyLock),
	}
}

// Lock locks the lock.
func (m *LockMap) Lock(key string) {
	m.mu.Lock()
	lock, ok := m.locks[key]
	if !ok {
		lock = &keyLock{}
		m.locks[key] = lock
	}
	lock.refs++
	m.mu.Unlock()

	lock.Lock()
}

// Unlock unlocks the locked LockMap-lock and drops the key
// once no caller holds or waits for it.
func (m *LockMap) Unlock(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	lock, ok := m.locks[key]
	if !ok {
		return // no lock
	}
	lock.refs--
	if lock.refs == 0 {
		delete(m.locks, key)
	}
	lock.Unlock()
}

// Wouldn't we all like to know?
func isProtocolUp(protocol string) bool {
	protocol = strings.ToLower(protocol)
	return protocol == "up"
}
