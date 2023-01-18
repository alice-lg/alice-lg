package http

import (
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestApiRoutesPagination(t *testing.T) {
	routes := api.Routes{
		&api.Route{Network: "r01"},
		&api.Route{Network: "r02"},
		&api.Route{Network: "r03"},
		&api.Route{Network: "r04"},
		&api.Route{Network: "r05"},
		&api.Route{Network: "r06"},
		&api.Route{Network: "r07"},
		&api.Route{Network: "r08"},
		&api.Route{Network: "r09"},
		&api.Route{Network: "r10"},
	}

	paginated, pagination := apiPaginateRoutes(routes, 0, 8)

	if pagination.TotalPages != 2 {
		t.Error("Expected total pages to be 2, got:", pagination.TotalPages)
	}

	if pagination.TotalResults != 10 {
		t.Error("Expected total results to be 10, got:", pagination.TotalResults)
	}

	if pagination.Page != 0 {
		t.Error("Exptected current page to be 0, got:", pagination.Page)
	}

	// Check paginated slicing
	r := paginated[0]
	if r.Network != "r01" {
		t.Error("First route on page 0 should be r01, got:", r.Network)
	}

	r = paginated[len(paginated)-1]
	if r.Network != "r08" {
		t.Error("Last route should be r08, but got:", r.Network)
	}

	// Second page
	paginated, _ = apiPaginateRoutes(routes, 1, 8)
	if len(paginated) != 2 {
		t.Error("There should be 2 routes left on page 1, got:", len(paginated))
	}

	r = paginated[0]
	if r.Network != "r09" {
		t.Error("First route on page 1 should be r09, got:", r.Network)
	}

	r = paginated[len(paginated)-1]
	if r.Network != "r10" {
		t.Error("Last route should be r10, but got:", r.Network)
	}

	// Access out of bound page
	paginated, _ = apiPaginateRoutes(routes, 1000, 8)
	if len(paginated) > 0 {
		t.Error("There should be nothing on this page")
	}
}
