package openbgpd

import (
	"context"
	"net/http"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/caches"
	"github.com/alice-lg/alice-lg/pkg/decoders"
)

const (
	// StateServerSourceVersion is currently fixed at 1.0
	StateServerSourceVersion = "1.0"
)

// StateServerSource implements the OpenBGPD source for Alice.
// It is intendet to consume structured bgpctl output
// queried over HTTP using the:
//
//    openbgpd-state-server
//    https://github.com/alice-lg/openbgpd-state-server
//
type StateServerSource struct {
	// cfg is the source configuration retrieved
	// from the alice config file.
	cfg *Config

	// Store the neighbor responses from the server here
	neighborsCache *caches.NeighborsCache

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
	rc := caches.NewRoutesCache(cacheDisabled, cfg.RoutesCacheSize)
	rrc := caches.NewRoutesCache(cacheDisabled, cfg.RoutesCacheSize)
	rfc := caches.NewRoutesCache(cacheDisabled, cfg.RoutesCacheSize)

	return &StateServerSource{
		cfg:                 cfg,
		neighborsCache:      nc,
		routesCache:         rc,
		routesReceivedCache: rrc,
		routesFilteredCache: rfc,
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

// ShowNeighborRIBRequest retrives the routes accepted from the neighbor
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

// makeCacheStatus will create a new api status with cache infos
func (src *StateServerSource) makeCacheStatus() api.ApiStatus {
	return api.ApiStatus{
		CacheStatus: api.CacheStatus{
			CachedAt: time.Now().UTC(),
		},
		Version:         StateServerSourceVersion,
		ResultFromCache: false,
		Ttl:             time.Now().UTC().Add(src.cfg.CacheTTL),
	}
}

// Status returns an API status response. In our case
// this is pretty much only that the service is available.
func (src *StateServerSource) Status() (*api.StatusResponse, error) {
	// Make API request and read response. We do not cache the result.
	req, err := src.StatusRequest(context.Background())
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
		Api:    src.makeCacheStatus(),
		Status: status,
	}
	return response, nil
}

// Neighbors retrievs a full list of all neighbors
func (src *StateServerSource) Neighbors() (*api.NeighborsResponse, error) {
	// Query cache and see if we have a hit
	response := src.neighborsCache.Get()
	if response != nil {
		response.Api.ResultFromCache = true
		return response, nil
	}

	// Make API request and read response
	req, err := src.ShowNeighborsRequest(context.Background())
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
		n.RouteServerId = src.cfg.ID

		rejectedRes, err := src.RoutesFiltered(n.Id)
		if err != nil {
			return nil, err
		}
		rejectCount := len(rejectedRes.Filtered)
		n.RoutesFiltered = rejectCount

	}
	response = &api.NeighborsResponse{
		Api:        src.makeCacheStatus(),
		Neighbors: nb,
	}
	src.neighborsCache.Set(response)

	return response, nil
}

// NeighborsStatus retrives the status summary
// for all neightbors
func (src *StateServerSource) NeighborsStatus() (*api.NeighborsStatusResponse, error) {
	// Make API request and read response
	req, err := src.ShowNeighborsSummaryRequest(context.Background())
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

	response := &api.NeighborsStatusResponse{
		Api:        src.makeCacheStatus(),
		Neighbors: nb,
	}
	return response, nil
}

// Routes retrieves the routes for a specific neighbor
// identified by ID.
func (src *StateServerSource) Routes(neighborID string) (*api.RoutesResponse, error) {
	response := src.routesCache.Get(neighborID)
	if response != nil {
		response.Api.ResultFromCache = true
		return response, nil
	}

	// Query RIB for routes received
	req, err := src.ShowNeighborRIBRequest(context.Background(), neighborID)
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
		Api:         src.makeCacheStatus(),
		Imported:    received,
		NotExported: api.Routes{},
		Filtered:    rejected,
	}
	src.routesCache.Set(neighborID, response)

	return response, nil
}

// RoutesReceived returns the routes exported by the neighbor.
func (src *StateServerSource) RoutesReceived(neighborID string) (*api.RoutesResponse, error) {
	response := src.routesReceivedCache.Get(neighborID)
	if response != nil {
		response.Api.ResultFromCache = true
		return response, nil
	}

	// Query RIB for routes received
	req, err := src.ShowNeighborRIBRequest(context.Background(), neighborID)
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
		Api:         src.makeCacheStatus(),
		Imported:    received,
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	src.routesReceivedCache.Set(neighborID, response)

	return response, nil
}

// RoutesFiltered retrieves the routes filtered / not valid
func (src *StateServerSource) RoutesFiltered(neighborID string) (*api.RoutesResponse, error) {
	response := src.routesFilteredCache.Get(neighborID)
	if response != nil {
		response.Api.ResultFromCache = true
		return response, nil
	}

	// Query RIB for routes received
	req, err := src.ShowNeighborRIBRequest(context.Background(), neighborID)
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
		Api:         src.makeCacheStatus(),
		Imported:    api.Routes{},
		NotExported: api.Routes{},
		Filtered:    rejected,
	}
	src.routesFilteredCache.Set(neighborID, response)

	return response, nil
}

// RoutesNotExported retrievs the routes not exported
// from the rs for a neighbor.
func (src *StateServerSource) RoutesNotExported(neighborID string) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{
		Api: src.makeCacheStatus(),

		Imported:    api.Routes{},
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	return response, nil
}

// AllRoutes retrievs the entire RIB from the source. This is never
// cached as it is processed by the store.
func (src *StateServerSource) AllRoutes() (*api.RoutesResponse, error) {
	req, err := src.ShowRIBRequest(context.Background())
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
		Api:         src.makeCacheStatus(),
		Imported:    received,
		NotExported: api.Routes{},
		Filtered:    rejected,
	}
	return response, nil
}
