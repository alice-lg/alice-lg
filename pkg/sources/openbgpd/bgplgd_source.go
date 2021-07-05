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
	// BgplgdSourceVersion is currently fixed at 1.0
	BgplgdSourceVersion = "1.0"
)

// BgplgdSource implements a source for Alice, consuming
// the openbgp bgplgd.
type BgplgdSource struct {
	// cfg is the source configuration retrieved
	// from the alice config file.
	cfg *Config

	// Store the neighbor responses from the server here
	neighborsCache *caches.NeighborsCache

	// Store the routes responses from the server
	// here identified by neighborID
	routesReceivedCache *caches.RoutesCache
}

// NewBgplgdSource creates a new source instance with a configuration.
func NewBgplgdSource(cfg *Config) *BgplgdSource {
	cacheDisabled := cfg.CacheTTL == 0

	// Initialize caches
	nc := caches.NewNeighborsCache(cacheDisabled)
	rrc := caches.NewRoutesCache(cacheDisabled, 128) // configure this?

	return &BgplgdSource{
		cfg:                 cfg,
		neighborsCache:      nc,
		routesReceivedCache: rrc,
	}
}

// ExpireCaches ... will flush the cache.
func (src *BgplgdSource) ExpireCaches() int {
	totalExpired := src.routesReceivedCache.Expire()
	return totalExpired
}

// Requests
// ========

// ShowNeighborsRequest makes an all neighbors request
func (src *BgplgdSource) ShowNeighborsRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/neighbors")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborsSummaryRequest builds an neighbors status request
func (src *BgplgdSource) ShowNeighborsSummaryRequest(
	ctx context.Context,
) (*http.Request, error) {
	url := src.cfg.APIURL("/summary")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborRIBRequest retrives the routes accepted from the neighbor
// identified by bgp-id.
func (src *BgplgdSource) ShowNeighborRIBRequest(
	ctx context.Context,
	neighborID string,
) (*http.Request, error) {
	url := src.cfg.APIURL("/rib?neighbor=%s", neighborID)
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowRIBRequest makes a request for retrieving all routes imported
// from all peers
func (src *BgplgdSource) ShowRIBRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/rib")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// Datasource
// ==========

// makeCacheStatus will create a new api status with cache infos
func (src *BgplgdSource) makeCacheStatus() api.ApiStatus {
	return api.ApiStatus{
		CacheStatus: api.CacheStatus{
			CachedAt: time.Now().UTC(),
		},
		Version:         BgplgdSourceVersion,
		ResultFromCache: false,
		Ttl:             time.Now().UTC().Add(src.cfg.CacheTTL),
	}
}

// Status returns an API status response. In our case
// this is pretty much only that the service is available.
func (src *BgplgdSource) Status() (*api.StatusResponse, error) {
	// Make API request and read response. We do not cache the result.
	response := &api.StatusResponse{
		Api: src.makeCacheStatus(),
		Status: api.Status{
			Version: "openbgpd",
			Message: "openbgpd up and running",
		},
	}
	return response, nil
}

// Neighbours retrievs a full list of all neighbors
func (src *BgplgdSource) Neighbours() (*api.NeighboursResponse, error) {
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
	}
	response = &api.NeighboursResponse{
		Api:        src.makeCacheStatus(),
		Neighbours: nb,
	}
	src.neighborsCache.Set(response)

	return response, nil
}

// NeighboursStatus retrives the status summary
// for all neightbors
func (src *BgplgdSource) NeighboursStatus() (*api.NeighboursStatusResponse, error) {
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

	response := &api.NeighboursStatusResponse{
		Api:        src.makeCacheStatus(),
		Neighbours: nb,
	}
	return response, nil
}

// Routes retrieves the routes for a specific neighbor
// identified by ID.
func (src *BgplgdSource) Routes(neighborID string) (*api.RoutesResponse, error) {
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

	recv, err := decodeRoutes(body)
	if err != nil {
		return nil, err
	}

	response = &api.RoutesResponse{
		Api:         src.makeCacheStatus(),
		Imported:    recv,
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	src.routesReceivedCache.Set(neighborID, response)

	return response, nil
}

// RoutesReceived returns the routes exported by the neighbor.
func (src *BgplgdSource) RoutesReceived(neighborID string) (*api.RoutesResponse, error) {
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

	recv, err := decodeRoutes(body)
	if err != nil {
		return nil, err
	}

	response = &api.RoutesResponse{
		Api:         src.makeCacheStatus(),
		Imported:    recv,
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	src.routesReceivedCache.Set(neighborID, response)

	return response, nil
}

// RoutesFiltered retrieves the routes filtered / not valid
func (src *BgplgdSource) RoutesFiltered(neighborID string) (*api.RoutesResponse, error) {
	response := &api.RoutesResponse{
		Api: src.makeCacheStatus(),

		Imported:    api.Routes{},
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	return response, nil
}

// RoutesNotExported retrievs the routes not exported
// from the rs for a neighbor.
func (src *BgplgdSource) RoutesNotExported(neighborID string) (*api.RoutesResponse, error) {
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
func (src *BgplgdSource) AllRoutes() (*api.RoutesResponse, error) {
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

	recv, err := decodeRoutes(body)
	if err != nil {
		return nil, err
	}

	response := &api.RoutesResponse{
		Api:         src.makeCacheStatus(),
		Imported:    recv,
		NotExported: api.Routes{},
		Filtered:    api.Routes{},
	}
	return response, nil
}
