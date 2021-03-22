package backend

import (
	"time"
)

// Store State Constants
const (
	STATE_INIT = iota
	STATE_READY
	STATE_UPDATING
	STATE_ERROR
)

// StoreStatus defines a status the store can be in
type StoreStatus struct {
	LastRefresh time.Time
	LastError   error
	State       int
}

// Helper: stateToString
func stateToString(state int) string {
	switch state {
	case STATE_INIT:
		return "INIT"
	case STATE_READY:
		return "READY"
	case STATE_UPDATING:
		return "UPDATING"
	case STATE_ERROR:
		return "ERROR"
	}
	return "INVALID"
}
