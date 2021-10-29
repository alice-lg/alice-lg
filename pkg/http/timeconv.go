package http

import (
	"time"
)

// DurationMs converts time.Duration to milliseconds
func DurationMs(d time.Duration) float64 {
	return float64(d) / 1000.0 / 1000.0 // nano -> micro -> milli
}
