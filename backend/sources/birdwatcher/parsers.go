package birdwatcher

// Parsers and helpers

import (
	"time"

	"github.com/ecix/alice-lg/backend/api"
)

const SERVER_TIME = time.RFC3339Nano
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

// Make api status from response:
// The api status is always included in a birdwatcher response
func parseApiStatus(bird ClientResponse, config Config) (api.ApiStatus, error) {
	birdApi := bird["api"].(map[string]interface{})

	ttl, err := parseServerTime(
		bird["ttl"],
		SERVER_TIME,
		config.Timezone,
	)
	if err != nil {
		return api.ApiStatus{}, err
	}

	status := api.ApiStatus{
		Version:         birdApi["Version"].(string),
		ResultFromCache: birdApi["result_from_cache"].(bool),
		Ttl:             ttl,
	}

	return status, nil
}

// Parse neighbours response
func parseNeighbours(bird ClientResponse) ([]api.Neighbour, error) {

	return []api.Neighbour{}, nil
}
