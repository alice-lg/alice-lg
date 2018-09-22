package main

import (
	"compress/gzip"
	"encoding/json"
	"net/http"

	"log"
	"strings"

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
//     Neighbors    /api/routeservers/:id/neighbours
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
			// Get affected rs id
			rsId, paramErr := validateSourceId(params.ByName("id"))
			if paramErr != nil {
				rsId = -1
			}

			// Make error response
			result = apiErrorResponse(rsId, err)
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
		endpoint(apiNeighborsList))
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
