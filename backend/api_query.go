package main

import (
	"net/http"
	"strconv"
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
