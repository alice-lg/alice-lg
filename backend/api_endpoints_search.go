package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
	"time"
)

// Handle global lookup
func apiLookupPrefixGlobal(req *http.Request, params httprouter.Params) (api.Response, error) {
	// Get prefix to query
	q, err := validateQueryString(req, "q")
	if err != nil {
		return nil, err
	}

	q, err = validatePrefixQuery(q)
	if err != nil {
		return nil, err
	}
	// Get pagination params
	limit, offset, err := validatePaginationParams(req, 50, 0)
	if err != nil {
		return nil, err
	}

	// Check what we want to query
	//  Prefix -> fetch prefix
	//       _ -> fetch neighbours and routes
	lookupPrefix := MaybePrefix(q)

	// Measure response time
	t0 := time.Now()

	// Perform query
	var routes api.LookupRoutes
	if lookupPrefix {
		routes = AliceRoutesStore.LookupPrefix(q)

	} else {
		neighbours := AliceNeighboursStore.LookupNeighbours(q)
		routes = AliceRoutesStore.LookupPrefixForNeighbours(neighbours)
	}

	// Paginate result
	totalRoutes := len(routes)
	cap := offset + limit
	if cap > totalRoutes {
		cap = totalRoutes
	}

	queryDuration := time.Since(t0)
	response := api.RoutesLookupResponseGlobal{
		Routes: routes[offset:cap],

		TotalRoutes: totalRoutes,
		Limit:       limit,
		Offset:      offset,

		Time: float64(queryDuration) / 1000.0 / 1000.0, // nano -> micro -> milli
	}

	return response, nil
}
