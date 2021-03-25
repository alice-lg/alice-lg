package openbgpd

import (
	"context"
	"net/http"
	"strings"
)

func joinURL(prefix, path string) string {
	u := prefix
	if !strings.HasSuffix(prefix, "/") {
		u += "/"
	}
	u += path
	return u
}

// StatusRequest makes status request from source
func StatusRequest(ctx context.Context, s *Source) (*http.Request, error) {
	url := joinURL(s.API, "/v1/status")
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}
