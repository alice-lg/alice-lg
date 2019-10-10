package api

import (
	"net/url"
	"testing"
)

func makeTestRoute() *Route {
	route := &Route{
		Bgp: BgpInfo{
			Communities: []Community{
				Community{23, 42},
				Community{111, 11},
			},
			ExtCommunities: []ExtCommunity{
				ExtCommunity{"ro", "23", "123"},
			},
			LargeCommunities: []Community{
				Community{1000, 23, 42},
			},
		},
	}

	return route
}

func makeTestLookupRoute() *LookupRoute {
	route := &LookupRoute{
		Bgp: BgpInfo{
			Communities: []Community{
				Community{23, 42},
				Community{111, 11},
			},
			ExtCommunities: []ExtCommunity{
				ExtCommunity{"ro", "23", "123"},
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
			Id:   "3",
			Name: "test.rs.ixp",
		},
	}

	return route
}

func TestSearchFilterCmpInt(t *testing.T) {
	if searchFilterCmpInt(23, 23) != true {
		t.Error("23 == 23 should be true")
	}
	if searchFilterCmpInt(23, 42) == true {
		t.Error("23 == 42 should be false")
	}
}

func TestSearchFilterCmpCommunity(t *testing.T) {
	// Standard communities
	if searchFilterCmpCommunity(Community{23, 42}, Community{23, 42}) != true {
		t.Error("23:42 == 23:42 should be true")
	}
	if searchFilterCmpCommunity(Community{23, 42}, Community{42, 23}) == true {
		t.Error("23:42 == 42:23 should be false")
	}

	// Large communities
	if searchFilterCmpCommunity(Community{1000, 23, 42}, Community{1000, 23, 42}) != true {
		t.Error("1000:23:42 == 1000:23:42 should be true")
	}
	if searchFilterCmpCommunity(Community{1000, 23, 42}, Community{1111, 42, 23}) == true {
		t.Error("1000:23:42 == 1111:42:23 should be false")
	}

	// Length missmatch
	if searchFilterCmpCommunity(Community{1000, 23, 42}, Community{42, 23}) == true {
		t.Error("1000:23:42 == 42:23 should be false")
	}
}

func TestSearchFilterEqual(t *testing.T) {
	// Int values (ASNS)
	a := &SearchFilter{Value: 23}
	b := &SearchFilter{Value: 23}
	c := &SearchFilter{Value: 42}

	if a.Equal(b) == false {
		t.Error("filter[23] == filter[23] should be true")
	}

	if a.Equal(c) {
		t.Error("filter[23] == filter[42] should be false")
	}

	// String values (sources)
	a = &SearchFilter{Value: "rs-foo"}
	b = &SearchFilter{Value: "rs-foo"}
	c = &SearchFilter{Value: "rs-bar"}

	if a.Equal(b) == false {
		t.Error("filter['rs-foo'] == filter['rs-foo'] should be true")
	}

	if a.Equal(c) {
		t.Error("filter['rs-foo'] == filter['rs-bar'] should be false")
	}

	// Communities
	a = &SearchFilter{Value: Community{23, 42}}
	b = &SearchFilter{Value: Community{23, 42}}
	c = &SearchFilter{Value: Community{42, 23}}

	if a.Equal(b) == false {
		t.Error("filter[23:42] == filter[23:42] should be true")
	}

	if a.Equal(c) {
		t.Error("filter[23:42] == filter[42:23] should be false")
	}

	// Ext. Communities
	a = &SearchFilter{Value: ExtCommunity{"ro", "23", "42"}}
	b = &SearchFilter{Value: ExtCommunity{"ro", "23", "42"}}
	c = &SearchFilter{Value: ExtCommunity{"rt", "42", "23"}}

	if a.Equal(b) == false {
		t.Error("filter[ro:23:42] == filter[ro:23:42] should be true")
	}

	if a.Equal(c) {
		t.Error("filter[ro:23:42] == filter[rt:42:23] should be false")
	}

	// Large communities
	a = &SearchFilter{Value: Community{1000, 23, 42}}
	b = &SearchFilter{Value: Community{1000, 23, 42}}
	c = &SearchFilter{Value: Community{1111, 42, 23}}

	if a.Equal(b) == false {
		t.Error("filter[23:42] == filter[23:42] should be true")
	}

	if a.Equal(c) {
		t.Error("filter[23:42] == filter[42:23] should be false")
	}
}

func TestSearchFilterGroupContains(t *testing.T) {
	group := SearchFilterGroup{
		Filters: []*SearchFilter{
			&SearchFilter{Value: Community{1000, 23, 42}},
			&SearchFilter{Value: Community{1001, 24, 43}},
		},
	}

	f := &SearchFilter{Value: Community{1001, 24, 43}}
	if group.Contains(f) == false {
		t.Error("Group should contain filter.")
	}

	f = &SearchFilter{Value: Community{1111, 24, 43}}
	if group.Contains(f) {
		t.Error("Group should not contain filter.")
	}
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
	route := makeTestLookupRoute()

	// Source
	if searchFilterMatchSource(route, "3") != true {
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
	if searchFilterMatchCommunity(route, Community{111, 11}) != true {
		t.Error("Route should have community 111:11")
	}
	if searchFilterMatchCommunity(route, Community{42, 111}) == true {
		t.Error("Route should not have community 42:111")
	}

	// Ext. Communities
	if searchFilterMatchExtCommunity(route, ExtCommunity{"ro", "23", "123"}) != true {
		t.Error("Route should have community ro:23:123")
	}
	if searchFilterMatchExtCommunity(route, ExtCommunity{"rt", "42", "111"}) == true {
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
	route := makeTestLookupRoute()

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
	route := makeTestLookupRoute()

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

// Communities should match all aswell
func testSearchFilterCommunities(route Filterable, t *testing.T) {
	query := "communities=23:42,111:11"
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
		t.Error("Route should have matched filters!")
	}

	// Now check that all communities need to match
	query = "communities=23:42,111:12"
	values, err = url.ParseQuery(query)
	if err != nil {
		t.Error(err)
		return
	}

	filters, err = FiltersFromQuery(values)
	if err != nil {
		t.Error(err)
		return
	}

	if filters.MatchRoute(route) != false {
		t.Error("Route should not have matched filters!")
	}
}

func TestSearchFilterLookupRouteCommunity(t *testing.T) {
	route := makeTestLookupRoute()
	testSearchFilterCommunities(route, t)
}

// Check that ext. communities work
func testSearchFilterExtCommunities(route Filterable, t *testing.T) {
	query := "ext_communities=ro:23:123"
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
		t.Error("Route should have matched filters!")
	}

	// Now check that all communities need to match
	query = "ext_communities=ro:23:142"
	values, err = url.ParseQuery(query)
	if err != nil {
		t.Error(err)
		return
	}

	filters, err = FiltersFromQuery(values)
	if err != nil {
		t.Error(err)
		return
	}

	if filters.MatchRoute(route) != false {
		t.Error("Route should not have matched filters!")
	}
}

func TestSearchFilterRouteExtCommunities(t *testing.T) {
	route := makeTestRoute()
	testSearchFilterExtCommunities(route, t)
}

func TestSearchFilterLookupRouteExtCommunities(t *testing.T) {
	route := makeTestLookupRoute()
	testSearchFilterExtCommunities(route, t)
}

// Check that large communities work aswell
func testSearchFilterLargeCommunities(route Filterable, t *testing.T) {
	query := "large_communities=1000:23:42"
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
		t.Error("Route should have matched filters!")
	}

	// Now check that all communities need to match
	query = "large_communities=1002:111:11"
	values, err = url.ParseQuery(query)
	if err != nil {
		t.Error(err)
		return
	}

	filters, err = FiltersFromQuery(values)
	if err != nil {
		t.Error(err)
		return
	}

	if filters.MatchRoute(route) != false {
		t.Error("Route should not have matched filters!")
	}
}

func TestSearchFilterRouteLargeCommunities(t *testing.T) {
	route := makeTestRoute()
	testSearchFilterLargeCommunities(route, t)
}

func TestSearchFilterLookupRouteLargeCommunities(t *testing.T) {
	route := makeTestLookupRoute()
	testSearchFilterLargeCommunities(route, t)
}

// Subtract other
func TestSearchFiltersSub(t *testing.T) {
	query := "asns=2342,23042&communities=23:42&large_communities=42:23:42&sources=1,2,3&q=foo"
	values, err := url.ParseQuery(query)
	if err != nil {
		t.Error(err)
		return
	}

	a, err := FiltersFromQuery(values)
	if err != nil {
		t.Error(err)
		return
	}

	query = "asns=2342,10&large_communities=42:23:42&sources=1,2,3&q=foo"
	values, err = url.ParseQuery(query)
	if err != nil {
		t.Error(err)
		return
	}

	b, err := FiltersFromQuery(values)
	if err != nil {
		t.Error(err)
		return
	}

	// Modify some filters
	g := a.GetGroupByKey(SEARCH_KEY_ASNS)
	g.Filters[1].Cardinality = 9001

	t.Log(a)
	t.Log(b)

	c := a.Sub(b)

	// Check diff
	g = c.GetGroupByKey(SEARCH_KEY_ASNS)
	if len(g.Filters) != 1 {
		t.Error("There should be now only be one filter")
	}

	if g.Filters[0].Cardinality != 9001 {
		t.Error("This should be the modified filter")
	}

	if g.Filters[0].Value != 23042 {
		t.Error("This should be the modified filter")
	}

	// Should still contain community filter
	g = c.GetGroupByKey(SEARCH_KEY_COMMUNITIES)
	if len(g.Filters) != 1 {
		t.Error("The community filter should not have been touched")
	}

	// The large community filter should have been removed
	g = c.GetGroupByKey(SEARCH_KEY_LARGE_COMMUNITIES)
	if len(g.Filters) != 0 {
		t.Error("The large community filter is not removed")
	}

}

func TestSearchFiltersMergeProperties(t *testing.T) {
	filtering := NewSearchFilters()
	group := filtering.GetGroupByKey(SEARCH_KEY_ASNS)

	group.AddFilter(&SearchFilter{
		Name:  "Tech Inc. Solutions GmbH",
		Value: 23042})

	group.AddFilter(&SearchFilter{
		Name:  "Offline.net",
		Value: 1119})

	offlineNet := group.Filters[1]
	offlineNet.Cardinality = 9001

	other := NewSearchFilters()
	otherGroup := other.GetGroupByKey(SEARCH_KEY_ASNS)

	otherGroup.AddFilter(&SearchFilter{
		Value: 1119})

	other.MergeProperties(filtering)

	filter := otherGroup.Filters[0]

	if filter.Value != 1119 {
		t.Error("Expected filter value shoud be 1119, got:", filter.Value)
	}

	if filter.Cardinality < 9000 {
		t.Error("Expected cardinatlity property to be set")
	}

	if filter.Name == "" {
		t.Error("Filter name should have been merged")
	}

}

func TestNeighborFilterMatch(t *testing.T) {
	n1 := &Neighbour{
		Asn:         2342,
		Description: "Foo Networks AB",
	}
	n2 := &Neighbour{
		Asn:         42,
		Description: "Bar Communications Inc.",
	}

	filter := &NeighborFilter{
		asn: 42,
	}
	if filter.Match(n1) != false {
		t.Error("Expected n1 not to match filter")
	}
	if filter.Match(n2) == false {
		t.Error("Expected n2 to match filter")
	}

	filter = &NeighborFilter{
		name: "network",
	}
	if filter.Match(n1) == false {
		t.Error("Expected n1 to match filter")
	}
	if filter.Match(n2) != false {
		t.Error("Expected n2 not to match filter")
	}

	filter = &NeighborFilter{
		asn:  42,
		name: "network",
	}

	if filter.Match(n1) == false || filter.Match(n2) == false {
		t.Error("Expected filter to match both neighbors.")
	}
}

func TestNeighborFilterFromQuery(t *testing.T) {
	query := "asn=2342&name=foo"
	filter := NeighborFilterFromQueryString(query)

	if filter.asn != 2342 {
		t.Error("Unexpected asn filter:", filter.asn)
	}
	if filter.name != "foo" {
		t.Error("Unexpected name filter:", filter.name)
	}

	filter = NeighborFilterFromQueryString(values)
	if filter.asn != 0 {
		t.Error("Unexpected asn:", filter.asn)
	}
	if filter.name != "" {
		t.Error("Unexpected name:", filter.name)
	}
}
