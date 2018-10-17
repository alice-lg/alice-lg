package api

import (
	"fmt"
	"net/url"
)

const (
	SEARCH_KEY_SOURCES           = "sources"
	SEARCH_KEY_ASNS              = "asns"
	SEARCH_KEY_COMMUNITIES       = "communities"
	SEARCH_KEY_EXT_COMMUNITIES   = "ext_communities"
	SEARCH_KEY_LARGE_COMMUNITIES = "large_communities"
)

/*
API Search

* Helper methods for searching
* Handle filter criteria

*/
type SearchFilter struct {
	Cardinality int         `json:"cardinality"`
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
}

type SearchFilterGroup struct {
	Key string `json:"key"`

	Filters    []*SearchFilter `json:"filters"`
	filtersIdx map[string]int
}

type SearchFilters []*SearchFilterGroup

func NewSearchFilters() *SearchFilters {
	// Define groups: CAVEAT! the order is relevant
	groups := &SearchFilters{
		&SearchFilterGroup{
			Key:        SEARCH_KEY_SOURCES,
			Filters:    []*SearchFilter{},
			filtersIdx: make(map[string]int),
		},
		&SearchFilterGroup{
			Key:        SEARCH_KEY_ASNS,
			Filters:    []*SearchFilter{},
			filtersIdx: make(map[string]int),
		},
		&SearchFilterGroup{
			Key:        SEARCH_KEY_COMMUNITIES,
			Filters:    []*SearchFilter{},
			filtersIdx: make(map[string]int),
		},
		&SearchFilterGroup{
			Key:        SEARCH_KEY_EXT_COMMUNITIES,
			Filters:    []*SearchFilter{},
			filtersIdx: make(map[string]int),
		},
		&SearchFilterGroup{
			Key:        SEARCH_KEY_LARGE_COMMUNITIES,
			Filters:    []*SearchFilter{},
			filtersIdx: make(map[string]int),
		},
	}

	return groups
}

func (self *SearchFilters) GetGroupByKey(key string) *SearchFilterGroup {
	// This is an optimization (this is basically a fixed hash map,
	// with hash(key) = position(key)
	switch key {
	case SEARCH_KEY_SOURCES:
		return (*self)[0]
	case SEARCH_KEY_ASNS:
		return (*self)[1]
	case SEARCH_KEY_COMMUNITIES:
		return (*self)[2]
	case SEARCH_KEY_EXT_COMMUNITIES:
		return (*self)[3]
	case SEARCH_KEY_LARGE_COMMUNITIES:
		return (*self)[4]
	}
	return nil
}

func (self *SearchFilterGroup) GetFilterByValue(value interface{}) *SearchFilter {
	// I've tried it with .(fmt.Stringer), but int does not implement this...
	// So whatever. I'm using the trick of letting Sprintf choose the right
	// conversion. If this is too expensive, we need to refactor this.
	// TODO: profile this.
	idx, ok := self.filtersIdx[fmt.Sprintf("%v", value)]
	if !ok {
		return nil // We don't have this particular filter
	}

	return self.Filters[idx]
}

func (self *SearchFilterGroup) AddFilter(filter *SearchFilter) {
	// Check if a filter with this value is present, if not:
	// append and update index; otherwise incrementc cardinality
	if presentFilter := self.GetFilterByValue(filter.Value); presentFilter != nil {
		presentFilter.Cardinality++
		return
	}

	// Insert filter
	idx := len(self.Filters)
	filter.Cardinality = 1
	self.Filters = append(self.Filters, filter)
	self.filtersIdx[fmt.Sprintf("%v", filter.Value)] = idx
}

func (self *SearchFilterGroup) AddFilters(filters []*SearchFilter) {
	for _, filter := range filters {
		self.AddFilter(filter)
	}
}

/*
 Update filter struct to include route:
  - Extract ASN, source, bgp communites,
  - Find Filter in group, increment result count if required.
*/
func (self *SearchFilters) UpdateFromRoute(route *LookupRoute) {
	// Add source
	self.GetGroupByKey(SEARCH_KEY_SOURCES).AddFilter(&SearchFilter{
		Name:  route.Routeserver.Name,
		Value: route.Routeserver.Id,
	})

	// Add ASN from neighbor
	self.GetGroupByKey(SEARCH_KEY_ASNS).AddFilter(&SearchFilter{
		Name:  route.Neighbour.Description,
		Value: route.Neighbour.Asn,
	})

	// Add communities
	communities := self.GetGroupByKey(SEARCH_KEY_COMMUNITIES)
	for _, c := range route.Bgp.Communities {
		communities.AddFilter(&SearchFilter{
			Name:  c.String(),
			Value: c,
		})
	}
	extCommunities := self.GetGroupByKey(SEARCH_KEY_EXT_COMMUNITIES)
	for _, c := range route.Bgp.ExtCommunities {
		extCommunities.AddFilter(&SearchFilter{
			Name:  c.String(),
			Value: c,
		})
	}
	largeCommunities := self.GetGroupByKey(SEARCH_KEY_LARGE_COMMUNITIES)
	for _, c := range route.Bgp.LargeCommunities {
		largeCommunities.AddFilter(&SearchFilter{
			Name:  c.String(),
			Value: c,
		})
	}
}

/*
 Build filter struct from query params:
 For example a query string of:
    asns=2342,23123&communities=23:42&large_communities=23:42:42
 yields a filtering struct of
    Groups[
        Group{"sources", []},
        Group{"asns", [Filter{Value: 2342},
                       Filter{Value: 23123}]},
        Group{"communities", ...
    }

*/
func FiltersFromQuery(query url.Values) (*SearchFilters, error) {
	queryFilters := NewSearchFilters()
	for key, _ := range query {
		value := query.Get(key)
		switch key {
		case SEARCH_KEY_SOURCES:
			filters, err := parseQueryValueList(parseIntValue, value)
			if err != nil {
				return nil, err
			}
			queryFilters.GetGroupByKey(SEARCH_KEY_SOURCES).AddFilters(filters)
			break

		case SEARCH_KEY_ASNS:
			filters, err := parseQueryValueList(parseIntValue, value)
			if err != nil {
				return nil, err
			}
			queryFilters.GetGroupByKey(SEARCH_KEY_ASNS).AddFilters(filters)
			break

		case SEARCH_KEY_COMMUNITIES:
			filters, err := parseQueryValueList(parseCommunityValue, value)
			if err != nil {
				return nil, err
			}
			queryFilters.GetGroupByKey(SEARCH_KEY_COMMUNITIES).AddFilters(filters)
			break

		case SEARCH_KEY_EXT_COMMUNITIES:
			filters, err := parseQueryValueList(parseExtCommunityValue, value)
			if err != nil {
				return nil, err
			}
			queryFilters.GetGroupByKey(SEARCH_KEY_EXT_COMMUNITIES).AddFilters(filters)
			break

		case SEARCH_KEY_LARGE_COMMUNITIES:
			filters, err := parseQueryValueList(parseCommunityValue, value)
			if err != nil {
				return nil, err
			}
			queryFilters.GetGroupByKey(SEARCH_KEY_LARGE_COMMUNITIES).AddFilters(filters)
			break
		}
	}

	return queryFilters, nil
}

/*
 Match a route. Check if route matches all filters.
 Unless all filters are blank.
*/
func (self *SearchFilters) MatchRoute(route LookupRoute) bool {
	return false
}
