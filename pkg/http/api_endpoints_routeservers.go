package http

import (
	"context"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Handle RouteServers List
func (s *Server) apiRouteServersList(
	ctx context.Context,
	_req *http.Request,
	_params httprouter.Params,
) (response, error) {
	// Get list of sources from config,
	routeservers := api.RouteServers{}

	sources := s.cfg.Sources
	for _, source := range sources {
		routeservers = append(routeservers, api.RouteServer{
			ID:         source.ID,
			Type:       source.Type,
			Name:       source.Name,
			Group:      source.Group,
			Blackholes: source.Blackholes,
			Order:      source.Order,
		})
	}

	// Assert routeserver ordering
	sort.Sort(routeservers)

	// Make routeservers response
	response := api.RouteServersResponse{
		RouteServers: routeservers,
	}

	return response, nil
}
