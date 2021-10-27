package http

import (
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestApiRoutesPagination(t *testing.T) {
	routes := api.Routes{
		&api.Route{ID: "r01"},
		&api.Route{ID: "r02"},
		&api.Route{ID: "r03"},
		&api.Route{ID: "r04"},
		&api.Route{ID: "r05"},
		&api.Route{ID: "r06"},
		&api.Route{ID: "r07"},
		&api.Route{ID: "r08"},
		&api.Route{ID: "r09"},
		&api.Route{ID: "r10"},
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
	if r.ID != "r01" {
		t.Error("First route on page 0 should be r01, got:", r.ID)
	}

	r = paginated[len(paginated)-1]
	if r.ID != "r08" {
		t.Error("Last route should be r08, but got:", r.ID)
	}

	// Second page
	paginated, _ = apiPaginateRoutes(routes, 1, 8)
	if len(paginated) != 2 {
		t.Error("There should be 2 routes left on page 1, got:", len(paginated))
	}

	r = paginated[0]
	if r.ID != "r09" {
		t.Error("First route on page 1 should be r09, got:", r.ID)
	}

	r = paginated[len(paginated)-1]
	if r.ID != "r10" {
		t.Error("Last route should be r10, but got:", r.ID)
	}

	// Access out of bound page
	paginated, _ = apiPaginateRoutes(routes, 1000, 8)
	if len(paginated) > 0 {
		t.Error("There should be nothing on this page")
	}
}
