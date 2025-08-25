package http

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
)


// Handle global lookup
func (s *Server) apiLookupPrefixGlobal(
	ctx context.Context,
	req *http.Request,
	params httprouter.Params,
) (response, error) {
	// TODO: This function is way too long
	statusMeta := &api.StoreStatusMeta{
		Neighbors: s.neighborsStore.Status(ctx),
		Routes:    s.routesStore.Status(ctx),
	}

	// Get prefix to query
	q, err := validateQueryString(req, "q")
	if err != nil {
		return nil, err
	}

	q, filterTokens := QueryString(q).ExtractFilters()

	// Get filters from query string
	queryFilters, err := api.FiltersFromTokens(filterTokens)
	if err != nil {
		return nil, &ErrValidationFailed{
			Param:  "q",
			Reason: err.Error(),
		}
	}

	// Get additional filter criteria
	filtersApplied, err := api.FiltersFromQuery(req.URL.Query())
	if err != nil {
		return nil, err
	}

	// Merge query filters into applied filters
	filtersApplied = filtersApplied.Combine(queryFilters)

	// Select the query strategy:
	//  Prefix -> fetch prefix
	//       _ -> fetch neighbors and routes
	//
	lookupPrefix := decoders.MaybePrefix(q)
	lookupEmptyQuery := false
	if q == "" && (filtersApplied.HasGroup(api.SearchKeyCommunities) ||
		filtersApplied.HasGroup(api.SearchKeyExtCommunities) ||
		filtersApplied.HasGroup(api.SearchKeyLargeCommunities)) {
		lookupPrefix = true
		lookupEmptyQuery = true
	}

	// Measure response time
	t0 := time.Now()

	// Perform query
	var routes api.LookupRoutes
	if lookupPrefix {
		if !lookupEmptyQuery {
			q, err = validatePrefixQuery(q)
			if err != nil {
				return nil, err
			}
		}
		routes, err = s.routesStore.LookupPrefix(ctx, q, filtersApplied)
		if err != nil {
			return nil, err
		}

	} else {
		// Query by neighbors
		q, err = validateNeighborsQuery(q)
		if err != nil {
			return nil, err
		}
		neighbors, err := s.neighborsStore.LookupNeighbors(ctx, q)
		if err != nil {
			return nil, err
		}
		routes, err = s.routesStore.LookupPrefixForNeighbors(ctx, neighbors, filtersApplied)
		if err != nil {
			return nil, err
		}
	}

	// Split routes
	// TODO: Refactor at neighbors store
	totalResults := len(routes)
	imported := make(api.LookupRoutes, 0, totalResults)
	filtered := make(api.LookupRoutes, 0, totalResults)

	// Check if we should calculate community filter
	// cardinalities.
	filterCutoff := s.cfg.Server.PrefixLookupCommunityFilterCutoff
	canFilterCommunities := totalResults <= filterCutoff

	// In case there is a source filter applied, we can filter communities
	if filtersApplied.HasGroup(api.SearchKeySources) {
		canFilterCommunities = true
	}

	filtersNotAvailable := []string{}
	if !canFilterCommunities {
		filtersNotAvailable = append(
			filtersNotAvailable, api.SearchKeyCommunities)
	}

	// Now, as we have allocated even more space split routes,
	// and update the available filters...
	filtersAvailable := api.NewSearchFilters()
	populateAddrFamilyFilters(filtersAvailable)
	for _, r := range routes {

		switch r.State {
		case api.RouteStateFiltered:
			filtered = append(filtered, r)
		case api.RouteStateImported:
			imported = append(imported, r)
		}

		// Update available filters for sources and asns,
		// conditionally for communities.
		filtersAvailable.UpdateSourcesFromLookupRoute(r)
		filtersAvailable.UpdateASNSFromLookupRoute(r)

		if canFilterCommunities {
			filtersAvailable.UpdateCommunitiesFromLookupRoute(r)
		}
	}

	// Remove applied filters from available
	filtersApplied.MergeProperties(filtersAvailable)
	filtersAvailable = filtersAvailable.Sub(filtersApplied)

	// Homogenize results
	sort.Sort(imported)
	sort.Sort(filtered)

	// Paginate results
	pageImported := apiQueryMustInt(req, "page_imported", 0)
	pageSizeImported := s.cfg.UI.Pagination.RoutesAcceptedPageSize
	routesImported, paginationImported := apiPaginateLookupRoutes(
		imported, pageImported, pageSizeImported,
	)

	pageFiltered := apiQueryMustInt(req, "page_filtered", 0)
	pageSizeFiltered := s.cfg.UI.Pagination.RoutesFilteredPageSize
	routesFiltered, paginationFiltered := apiPaginateLookupRoutes(
		filtered, pageFiltered, pageSizeFiltered,
	)

	// Calculate query duration
	queryDuration := time.Since(t0)

	cachedAt := s.routesStore.CachedAt(ctx)
	ttl := s.routesStore.CacheTTL(ctx)

	// Make response
	response := api.PaginatedRoutesLookupResponse{
		Response: api.Response{
			Meta: &api.Meta{
				CacheStatus: api.CacheStatus{
					CachedAt: cachedAt,
				},
				StoreStatus:     statusMeta,
				ResultFromCache: true, // Well.
				TTL:             ttl,
			},
		},
		TimedResponse: api.TimedResponse{
			RequestDuration: DurationMs(queryDuration),
		},
		Imported: &api.RoutesLookup{
			Routes:     routesImported,
			Pagination: paginationImported,
		},
		Filtered: &api.RoutesLookup{
			Routes:     routesFiltered,
			Pagination: paginationFiltered,
		},
		FilteredResponse: api.FilteredResponse{
			FiltersAvailable:    filtersAvailable,
			FiltersNotAvailable: filtersNotAvailable,
			FiltersApplied:      filtersApplied,
		},
	}

	return response, nil
}

func (s *Server) apiLookupNeighborsGlobal(
	ctx context.Context,
	req *http.Request,
	params httprouter.Params,
) (response, error) {
	statusMeta := &api.StoreStatusMeta{
		Neighbors: s.neighborsStore.Status(ctx),
	}

	// Query neighbors store
	filter := api.NeighborFilterFromQuery(req.URL.Query())
	neighbors, err := s.neighborsStore.FilterNeighbors(ctx, filter)
	if err != nil {
		return nil, err
	}

	sort.Sort(neighbors)

	// Make response
	response := &api.NeighborsResponse{
		Response: api.Response{
			Meta: &api.Meta{
				CacheStatus: api.CacheStatus{
					CachedAt: s.neighborsStore.CachedAt(ctx),
				},
				StoreStatus:     statusMeta,
				ResultFromCache: true, // You would not have guessed.
				TTL:             s.neighborsStore.CacheTTL(ctx),
			},
		},
		Neighbors: neighbors,
	}
	return response, nil
}
