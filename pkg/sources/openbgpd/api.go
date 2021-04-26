package openbgpd

import (
	"context"
	"net/http"
	"strings"
)

func apiURL(prefix, path string) string {
	u := prefix
	if !strings.HasSuffix(prefix, "/") {
		u += "/"
	}
	u += path
	return u
}

// StatusRequest makes status request from source
func StatusRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/status")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// NeighborsRequest makes an all neighbors request
func NeighborsRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/show/neighbor")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// NeighborsSummaryRequest builds an neighbors status request
func NeighborsSummaryRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/show/summary")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}
