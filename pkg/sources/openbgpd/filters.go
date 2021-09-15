package openbgpd

import (
	"github.com/alice-lg/alice-lg/pkg/api"
)

func filterReceivedRoutes(
	rejectCommunities api.Communities,
	routes api.Routes,
) api.Routes {
	filtered := make(api.Routes, 0, len(routes))
	for _, r := range routes {
		received := true
		for _, c := range rejectCommunities {
			if r.Bgp.HasLargeCommunity(c) {
				received = false
				break
			}
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
		rejected := false
		for _, c := range rejectCommunities {
			if r.Bgp.HasLargeCommunity(c) {
				rejected = true
				break
			}
		}
		if rejected {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
