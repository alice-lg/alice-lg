package http

import (
	"errors"
	"fmt"
	"strings"

	"net/http"
)

var (
	// ErrQueryTooShort will be returned when the query
	// is less than 2 characters.
	ErrQueryTooShort = errors.New("query too short")

	// ErrQueryIncomplete will be returned when the
	// prefix query lacks a : or .
	ErrQueryIncomplete = errors.New(
		"prefix query must contain at least on '.' or ':'")
)

// Helper: Validate source Id
func validateSourceID(id string) (string, error) {
	if len(id) > 42 {
		return "unknown", fmt.Errorf("source ID too long with length: %d", len(id))
	}
	return id, nil
}

// Helper: Validate query string
func validateQueryString(req *http.Request, key string) (string, error) {
	query := req.URL.Query()
	values, ok := query[key]
	if !ok {
		return "", fmt.Errorf("query param %s is missing", key)
	}

	if len(values) != 1 {
		return "", fmt.Errorf("query param %s is ambigous", key)
	}

	value := values[0]
	if value == "" {
		return "", fmt.Errorf("query param %s may not be empty", key)
	}

	return value, nil
}

// Helper: Validate prefix query. It should contain
// at least one dot or :
func validatePrefixQuery(value string) (string, error) {
	// We should at least provide 2 chars
	if len(value) < 2 {
		return "", ErrQueryTooShort
	}
	if !strings.Contains(value, ":") && !strings.Contains(value, ".") {
		return "", ErrQueryIncomplete
	}
	return value, nil
}

// Helper: Validate neighbors query. A valid query should have
// at least 4 chars.
func validateNeighborsQuery(value string) (string, error) {
	if len(value) < 3 {
		// Maybe make configurable,
		// A length of 3 would be sufficient for "DFN" and
		// other shorthands.
		return "", ErrQueryTooShort
	}
	return value, nil
}
