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

	// In case pageSize is 0, we assume pagination
	// is disabled.
	if pageSize == 0 {
		pagination := api.Pagination{
			Page:         page,
			PageSize:     pageSize,
			TotalPages:   0,
			TotalResults: totalResults,
		}
		return routes, pagination
	}

	// Calculate the number of pages we get
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
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	// Safeguards
	if offset >= totalResults {
		return api.Routes{}, pagination
	}

	return routes[offset:rindex], pagination
}

func apiPaginateLookupRoutes(
	routes api.LookupRoutes,
	page, pageSize int,
) (api.LookupRoutes, api.Pagination) {
	totalResults := len(routes)

	// In case pageSize is 0, we assume pagination
	// is disabled.
	if pageSize == 0 {
		pagination := api.Pagination{
			Page:         page,
			PageSize:     pageSize,
			TotalPages:   0,
			TotalResults: totalResults,
		}
		return routes, pagination
	}

	// Calculate the number of pages we get
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
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	// Safeguards
	if offset >= totalResults {
		return api.LookupRoutes{}, pagination
	}

	return routes[offset:rindex], pagination
}
