package http

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Handle routes
func (s *Server) apiRoutesList(
	_req *http.Request,
	params httprouter.Params,
) (response, error) {
	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	neighborID := params.ByName("neighborId")

	source := s.cfg.SourceInstanceByID(rsID)
	if source == nil {
		return nil, ErrSourceNotFound
	}

	result, err := source.Routes(neighborID)
	if err != nil {
		s.logSourceError("routes", rsID, neighborID, err)
	}

	return result, err
}

// Paginated Routes Respponse: Received routes
func (s *Server) apiRoutesListReceived(
	req *http.Request,
	params httprouter.Params,
) (response, error) {
	// Measure response time
	t0 := time.Now()

	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighborID := params.ByName("neighborId")
	source := s.cfg.SourceInstanceByID(rsID)
	if source == nil {
		return nil, ErrSourceNotFound
	}

	result, err := source.RoutesReceived(neighborID)
	if err != nil {
		s.logSourceError("routes_received", rsID, neighborID, err)
		return nil, err
	}

	// Filter routes based on criteria if present
	allRoutes := apiQueryFilterNextHopGateway(req, "q", result.Imported)
	routes := api.Routes{}

	// Apply other (commmunity) filters
	filtersApplied, err := api.FiltersFromQuery(req.URL.Query())
	if err != nil {
		return nil, err
	}

	filtersAvailable := api.NewSearchFilters()
	for _, r := range allRoutes {
		if !filtersApplied.MatchRoute(r) {
			continue // Exclude route from results set
		}
		routes = append(routes, r)
		filtersAvailable.UpdateFromRoute(r)
	}

	// Remove applied filters from available
	filtersApplied.MergeProperties(filtersAvailable)
	filtersAvailable = filtersAvailable.Sub(filtersApplied)

	// Paginate results
	page := apiQueryMustInt(req, "page", 0)
	pageSize := s.cfg.UI.Pagination.RoutesAcceptedPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Calculate query duration
	queryDuration := time.Since(t0)

	// Make paginated response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: api.RoutesResponse{
			Response: api.Response{
				Meta: result.Response.Meta,
			},
			Imported: routes,
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		FilteredResponse: api.FilteredResponse{
			FiltersAvailable: filtersAvailable,
			FiltersApplied:   filtersApplied,
		},
		PaginatedResponse: api.PaginatedResponse{
			Pagination: pagination,
		},
	}

	return response, nil
}

func (s *Server) apiRoutesListFiltered(
	req *http.Request,
	params httprouter.Params,
) (response, error) {
	t0 := time.Now()

	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighborID := params.ByName("neighborId")
	source := s.cfg.SourceInstanceByID(rsID)
	if source == nil {
		return nil, ErrSourceNotFound
	}

	result, err := source.RoutesFiltered(neighborID)
	if err != nil {
		s.logSourceError("routes_filtered", rsID, neighborID, err)
		return nil, err
	}

	// Filter routes based on criteria if present
	allRoutes := apiQueryFilterNextHopGateway(req, "q", result.Filtered)
	routes := api.Routes{}

	// Apply other (commmunity) filters
	filtersApplied, err := api.FiltersFromQuery(req.URL.Query())
	if err != nil {
		return nil, err
	}

	filtersAvailable := api.NewSearchFilters()
	for _, r := range allRoutes {
		if !filtersApplied.MatchRoute(r) {
			continue // Exclude route from results set
		}
		routes = append(routes, r)
		filtersAvailable.UpdateFromRoute(r)
	}

	// Remove applied filters from available
	filtersApplied.MergeProperties(filtersAvailable)
	filtersAvailable = filtersAvailable.Sub(filtersApplied)

	// Paginate results
	page := apiQueryMustInt(req, "page", 0)
	pageSize := s.cfg.UI.Pagination.RoutesFilteredPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Calculate query duration
	queryDuration := time.Since(t0)

	// Make response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: api.RoutesResponse{
			Response: api.Response{
				Meta: result.Response.Meta,
			},
			Filtered: routes,
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		FilteredResponse: api.FilteredResponse{
			FiltersAvailable: filtersAvailable,
			FiltersApplied:   filtersApplied,
		},
		PaginatedResponse: api.PaginatedResponse{
			Pagination: pagination,
		},
	}

	return response, nil
}

func (s *Server) apiRoutesListNotExported(
	req *http.Request,
	params httprouter.Params,
) (response, error) {
	t0 := time.Now()

	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighborID := params.ByName("neighborId")
	source := s.cfg.SourceInstanceByID(rsID)
	if source == nil {
		return nil, ErrSourceNotFound
	}

	result, err := source.RoutesNotExported(neighborID)
	if err != nil {
		s.logSourceError("routes_not_exported", rsID, neighborID, err)
		return nil, err
	}

	// Filter routes based on criteria if present
	allRoutes := apiQueryFilterNextHopGateway(req, "q", result.NotExported)
	routes := api.Routes{}

	// Apply other (commmunity) filters
	filtersApplied, err := api.FiltersFromQuery(req.URL.Query())
	if err != nil {
		return nil, err
	}

	filtersAvailable := api.NewSearchFilters()
	for _, r := range allRoutes {
		if !filtersApplied.MatchRoute(r) {
			continue // Exclude route from results set
		}
		routes = append(routes, r)
		filtersAvailable.UpdateFromRoute(r)
	}

	// Remove applied filters from available
	filtersApplied.MergeProperties(filtersAvailable)
	filtersAvailable = filtersAvailable.Sub(filtersApplied)

	// Paginate results
	page := apiQueryMustInt(req, "page", 0)
	pageSize := s.cfg.UI.Pagination.RoutesNotExportedPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Calculate query duration
	queryDuration := time.Since(t0)

	// Make response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: api.RoutesResponse{
			Response: api.Response{
				Meta: result.Response.Meta,
			},
			NotExported: routes,
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		FilteredResponse: api.FilteredResponse{
			FiltersAvailable: filtersAvailable,
			FiltersApplied:   filtersApplied,
		},
		PaginatedResponse: api.PaginatedResponse{
			Pagination: pagination,
		},
	}

	return response, nil
}
