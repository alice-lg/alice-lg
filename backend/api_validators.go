package main

import (
	"fmt"
	"strconv"

	"net/http"
)

// Helper: Validate source Id
func validateSourceId(id string) (string, error) {
	if len(id) > 42 {
		return "unknown", fmt.Errorf("Source ID too long with length: %d", len(id))
	}
	return id, nil
}

// Helper: Validate query string
func validateQueryString(req *http.Request, key string) (string, error) {
	query := req.URL.Query()
	values, ok := query[key]
	if !ok {
		return "", fmt.Errorf("Query param %s is missing.", key)
	}

	if len(values) != 1 {
		return "", fmt.Errorf("Query param %s is ambigous.", key)
	}

	value := values[0]
	if value == "" {
		return "", fmt.Errorf("Query param %s may not be empty.", key)
	}

	return value, nil
}

// Helper: Validate prefix query
func validatePrefixQuery(value string) (string, error) {

	// We should at least provide 2 chars
	if len(value) < 2 {
		return "", fmt.Errorf("Query too short")
	}

	// Query constraints: Should at least include a dot or colon
	/* let's try without this :)

	if strings.Index(value, ".") == -1 &&
		strings.Index(value, ":") == -1 {
		return "", fmt.Errorf("Query needs at least a ':' or '.'")
	}
	*/

	return value, nil
}

// Get pagination parameters: limit and offset
// Refer to defaults if none are given.
func validatePaginationParams(req *http.Request, limit, offset int) (int, int, error) {
	query := req.URL.Query()
	queryLimit, ok := query["limit"]
	if ok {
		limit, _ = strconv.Atoi(queryLimit[0])
	}

	queryOffset, ok := query["offset"]
	if ok {
		offset, _ = strconv.Atoi(queryOffset[0])
	}

	// Cap limit to [1, 1000]
	if limit < 1 {
		limit = 1
	}
	if limit > 500 {
		limit = 500
	}

	return limit, offset, nil
}
