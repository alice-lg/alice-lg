package http

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
//     Show         /api/v1/config
//
//   Routeservers
//     List         /api/v1/routeservers
//     Status       /api/v1/routeservers/:id/status
//     Neighbors    /api/v1/routeservers/:id/neighbors
//     Routes       /api/v1/routeservers/:id/neighbors/:neighborId/routes
//
//   Querying
//     LookupPrefix   /api/v1/lookup/prefix?q=<prefix>
//     LookupNeighbor /api/v1/lookup/neighbor?asn=1235

type response interface{}

type apiEndpoint func(*http.Request, httprouter.Params) (response, error)

// Wrap handler for access controll, throtteling and compression
func endpoint(wrapped apiEndpoint) httprouter.Handle {
	return func(res http.ResponseWriter,
		req *http.Request,
		params httprouter.Params) {

		// Get result from handler
		result, err := wrapped(req, params)
		if err != nil {
			// Get affected rs id
			rsID, paramErr := validateSourceID(params.ByName("id"))
			if paramErr != nil {
				rsID = "unknown"
			}

			// Make error response
			result, status := apiErrorResponse(rsID, err)
			payload, _ := json.Marshal(result)
			http.Error(res, string(payload), status)
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
func (s *Server) apiRegisterEndpoints(
	router *httprouter.Router,
) error {

	// Meta
	router.GET("/api/v1/status", endpoint(s.apiStatusShow))
	router.GET("/api/v1/config", endpoint(s.apiConfigShow))

	// Routeservers
	router.GET("/api/v1/routeservers",
		endpoint(s.apiRouteServersList))
	router.GET("/api/v1/routeservers/:id/status",
		endpoint(s.apiStatus))
	router.GET("/api/v1/routeservers/:id/neighbors",
		endpoint(s.apiNeighborsList))
	router.GET("/api/v1/routeservers/:id/neighbors/:neighborId/routes",
		endpoint(s.apiRoutesList))
	router.GET("/api/v1/routeservers/:id/neighbors/:neighborId/routes/received",
		endpoint(s.apiRoutesListReceived))
	router.GET("/api/v1/routeservers/:id/neighbors/:neighborId/routes/filtered",
		endpoint(s.apiRoutesListFiltered))
	router.GET("/api/v1/routeservers/:id/neighbors/:neighborId/routes/not-exported",
		endpoint(s.apiRoutesListNotExported))

	// Querying
	if s.cfg.Server.EnablePrefixLookup == true {
		router.GET("/api/v1/lookup/prefix",
			endpoint(s.apiLookupPrefixGlobal))
		router.GET("/api/v1/lookup/neighbors",
			endpoint(s.apiLookupNeighborsGlobal))
	}

	return nil
}
