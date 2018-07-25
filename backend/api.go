package main

import (
	"compress/gzip"
	"encoding/json"
	"net/http"

	"log"
	"strings"
	"time"

	"github.com/alice-lg/alice-lg/backend/api"

	"github.com/julienschmidt/httprouter"
)

// Alice LG Rest API
//
// The API provides endpoints for getting
// information from the routeservers / alice datasources.
//
// Endpoints:
//
//   Config
//     Show         /api/config
//
//   Routeservers
//     List         /api/routeservers
//     Status       /api/routeservers/:id/status
//     Neighbours   /api/routeservers/:id/neighbours
//     Routes       /api/routeservers/:id/neighbours/:neighbourId/routes
//
//   Querying
//     LookupPrefix /api/routeservers/:id/lookup/prefix?q=<prefix>
//

type apiEndpoint func(*http.Request, httprouter.Params) (api.Response, error)

// Wrap handler for access controll, throtteling and compression
func endpoint(wrapped apiEndpoint) httprouter.Handle {
	return func(res http.ResponseWriter,
		req *http.Request,
		params httprouter.Params) {

		// Get result from handler
		result, err := wrapped(req, params)
		if err != nil {
			result = api.ErrorResponse{
				Error: err.Error(),
			}
			payload, _ := json.Marshal(result)
			http.Error(res, string(payload), http.StatusInternalServerError)
			return
		}

		// Encode json
		payload, err := json.Marshal(result)
		if err != nil {
			msg := "Could not encode result as json"
			http.Error(res, msg, http.StatusInternalServerError)
			log.Println(err)
			log.Println("This is most likely due to an older version of go.")
			log.Println("Consider upgrading to golang > 1.8")
			return
		}

		// Set response header
		res.Header().Set("Content-Type", "application/json")

		// Check if compression is supported
		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			// Compress response
			res.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(res)
			defer gz.Close()
			gz.Write(payload)
		} else {
			res.Write(payload) // Fall back to uncompressed response
		}
	}
}

// Register api endpoints
func apiRegisterEndpoints(router *httprouter.Router) error {

	// Meta
	router.GET("/api/status", endpoint(apiStatusShow))
	router.GET("/api/config", endpoint(apiConfigShow))

	// Routeservers
	router.GET("/api/routeservers",
		endpoint(apiRouteserversList))
	router.GET("/api/routeservers/:id/status",
		endpoint(apiStatus))
	router.GET("/api/routeservers/:id/neighbours",
		endpoint(apiNeighboursList))
	router.GET("/api/routeservers/:id/neighbours/:neighbourId/routes",
		endpoint(apiRoutesList))
	router.GET("/api/routeservers/:id/neighbours/:neighbourId/routes/received",
		endpoint(apiRoutesListReceived))
	router.GET("/api/routeservers/:id/neighbours/:neighbourId/routes/filtered",
		endpoint(apiRoutesListFiltered))
	router.GET("/api/routeservers/:id/neighbours/:neighbourId/routes/not-exported",
		endpoint(apiRoutesListNotExported))

	// Querying
	if AliceConfig.Server.EnablePrefixLookup == true {
		router.GET("/api/lookup/prefix",
			endpoint(apiLookupPrefixGlobal))
	}

	return nil
}

// Handle Status Endpoint, this is intended for
// monitoring and service health checks
func apiStatusShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	status, err := NewAppStatus()
	return status, err
}

// Handle Config Endpoint
func apiConfigShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	result := api.ConfigResponse{
		Rejection: api.Rejection{
			Asn:      AliceConfig.Ui.RoutesRejections.Asn,
			RejectId: AliceConfig.Ui.RoutesRejections.RejectId,
		},
		RejectReasons: SerializeReasons(
			AliceConfig.Ui.RoutesRejections.Reasons),
		Noexport: api.Noexport{
			Asn:        AliceConfig.Ui.RoutesNoexports.Asn,
			NoexportId: AliceConfig.Ui.RoutesNoexports.NoexportId,
		},
		NoexportReasons: SerializeReasons(
			AliceConfig.Ui.RoutesNoexports.Reasons),
		RoutesColumns:          AliceConfig.Ui.RoutesColumns,
		RoutesColumnsOrder:     AliceConfig.Ui.RoutesColumnsOrder,
		NeighboursColumns:      AliceConfig.Ui.NeighboursColumns,
		NeighboursColumnsOrder: AliceConfig.Ui.NeighboursColumnsOrder,
		PrefixLookupEnabled:    AliceConfig.Server.EnablePrefixLookup,
	}
	return result, nil
}

// Handle Routeservers List
func apiRouteserversList(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	// Get list of sources from config,
	routeservers := []api.Routeserver{}

	sources := AliceConfig.Sources
	for _, source := range sources {
		routeservers = append(routeservers, api.Routeserver{
			Id:   source.Id,
			Name: source.Name,
		})
	}

	// Make routeservers response
	response := api.RouteserversResponse{
		Routeservers: routeservers,
	}

	return response, nil
}

// Handle status
func apiStatus(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Status()
	return result, err
}

// Handle get neighbours on routeserver
func apiNeighboursList(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Neighbours()
	return result, err
}

// Handle routes
func apiRoutesList(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	neighbourId := params.ByName("neighbourId")
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Routes(neighbourId)
	return result, err
}

// Paginated Routes Respponse: Received routes
func apiRoutesListReceived(
	req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighbourId := params.ByName("neighbourId")
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.RoutesReceived(neighbourId)
	if err != nil {
		return nil, err
	}

	// Filter routes based on criteria if present
	routes := apiQueryFilterNextHopGateway(req, "q", result.Imported)

	// Paginate results
	page := apiQueryMustInt(req, "page", 0)
	pageSize := AliceConfig.Ui.Pagination.RoutesAcceptedPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Make paginated response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: &api.RoutesResponse{
			Api:      result.Api,
			Imported: routes,
		},
		Pagination: pagination,
	}

	return response, nil
}

func apiRoutesListFiltered(
	req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighbourId := params.ByName("neighbourId")
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.RoutesFiltered(neighbourId)
	if err != nil {
		return nil, err
	}

	// Filter routes based on criteria if present
	routes := apiQueryFilterNextHopGateway(req, "q", result.Filtered)

	// Paginate results
	page := apiQueryMustInt(req, "page", 0)
	pageSize := AliceConfig.Ui.Pagination.RoutesFilteredPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Make response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: &api.RoutesResponse{
			Api:      result.Api,
			Filtered: routes,
		},
		Pagination: pagination,
	}

	return response, nil
}

func apiRoutesListNotExported(
	req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighbourId := params.ByName("neighbourId")
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.RoutesNotExported(neighbourId)
	if err != nil {
		return nil, err
	}

	routes := apiQueryFilterNextHopGateway(req, "q", result.NotExported)

	// Paginate results
	page := apiQueryMustInt(req, "page", 0)
	pageSize := AliceConfig.Ui.Pagination.RoutesNotExportedPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Make response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: &api.RoutesResponse{
			Api:         result.Api,
			NotExported: routes,
		},
		Pagination: pagination,
	}

	return response, nil
}

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
