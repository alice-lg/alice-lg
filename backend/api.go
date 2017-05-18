package main

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ecix/alice-lg/backend/api"

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

	// Config
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

	return nil
}

// Handle Config Endpoint
func apiConfigShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	result := api.ConfigResponse{
		Rejection: api.Rejection{
			Asn:      AliceConfig.Ui.RoutesRejections.Asn,
			RejectId: AliceConfig.Ui.RoutesRejections.RejectId,
		},
		RejectReasons: AliceConfig.Ui.RoutesRejections.Reasons,
		Noexport: api.Noexport{
			Asn:        AliceConfig.Ui.RoutesNoexports.Asn,
			NoexportId: AliceConfig.Ui.RoutesNoexports.NoexportId,
		},
		NoexportReasons: AliceConfig.Ui.RoutesNoexports.Reasons,
		RoutesColumns:   AliceConfig.Ui.RoutesColumns,
	}
	return result, nil
}

// Handle Routeservers List
func apiRouteserversList(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	// Get list of sources from config,
	routeservers := []api.Routeserver{}

	sources := AliceConfig.Sources
	for id, source := range sources {
		routeservers = append(routeservers, api.Routeserver{
			Id:   id,
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
	rsId, _ := strconv.Atoi(params.ByName("id"))
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Status()
	return result, err
}

// Handle get neighbours on routeserver
func apiNeighboursList(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, _ := strconv.Atoi(params.ByName("id"))
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Neighbours()
	return result, err
}

// Handle routes
func apiRoutesList(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, _ := strconv.Atoi(params.ByName("id"))
	neighbourId := params.ByName("neighbourId")
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Routes(neighbourId)
	return result, err
}
