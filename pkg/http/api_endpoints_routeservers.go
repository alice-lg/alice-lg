package http

import (
	"context"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
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

// Handle route server status
func (s *Server) apiRouteServerStatusShow(
	ctx context.Context,
	_req *http.Request,
	params httprouter.Params,
) (response, error) {
	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	source := s.cfg.SourceInstanceByID(rsID)
	if source == nil {
		return nil, ErrSourceNotFound
	}

	result, err := source.Status(ctx)
	if err != nil {
		s.logSourceError("status", rsID, err)
		return nil, err
	}
	if result != nil {
		// Prevent panic if *api.Meta is null
		if result.Meta != nil {
			result.Meta.Version = config.Version
		}
	}

	return result, nil
}
