package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/alice-lg/alice-lg/pkg/api"
)

/*
Convenience methods for accessing the query string
in the request object.
*/

/*
Get int value by name from query string
*/
func apiQueryMustInt(req *http.Request, param string, defaultValue int) int {
	query := req.URL.Query()
	strVal, ok := query[param]
	if !ok {
		return defaultValue
	}

	value, err := strconv.Atoi(strVal[0])
	if err != nil {
		return defaultValue
	}

	return value
}

/*
Filter response to match query criteria
*/

func apiQueryFilterNextHopGateway(
	req *http.Request, param string, routes api.Routes,
) api.Routes {
	query := req.URL.Query()
	queryParam, ok := query[param]
	if !ok {
		return routes
	}

	// Normalize to lowercase
	queryString := strings.ToLower(queryParam[0])

	results := make(api.Routes, 0, len(routes))
	for _, r := range routes {
		if strings.HasPrefix(strings.ToLower(r.Network), queryString) ||
			strings.HasPrefix(strings.ToLower(r.Gateway), queryString) {
			results = append(results, r)
		}
	}

	return results
}
