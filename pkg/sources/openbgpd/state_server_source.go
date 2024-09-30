package openbgpd

import (
	"context"
	"net/http"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/caches"
	"github.com/alice-lg/alice-lg/pkg/decoders"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// Ensure source interface is implemented
var _OpenBGPStateServerSource sources.Source = &StateServerSource{}

const (
	// StateServerSourceVersion is currently fixed at 1.0
	StateServerSourceVersion = "1.0"
)

// StateServerSource implements the OpenBGPD source for Alice.
// It is intended to consume structured bgpctl output
// queried over HTTP using the:
//
//	openbgpd-state-server
//	https://github.com/alice-lg/openbgpd-state-server
type StateServerSource struct {
	// cfg is the source configuration retrieved
	// from the alice config file.
	cfg *Config

	// Store the neighbor responses from the server here
	neighborsCache        *caches.NeighborsCache
	neighborsSummaryCache *caches.NeighborsCache

	// Store the routes responses from the server
	// here identified by neighborID
	routesCache         *caches.RoutesCache
	routesReceivedCache *caches.RoutesCache
	routesFilteredCache *caches.RoutesCache
}

// NewStateServerSource creates a new source instance with a
// configuration.
func NewStateServerSource(cfg *Config) *StateServerSource {
	cacheDisabled := cfg.CacheTTL == 0

	// Initialize caches
	nc := caches.NewNeighborsCache(cacheDisabled)
	nsc := caches.NewNeighborsCache(cacheDisabled)
	rc := caches.NewRoutesCache(cacheDisabled, cfg.RoutesCacheSize)
	rrc := caches.NewRoutesCache(cacheDisabled, cfg.RoutesCacheSize)
	rfc := caches.NewRoutesCache(cacheDisabled, cfg.RoutesCacheSize)

	return &StateServerSource{
		cfg:                   cfg,
		neighborsCache:        nc,
		neighborsSummaryCache: nsc,
		routesCache:           rc,
		routesReceivedCache:   rrc,
		routesFilteredCache:   rfc,
	}
}

// ExpireCaches ... will flush the cache. Seriously this needs
// a renaming.
func (src *StateServerSource) ExpireCaches() int {
	totalExpired := src.routesCache.Expire()
	return totalExpired
}

// Requests
// ========

// StatusRequest makes status request from source
func (src *StateServerSource) StatusRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/status")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborsRequest makes an all neighbors request
func (src *StateServerSource) ShowNeighborsRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/neighbor")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborsSummaryRequest builds an neighbors status request
func (src *StateServerSource) ShowNeighborsSummaryRequest(
	ctx context.Context,
) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/summary")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborRIBRequest retrieves the routes accepted from the neighbor
// identified by bgp-id.
func (src *StateServerSource) ShowNeighborRIBRequest(
	ctx context.Context,
	neighborID string,
) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/rib/neighbor/%s/detail", neighborID)
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowRIBRequest makes a request for retrieving all routes imported
// from all peers
func (src *StateServerSource) ShowRIBRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/rib/detail")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// Datasource
// ==========

// makeResponseMeta will create a new api status with cache infos
func (src *StateServerSource) makeResponseMeta() *api.Meta {
	return &api.Meta{
		CacheStatus: api.CacheStatus{
			CachedAt: time.Now().UTC(),
		},
		Version:         StateServerSourceVersion,
		ResultFromCache: false,
		TTL:             time.Now().UTC().Add(src.cfg.CacheTTL),
	}
}

// Status returns an API status response. In our case
// this is pretty much only that the service is available.
func (src *StateServerSource) Status(
	ctx context.Context,
) (*api.StatusResponse, error) {
	// Make API request and read response. We do not cache the result.
	req, err := src.StatusRequest(ctx)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}
	status := decodeAPIStatus(body)
	response := &api.StatusResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Status: status,
	}
	return response, nil
}

// Neighbors retrieves a full list of all neighbors
func (src *StateServerSource) Neighbors(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	// Query cache and see if we have a hit
	response := src.neighborsCache.Get()
	if response != nil {
		response.Response.Meta.ResultFromCache = true
		return response, nil
	}

	// Make API request and read response
	req, err := src.ShowNeighborsRequest(ctx)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	nb, err := decodeNeighbors(body)
	if err != nil {
		return nil, err
	}
	// Set route server id (sourceID) for all neighbors
	for _, n := range nb {
		n.RouteServerID = src.cfg.ID

		rejectedRes, err := src.RoutesFiltered(ctx, n.ID)
		if err != nil {
			return nil, err
		}
		rejectCount := len(rejectedRes.Filtered)
		n.RoutesFiltered = rejectCount

	}
	response = &api.NeighborsResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Neighbors: nb,
	}
	response.Neighbors, err = sources.FilterHiddenNeighbors(response.Neighbors, src.cfg.HiddenNeighbors)
	if err != nil {
		return nil, err
	}
	src.neighborsCache.Set(response)

	return response, nil
}

// NeighborsSummary retrieves the neighbors without additional
// information but as quickly as possible. The result will lack
// a reject count.
func (src *StateServerSource) NeighborsSummary(
	ctx context.Context,
) (*api.NeighborsResponse, error) {
	response := src.neighborsSummaryCache.Get()
	if response != nil {
		response.Meta.ResultFromCache = true
		return response, nil
	}

	// Make API request and read response
	req, err := src.ShowNeighborsRequest(ctx)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	nb, err := decodeNeighbors(body)
	if err != nil {
		return nil, err
	}
	nb, err = sources.FilterHiddenNeighbors(nb, src.cfg.HiddenNeighbors)
	if err != nil {
		return nil, err
	}
	// Set route server id (sourceID) for all neighbors
	for _, n := range nb {
		n.RouteServerID = src.cfg.ID
	}

	response = &api.NeighborsResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Neighbors: nb,
	}
	src.neighborsSummaryCache.Set(response)
	return response, nil
}

// NeighborsStatus retrieves the status summary
// for all neighbors
func (src *StateServerSource) NeighborsStatus(
	ctx context.Context,
) (*api.NeighborsStatusResponse, error) {
	// Make API request and read response
	req, err := src.ShowNeighborsSummaryRequest(ctx)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and decode response
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	nb, err := decodeNeighborsStatus(body)
	if err != nil {
		return nil, err
	}
	nb, err = sources.FilterHiddenNeighborsStatus(nb, src.cfg.HiddenNeighbors)
	if err != nil {
		return nil, err
	}

	response := &api.NeighborsStatusResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Neighbors: nb,
	}
	return response, nil
}

// Routes retrieves the routes for a specific neighbor
// identified by ID.
func (src *StateServerSource) Routes(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := src.routesCache.Get(neighborID)
	if response != nil {
		response.Response.Meta.ResultFromCache = true
		return response, nil
	}

	// Query RIB for routes received
	req, err := src.ShowNeighborRIBRequest(ctx, neighborID)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and decode response
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	routes, err := decodeRoutes(body)
	if err != nil {
		return nil, err
	}

	// Filtered routes are marked with a large BGP community
	// as defined in the reject reasons.
	received := filterReceivedRoutes(src.cfg.RejectCommunities, routes)
	rejected := filterRejectedRoutes(src.cfg.RejectCommunities, routes)

	response = &api.RoutesResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Imported:    received,
		NotExported: api.Routes{},
		Filtered:    rejected,
	}
	src.routesCache.Set(neighborID, response)

	return response, nil
}

// RoutesReceived returns the routes exported by the neighbor.
func (src *StateServerSource) RoutesReceived(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := src.routesReceivedCache.Get(neighborID)
	if response != nil {
		response.Response.Meta.ResultFromCache = true
		return response, nil
	}

	// Query RIB for routes received
	req, err := src.ShowNeighborRIBRequest(ctx, neighborID)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and decode response
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	routes, err := decodeRoutes(body)
	if err != nil {
		return nil, err
	}

	received := filterReceivedRoutes(src.cfg.RejectCommunities, routes)

	response = &api.RoutesResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Imported:    received,
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	src.routesReceivedCache.Set(neighborID, response)

	return response, nil
}

// RoutesFiltered retrieves the routes filtered / not valid
func (src *StateServerSource) RoutesFiltered(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := src.routesFilteredCache.Get(neighborID)
	if response != nil {
		response.Response.Meta.ResultFromCache = true
		return response, nil
	}

	// Query RIB for routes received
	req, err := src.ShowNeighborRIBRequest(ctx, neighborID)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and decode response
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	routes, err := decodeRoutes(body)
	if err != nil {
		return nil, err
	}

	rejected := filterRejectedRoutes(src.cfg.RejectCommunities, routes)

	response = &api.RoutesResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Imported:    api.Routes{},
		NotExported: api.Routes{},
		Filtered:    rejected,
	}
	src.routesFilteredCache.Set(neighborID, response)

	return response, nil
}

// RoutesNotExported retrieves the routes not exported
// from the rs for a neighbor.
func (src *StateServerSource) RoutesNotExported(
	ctx context.Context,
	neighborID string,
) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Imported:    api.Routes{},
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	return response, nil
}

// AllRoutes retrieves the entire RIB from the source. This is never
// cached as it is processed by the store.
func (src *StateServerSource) AllRoutes(
	ctx context.Context,
) (*api.RoutesResponse, error) {
	req, err := src.ShowRIBRequest(ctx)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and decode response
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	routes, err := decodeRoutes(body)
	if err != nil {
		return nil, err
	}

	// Filtered routes are marked with a large BGP community
	// as defined in the reject reasons.
	received := filterReceivedRoutes(src.cfg.RejectCommunities, routes)
	rejected := filterRejectedRoutes(src.cfg.RejectCommunities, routes)

	response := &api.RoutesResponse{
		Response: api.Response{
			Meta: src.makeResponseMeta(),
		},
		Imported:    received,
		NotExported: api.Routes{},
		Filtered:    rejected,
	}
	return response, nil
}
