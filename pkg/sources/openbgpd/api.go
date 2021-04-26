package openbgpd

import (
	"context"
	"net/http"
)

// StatusRequest makes status request from source
func (src *Source) StatusRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/status")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborsRequest makes an all neighbors request
func (src *Source) ShowNeighborsRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/neighbor")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborsSummaryRequest builds an neighbors status request
func (src *Source) ShowNeighborsSummaryRequest(
	ctx context.Context,
) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/summary")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborRIBInRequest retrives the routes accepted from the neighbor
// identified by bgp-id.
func (src *Source) ShowNeighborRIBInRequest(
	ctx context.Context,
	neighborID string,
) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/rib/in/neighbor/%s/detail", neighborID)
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowRIBRequest makes a request for retrieving all routes imported
// from all peers
func (src *Source) ShowRIBRequest(ctx context.Context) (*http.Request, error) {
	url := src.cfg.APIURL("/v1/bgpd/show/rib/in/detail")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}
