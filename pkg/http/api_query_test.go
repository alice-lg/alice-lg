package http

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/pools"
)

func makeQueryRequest(q string) *http.Request {
	url, _ := url.Parse("http://alice/api?q=" + q)
	req := &http.Request{
		URL: url,
	}
	return req
}

func makeQueryRoutes() api.Routes {
	routes := api.Routes{
		&api.Route{
			NeighborID: pools.Neighbors.Acquire("n01"),
			Network:    "123.42.43.0/24",
			Gateway:    pools.Gateways4.Acquire("23.42.42.1"),
		},
		&api.Route{
			NeighborID: pools.Neighbors.Acquire("n01"),
			Network:    "142.23.0.0/16",
			Gateway:    pools.Gateways4.Acquire("42.42.42.1"),
		},
		&api.Route{
			NeighborID: pools.Neighbors.Acquire("n01"),
			Network:    "123.43.0.0/16",
			Gateway:    pools.Gateways4.Acquire("23.42.43.1"),
		},
	}

	return routes
}

func TestApiQueryFilterNextHopGateway(t *testing.T) {

	routes := makeQueryRoutes()
	req := makeQueryRequest("123.")
	filtered := apiQueryFilterNextHopGateway(
		req, "q", routes,
	)

	if len(filtered) != 2 {
		t.Error("Expected 2 routes, got:", len(filtered))
	}

	// Check presence of route_01 and _03, matching prefix 123.
	if filtered[0].Network != "123.42.43.0/24" {
		t.Error("Expected 123.42.43.0/24 got:", filtered[0].Network)
	}
	if filtered[1].Network != "123.43.0.0/16" {
		t.Error("Expected 123.43.0.0/16, got:", filtered[1].Network)
	}

	// Test another query matching the gateway only
	req = makeQueryRequest("42.")
	filtered = apiQueryFilterNextHopGateway(
		req, "q", routes,
	)

	if len(filtered) != 1 {
		t.Error("Expected only one result")
	}

	if filtered[0].Network != "142.23.0.0/16" {
		t.Error("Expected 142.23.0.0/16 to match criteria, got:", filtered[0])
	}
}
