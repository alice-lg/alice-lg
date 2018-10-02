package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

// Handle routes
func apiRoutesList(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	neighbourId := params.ByName("neighbourId")
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Routes(neighbourId)
	if err != nil {
		apiLogSourceError("routes", rsId, neighbourId, err)
	}

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
		apiLogSourceError("routes_received", rsId, neighbourId, err)
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
		apiLogSourceError("routes_filtered", rsId, neighbourId, err)
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
		apiLogSourceError("routes_not_exported", rsId, neighbourId, err)
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
