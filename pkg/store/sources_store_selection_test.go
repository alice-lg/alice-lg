package store

import (
	"testing"
	"time"
)

func TestGetSourceIDsForRefreshSkipsBusyAndNotDue(t *testing.T) {
	now := time.Now().UTC()
	s := &SourcesStore{
		refreshInterval:    30 * time.Minute,
		refreshParallelism: 2,
		status: map[string]*Status{
			"busy-oldest": {
				SourceID:         "busy-oldest",
				State:            StateBusy,
				lastRefreshStart: now.Add(-time.Hour),
			},
			"ready-not-due": {
				SourceID:         "ready-not-due",
				State:            StateReady,
				LastRefresh:      now.Add(-time.Minute),
				lastRefreshStart: now.Add(-50 * time.Minute),
			},
			"error-retry-due": {
				SourceID:         "error-retry-due",
				State:            StateError,
				LastRefresh:      now.Add(-11 * time.Second),
				lastRefreshStart: now.Add(-12 * time.Second),
			},
		},
	}

	ids := s.GetSourceIDsForRefresh()
	if len(ids) != 1 {
		t.Fatalf("expected 1 id, got %v", ids)
	}
	if ids[0] != "error-retry-due" {
		t.Errorf("expected error-retry-due, got %s", ids[0])
	}
}
