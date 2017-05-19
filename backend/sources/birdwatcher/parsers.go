package birdwatcher

// Parsers and helpers

import (
	"time"
)

const SERVER_TIME_SHORT = "2006-01-02 15:04:05"
const SERVER_TIME_EXT = "Mon, 2 Jan 2006 15:04:05 +0000"

// Convert server time string to time
func parseServerTime(value interface{}, layout, timezone string) (time.Time, error) {
	svalue, ok := value.(string)
	if !ok {
		return time.Time{}, nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation(layout, svalue, loc)
	return t, err
}
