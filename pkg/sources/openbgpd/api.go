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

// NeighborsRequest makes an all neighbors request
func NeighborsRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/bgpd/show/neighbor")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// NeighborsSummaryRequest builds an neighbors status request
func NeighborsSummaryRequest(ctx context.Context, src *Source) (*http.Request, error) {
	url := apiURL(src.API, "v1/bgpd/show/summary")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// NeighborRoutesReceivedRequest retrives the RIB IN for neighbor identified
// by bgp-id.
func NeighborRoutesReceivedRequest(
	ctx context.Context,
	src *Source,
	neighborID string,
) (*http.Request, error) {
	url := apiURL(src.API, "v1/bgpd/show/rib/in/neighbor/%s/detail", neighborID)
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}
