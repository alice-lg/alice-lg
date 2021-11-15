package http

import (
	"fmt"

	"net/http"
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

// Helper: Validate prefix query
func validatePrefixQuery(value string) (string, error) {
	// We should at least provide 2 chars
	if len(value) < 2 {
		return "", fmt.Errorf("query too short")
	}
	return value, nil
}
