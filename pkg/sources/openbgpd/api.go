package openbgpd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func apiURL(prefix, path string, params ...interface{}) string {
	u := prefix
	if !strings.HasSuffix(prefix, "/") {
		u += "/"
	}
	u += fmt.Sprintf(path, params...)
	return u
}

// StatusRequest makes status request from source
func StatusRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/status")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborsRequest makes an all neighbors request
func ShowNeighborsRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/bgpd/show/neighbor")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborsSummaryRequest builds an neighbors status request
func ShowNeighborsSummaryRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/bgpd/show/summary")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowNeighborRIBInRequest retrives the routes accepted from the neighbor
// identified by bgp-id.
func ShowNeighborRIBInRequest(
	ctx context.Context,
	src *Source,
	neighborID string,
) (*http.Request, error) {
	url := apiURL(src.API, "v1/bgpd/show/rib/in/neighbor/%s/detail", neighborID)
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// ShowRIBRequest makes a request for retrieving all routes imported
// from all peers
func ShowRIBRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "vi/bgpd/show/rib/in/detail")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}
