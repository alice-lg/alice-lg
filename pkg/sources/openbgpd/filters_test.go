package openbgpd

import (
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestFilterReceivedRoutes(t *testing.T) {
	routes := api.Routes{
		&api.Route{
			Id: "1.2.3.4",
			Bgp: api.BgpInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 666, 1},
				},
			},
		},
		&api.Route{
			Id: "5.6.6.6",
			Bgp: api.BgpInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 5, 42},
					api.Community{9999, 666, 2},
				},
			},
		},
		&api.Route{
			Id: "5.6.7.8",
			Bgp: api.BgpInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 5, 42},
				},
			},
		},
	}
	c := api.Communities{
		api.Community{9999, 666, 1},
		api.Community{9999, 666, 2},
	}
	filtered := filterReceivedRoutes(c, routes)

	if filtered[0].Id != "5.6.7.8" {
		t.Error("unexpected route:", filtered[0])
	}
}

func TestFilterRejectedRoutes(t *testing.T) {
	routes := api.Routes{
		&api.Route{
			Id: "5.6.7.8",
			Bgp: api.BgpInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 5, 42},
				},
			},
		},
		&api.Route{
			Id: "1.2.3.4",
			Bgp: api.BgpInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 666, 1},
				},
			},
		},
		&api.Route{
			Id: "5.6.6.6",
			Bgp: api.BgpInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 5, 42},
					api.Community{9999, 666, 2},
				},
			},
		},
	}
	c := api.Communities{
		api.Community{9999, 666, 1},
		api.Community{9999, 666, 2},
	}
	filtered := filterRejectedRoutes(c, routes)

	if len(filtered) != 2 {
		t.Error("expected two filtered routes")
	}

	if filtered[0].Id != "1.2.3.4" {
		t.Error("unexpected route:", filtered[0])
	}
}
