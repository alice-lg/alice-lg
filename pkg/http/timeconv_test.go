package http

import (
	"testing"
	"time"
)

func TestDurationMs(t *testing.T) {
	if DurationMs(time.Second) != 1000 {
		t.Error("duration ms should return the duration in milliseconds")
	}
}
