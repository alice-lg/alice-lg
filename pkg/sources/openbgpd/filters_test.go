package openbgpd

import (
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestFilterReceivedRoutes(t *testing.T) {
	routes := api.Routes{
		&api.Route{
			Network: "1.2.3.4",
			BGP: &api.BGPInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 666, 1},
				},
			},
		},
		&api.Route{
			Network: "5.6.6.6",
			BGP: &api.BGPInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 5, 42},
					api.Community{9999, 666, 2},
				},
			},
		},
		&api.Route{
			Network: "5.6.7.8",
			BGP: &api.BGPInfo{
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

	if filtered[0].Network != "5.6.7.8" {
		t.Error("unexpected route:", filtered[0])
	}
}

func TestFilterRejectedRoutes(t *testing.T) {
	routes := api.Routes{
		&api.Route{
			Network: "5.6.7.8",
			BGP: &api.BGPInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 5, 42},
				},
			},
		},
		&api.Route{
			Network: "1.2.3.4",
			BGP: &api.BGPInfo{
				LargeCommunities: api.Communities{
					api.Community{9999, 23, 23},
					api.Community{9999, 666, 1},
				},
			},
		},
		&api.Route{
			Network: "5.6.6.6",
			BGP: &api.BGPInfo{
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

	if filtered[0].Network != "1.2.3.4" {
		t.Error("unexpected route:", filtered[0])
	}
}
