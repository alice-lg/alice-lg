package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

// Handle Routeservers List
func apiRouteserversList(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	// Get list of sources from config,
	routeservers := []api.Routeserver{}

	sources := AliceConfig.Sources
	for _, source := range sources {
		routeservers = append(routeservers, api.Routeserver{
			Id:         source.Id,
			Name:       source.Name,
			Blackholes: source.Blackholes,
		})
	}

	// Make routeservers response
	response := api.RouteserversResponse{
		Routeservers: routeservers,
	}

	return response, nil
}
