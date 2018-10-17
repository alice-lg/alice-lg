package api

import (
	"net/url"
	"testing"
)

func TestSearchFilterGetGroupsByKey(t *testing.T) {
	filtering := NewSearchFilters()

	group := filtering.GetGroupByKey(SEARCH_KEY_ASNS)
	if group == nil {
		t.Error(SEARCH_KEY_ASNS, "should exis")
		return
	}

	if group.Key != SEARCH_KEY_ASNS {
		t.Error("group should be:", SEARCH_KEY_ASNS, "but is:", group.Key)
	}
}

func TestSearchFilterManagement(t *testing.T) {
	filtering := NewSearchFilters()
	group := filtering.GetGroupByKey(SEARCH_KEY_ASNS)

	group.AddFilter(&SearchFilter{
		Name:  "Tech Inc. Solutions GmbH",
		Value: 23042})

	group.AddFilter(&SearchFilter{
		Name:  "T3ch Inc. Solutions GmbH",
		Value: 23042})

	group.AddFilter(&SearchFilter{
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

func TestSearchFiltersFromQuery(t *testing.T) {
	query := "asns=2342,23123&communities=23:42&large_communities=23:42:42&sources=1,2,3"
	values, err := url.ParseQuery(query)
	if err != nil {
		t.Error(err)
		return
	}

	filters, err := FiltersFromQuery(values)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(filters)
}
