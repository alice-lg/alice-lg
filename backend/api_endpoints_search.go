package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
	"sort"
	"time"
)

// Handle global lookup
func apiLookupPrefixGlobal(
	req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	// TODO: This function is too long

	// Get prefix to query
	q, err := validateQueryString(req, "q")
	if err != nil {
		return nil, err
	}

	q, err = validatePrefixQuery(q)
	if err != nil {
		return nil, err
	}

	// Check what we want to query
	//  Prefix -> fetch prefix
	//       _ -> fetch neighbours and routes
	lookupPrefix := MaybePrefix(q)

	// Measure response time
	t0 := time.Now()

	// Get additional filter criteria
	filtersApplied, err := api.FiltersFromQuery(req.URL.Query())
	if err != nil {
		return nil, err
	}

	// Perform query
	var routes api.LookupRoutes
	if lookupPrefix {
		routes = AliceRoutesStore.LookupPrefix(q)

	} else {
		neighbours := AliceNeighboursStore.LookupNeighbours(q)
		routes = AliceRoutesStore.LookupPrefixForNeighbours(neighbours)
	}

	// Split routes
	// TODO: Refactor at neighbors store
	totalResults := len(routes)
	imported := make(api.LookupRoutes, 0, totalResults)
	filtered := make(api.LookupRoutes, 0, totalResults)

	// Now, as we have allocated even more space process routes by, splitting,
	// filtering and updating the available filters...
	filtersAvailable := api.NewSearchFilters()
	for _, r := range routes {

		if !filtersApplied.MatchRoute(r) {
			continue // Exclude route from results set
		}

		switch r.State {
		case "filtered":
			filtered = append(filtered, r)
			break
		case "imported":
			imported = append(imported, r)
			break
		}

		filtersAvailable.UpdateFromLookupRoute(r)
	}

	// Remove applied filters from available
	filtersApplied.MergeProperties(filtersAvailable)
	filtersAvailable = filtersAvailable.Sub(filtersApplied)

	// Homogenize results
	sort.Sort(imported)
	sort.Sort(filtered)

	// Paginate results
	pageImported := apiQueryMustInt(req, "page_imported", 0)
	pageSizeImported := AliceConfig.Ui.Pagination.RoutesAcceptedPageSize
	routesImported, paginationImported := apiPaginateLookupRoutes(
		imported, pageImported, pageSizeImported,
	)

	pageFiltered := apiQueryMustInt(req, "page_filtered", 0)
	pageSizeFiltered := AliceConfig.Ui.Pagination.RoutesFilteredPageSize
	routesFiltered, paginationFiltered := apiPaginateLookupRoutes(
		filtered, pageFiltered, pageSizeFiltered,
	)

	// Calculate query duration
	queryDuration := time.Since(t0)

	// Make response
	response := api.PaginatedRoutesLookupResponse{
		Api: api.ApiStatus{
			CacheStatus: api.CacheStatus{
				CachedAt: AliceRoutesStore.CachedAt(),
			},
			ResultFromCache: true, // Well.
			Ttl:             AliceRoutesStore.CacheTtl(),
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		Imported: &api.LookupRoutesResponse{
			Routes: routesImported,
			PaginatedResponse: &api.PaginatedResponse{
				Pagination: paginationImported,
			},
		},
		Filtered: &api.LookupRoutesResponse{
			Routes: routesFiltered,
			PaginatedResponse: &api.PaginatedResponse{
				Pagination: paginationFiltered,
			},
		},
		FilterableResponse: api.FilterableResponse{
			FiltersAvailable: filtersAvailable,
			FiltersApplied:   filtersApplied,
		},
	}

	return response, nil
}
