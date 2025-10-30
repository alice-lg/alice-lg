package openbgpd

import (
	"github.com/alice-lg/alice-lg/pkg/api"
	"slices"
)

func filterReceivedRoutes(
	rejectCommunities api.Communities,
	routes api.Routes,
) api.Routes {
	filtered := make(api.Routes, 0, len(routes))
	for _, r := range routes {
		received := true
		if slices.ContainsFunc(rejectCommunities, r.BGP.HasLargeCommunity) {
			received = false
		}
		if received {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func filterRejectedRoutes(
	rejectCommunities api.Communities,
	routes api.Routes,
) api.Routes {
	filtered := make(api.Routes, 0, len(routes))
	for _, r := range routes {
		rejected := slices.ContainsFunc(rejectCommunities, r.BGP.HasLargeCommunity)
		if rejected {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
