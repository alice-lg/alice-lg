package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
)

/*
API Search

* Helper methods for searching
* Handle filter criteria

*/
type ApiSearchFilter struct {
	ResultsCount int         `json:"results"`
	Name         string      `json:"name"`
	Value        interface{} `json:"value"`
}

type ApiSearchFilterGroup struct {
	Name string `json:"name"`
	Key  string `json:"key"`

	Filters []ApiSearchFilters `json:"filters"`
}

type ApiSearchFilters []ApiSearchFilterGroup

/*
 Show filter criteria available
*/
func apiSearchAvailableFilters(routes *api.LookupRoutes) {
	filterSources := []ApiSearchFilter{}
	filterNeighbors := []ApiSearchFilter{}

	// Groups
	groupSources := ApiSearchFilterGroup{
		Name: "Routeservers",
		Key:  "rsid",

		filters: filterSources,
	}

}
