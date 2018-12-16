package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
	"sort"
)

// Handle Routeservers List
func apiRouteserversList(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	// Get list of sources from config,
	routeservers := api.Routeservers{}

	sources := AliceConfig.Sources
	for _, source := range sources {
		routeservers = append(routeservers, api.Routeserver{
			Id:         source.Id,
			Name:       source.Name,
			Group:      source.Group,
			Blackholes: source.Blackholes,
			Order:      source.Order,
		})
	}

	// Assert routeserver ordering
	sort.Sort(routeservers)

	// Make routeservers response
	response := api.RouteserversResponse{
		Routeservers: routeservers,
	}

	return response, nil
}
