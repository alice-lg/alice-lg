// Package store provides a store for persisting and
// querying neighbors and routes.
package store

import (
	"time"
)

// Store State Constants
const (
	StateInit = iota
	StateReady
	StateUpdating
	StateError
)

// Status defines a status the store can be in
type Status struct {
	LastRefresh time.Time
	LastError   error
	State       int
}

// Helper: stateToString
func stateToString(state int) string {
	switch state {
	case StateInit:
		return "INIT"
	case StateReady:
		return "READY"
	case StateUpdating:
		return "UPDATING"
	case StateError:
		return "ERROR"
	}
	return "INVALID"
}
