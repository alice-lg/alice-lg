package api

import (
	"testing"
)

func TestApiSearchFilterGetGroupsByKey(t *testing.T) {
	filtering := NewApiSearchFilters()

	group := filtering.GetGroupByKey(API_SEARCH_KEY_ASNS)
	if group == nil {
		t.Error(API_SEARCH_KEY_ASNS, "should exis")
		return
	}

	if group.Key != API_SEARCH_KEY_ASNS {
		t.Error("group should be:", API_SEARCH_KEY_ASNS, "but is:", group.Key)
	}
}

func TestApiSearchFilterManagement(t *testing.T) {
	filtering := NewApiSearchFilters()
	group := filtering.GetGroupByKey(API_SEARCH_KEY_ASNS)

	group.AddFilter(&ApiSearchFilter{
		Name:  "Tech Inc. Solutions GmbH",
		Value: 23042})

	group.AddFilter(&ApiSearchFilter{
		Name:  "T3ch Inc. Solutions GmbH",
		Value: 23042})

	group.AddFilter(&ApiSearchFilter{
		Name:  "Foocom Telecommunications Ltd.",
		Value: 424242})

	// Check filters

	filter := group.GetFilterByValue(424242)
	if filter == nil {
		t.Error("Expected filter for AS424242")
	}

	filter = group.GetFilterByValue(23042)
	if filter == nil {
		t.Error("Expected filter for AS23042")
		return
	}

	if filter.Cardinality != 2 {
		t.Error("Expected a cardinality of 2, got:", filter.Cardinality)
	}
}
