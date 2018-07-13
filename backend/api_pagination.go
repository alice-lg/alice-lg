package main

/*
Paginate api routes responses
*/

import (
	"github.com/alice-lg/alice-lg/backend/api"

	"math"
)

func apiPaginateRoutes(
	routes api.Routes, page, pageSize int,
) (api.Routes, api.Pagination) {

	totalResults := len(routes)
	totalPages := int(math.Ceil(float64(totalResults) / float64(pageSize)))

	offset := page * pageSize
	rindex := offset + pageSize

	// Don't access out of bounds
	if rindex > totalResults {
		rindex = totalResults
	}
	if offset < 0 {
		offset = 0
	}

	pagination := api.Pagination{
		Page:         page,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	// Safeguards
	if offset >= totalResults {
		return api.Routes{}, pagination
	}

	return routes[offset:rindex], pagination
}
