package api

import (
	"net/url"
	"testing"
)

func makeTestRoute() *LookupRoute {
	route := &LookupRoute{
		Bgp: BgpInfo{
			Communities: []Community{
				Community{23, 42},
				Community{11, 111},
			},
			ExtCommunities: []ExtCommunity{
				ExtCommunity{"ro", 23, 123},
			},
			LargeCommunities: []Community{
				Community{1000, 23, 42},
			},
		},
		Neighbour: &Neighbour{
			Asn:         23042,
			Description: "Security Solutions Ltd.",
		},
		Routeserver: Routeserver{
			Id:   3,
			Name: "test.rs.ixp",
		},
	}

	return route
}

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
	query := "asns=2342,23123&large_communities=23:42:42&sources=1,2,3&q=foo"
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

	// We should have 2 asns present
	asns := filters.GetGroupByKey(SEARCH_KEY_ASNS).Filters
	if asns[0].Value != 2342 {
		t.Error("Expected asn.filter[0].Value to be 2342, but got:", asns[0].Value)
	}
	if asns[1].Value != 23123 {
		t.Error("Expected asn.filter[1].Value to be 23123, but got:", asns[1].Value)
	}

	// Check communities
	communities := filters.GetGroupByKey(SEARCH_KEY_COMMUNITIES).Filters
	if len(communities) != 0 {
		t.Error("There should be no communities filters")
	}

	largeCommunities := filters.GetGroupByKey(SEARCH_KEY_LARGE_COMMUNITIES).Filters
	if len(largeCommunities) != 1 {
		t.Error("There should have been 1 large community")
	}

	if largeCommunities[0].Name != "23:42:42" {
		t.Error("There should have been 23:42:42 as a large community filter")
	}

	t.Log(largeCommunities[0].Value)

	// Check Sources
	sources := filters.GetGroupByKey(SEARCH_KEY_SOURCES).Filters
	if len(sources) != 3 {
		t.Error("Expected 3 source id filters")
	}
}

func TestSearchFilterCompareRoute(t *testing.T) {
	// Check filter matches
	route := makeTestRoute()

	// Source
	if searchFilterMatchSource(route, 3) != true {
		t.Error("Route should have sourceId 3")
	}
	if searchFilterMatchSource(route, 23) == true {
		t.Error("Route should not have sourceId 23")
	}

	// Asn
	if searchFilterMatchAsn(route, 23042) != true {
		t.Error("Route should have ASN 23042")
	}
	if searchFilterMatchAsn(route, 123) == true {
		t.Error("Route should not have ASN 123")
	}

	// Communities
	if searchFilterMatchCommunity(route, Community{11, 111}) != true {
		t.Error("Route should have community 11:111")
	}
	if searchFilterMatchCommunity(route, Community{42, 111}) == true {
		t.Error("Route should not have community 42:111")
	}

	// Ext. Communities
	if searchFilterMatchExtCommunity(route, ExtCommunity{"ro", 23, 123}) != true {
		t.Error("Route should have community ro:23:123")
	}
	if searchFilterMatchExtCommunity(route, ExtCommunity{"rt", 42, 111}) == true {
		t.Error("Route should not have community rt:42:111")
	}

	// Large Communities
	if searchFilterMatchLargeCommunity(route, Community{1000, 23, 42}) != true {
		t.Error("Route should have community 1000:23:42")
	}
	if searchFilterMatchLargeCommunity(route, Community{42, 111, 111}) == true {
		t.Error("Route should not have community 42:111:111")
	}
}

func TestSearchFilterMatchRoute(t *testing.T) {
	route := makeTestRoute()

	query := "asns=2342,23042&large_communities=1000:23:42&sources=1,2,3&q=foo"
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

	if filters.MatchRoute(route) == false {
		t.Error("Route should have matched filters...")
	}

}

func TestSearchFilterExcludeRoute(t *testing.T) {
	route := makeTestRoute()

	query := "asns=2342,23042&large_communities=42:23:42&sources=1,2,3&q=foo"
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

	if filters.MatchRoute(route) != false {
		t.Error("Route should not have matched filters...")
	}
}
