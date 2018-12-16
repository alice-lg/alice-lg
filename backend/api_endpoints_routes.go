package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
	"time"
)

// Handle routes
func apiRoutesList(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	neighborId := params.ByName("neighborId")
	source := AliceConfig.SourceById(rsId).getInstance()
	result, err := source.Routes(neighborId)
	if err != nil {
		apiLogSourceError("routes", rsId, neighborId, err)
	}

	return result, err
}

// Paginated Routes Respponse: Received routes
func apiRoutesListReceived(
	req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	// Measure response time
	t0 := time.Now()

	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighborId := params.ByName("neighborId")
	source := AliceConfig.SourceById(rsId).getInstance()
	result, err := source.RoutesReceived(neighborId)
	if err != nil {
		apiLogSourceError("routes_received", rsId, neighborId, err)
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
	pageSize := AliceConfig.Ui.Pagination.RoutesAcceptedPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Calculate query duration
	queryDuration := time.Since(t0)

	// Make paginated response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: &api.RoutesResponse{
			Api:      result.Api,
			Imported: routes,
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		FilterableResponse: api.FilterableResponse{
			FiltersAvailable: filtersAvailable,
			FiltersApplied:   filtersApplied,
		},
		Pagination: pagination,
	}

	return response, nil
}

func apiRoutesListFiltered(
	req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	t0 := time.Now()

	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighborId := params.ByName("neighborId")
	source := AliceConfig.SourceById(rsId).getInstance()
	result, err := source.RoutesFiltered(neighborId)
	if err != nil {
		apiLogSourceError("routes_filtered", rsId, neighborId, err)
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
	pageSize := AliceConfig.Ui.Pagination.RoutesFilteredPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Calculate query duration
	queryDuration := time.Since(t0)

	// Make response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: &api.RoutesResponse{
			Api:      result.Api,
			Filtered: routes,
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		FilterableResponse: api.FilterableResponse{
			FiltersAvailable: filtersAvailable,
			FiltersApplied:   filtersApplied,
		},
		Pagination: pagination,
	}

	return response, nil
}

func apiRoutesListNotExported(
	req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	t0 := time.Now()

	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	neighborId := params.ByName("neighborId")
	source := AliceConfig.SourceById(rsId).getInstance()
	result, err := source.RoutesNotExported(neighborId)
	if err != nil {
		apiLogSourceError("routes_not_exported", rsId, neighborId, err)
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
	pageSize := AliceConfig.Ui.Pagination.RoutesNotExportedPageSize
	routes, pagination := apiPaginateRoutes(routes, page, pageSize)

	// Calculate query duration
	queryDuration := time.Since(t0)

	// Make response
	response := api.PaginatedRoutesResponse{
		RoutesResponse: &api.RoutesResponse{
			Api:         result.Api,
			NotExported: routes,
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		FilterableResponse: api.FilterableResponse{
			FiltersAvailable: filtersAvailable,
			FiltersApplied:   filtersApplied,
		},
		Pagination: pagination,
	}

	return response, nil
}
