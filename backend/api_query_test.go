package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/alice-lg/alice-lg/backend/api"
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
			Id:          "route_01",
			NeighbourId: "n01",
			Network:     "123.42.43.0/24",
			Gateway:     "23.42.42.1",
		},
		&api.Route{
			Id:          "route_02",
			NeighbourId: "n01",
			Network:     "142.23.0.0/16",
			Gateway:     "42.42.42.1",
		},
		&api.Route{
			Id:          "route_03",
			NeighbourId: "n01",
			Network:     "123.43.0.0/16",
			Gateway:     "23.42.43.1",
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
		t.Error("Exptected 2 routes, got:", len(filtered))
	}

	// Check presence of route_01 and _03, matching prefix 123.
	if filtered[0].Id != "route_01" {
		t.Error("Expected route_01, got:", filtered[0].Id)
	}
	if filtered[1].Id != "route_03" {
		t.Error("Expected route_03, got:", filtered[1].Id)
	}

	// Test another query matching the gateway only
	req = makeQueryRequest("42.")
	filtered = apiQueryFilterNextHopGateway(
		req, "q", routes,
	)

	if len(filtered) != 1 {
		t.Error("Expected only one result")
	}

	if filtered[0].Id != "route_02" {
		t.Error("Expected route_02 to match criteria, got:", filtered[0])
	}
}
