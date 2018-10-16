package api

import ()

const (
	API_SEARCH_KEY_SOURCES           = "sources"
	API_SEARCH_KEY_ASNS              = "asns"
	API_SEARCH_KEY_COMMUNITIES       = "communities"
	API_SEARCH_KEY_EXT_COMMUNITIES   = "ext_communities"
	API_SEARCH_KEY_LARGE_COMMUNITIES = "large_communities"
)

/*
API Search

* Helper methods for searching
* Handle filter criteria

*/
type ApiSearchFilter struct {
	Cardinality int         `json:"cardinality"`
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
}

type ApiSearchFilterGroup struct {
	Key string `json:"key"`

	Filters    []*ApiSearchFilter `json:"filters"`
	filtersIdx map[interface{}]int
}

type ApiSearchFilters []*ApiSearchFilterGroup

func NewApiSearchFilters() *ApiSearchFilters {
	// Define groups: CAVEAT! the order is relevant
	groups := &ApiSearchFilters{
		&ApiSearchFilterGroup{
			Key:        API_SEARCH_KEY_SOURCES,
			Filters:    []*ApiSearchFilter{},
			filtersIdx: make(map[interface{}]int),
		},
		&ApiSearchFilterGroup{
			Key:        API_SEARCH_KEY_ASNS,
			Filters:    []*ApiSearchFilter{},
			filtersIdx: make(map[interface{}]int),
		},
		&ApiSearchFilterGroup{
			Key:        API_SEARCH_KEY_COMMUNITIES,
			Filters:    []*ApiSearchFilter{},
			filtersIdx: make(map[interface{}]int),
		},
		&ApiSearchFilterGroup{
			Key:        API_SEARCH_KEY_EXT_COMMUNITIES,
			Filters:    []*ApiSearchFilter{},
			filtersIdx: make(map[interface{}]int),
		},
		&ApiSearchFilterGroup{
			Key:        API_SEARCH_KEY_LARGE_COMMUNITIES,
			Filters:    []*ApiSearchFilter{},
			filtersIdx: make(map[interface{}]int),
		},
	}

	return groups
}

func (self *ApiSearchFilters) GetGroupByKey(key string) *ApiSearchFilterGroup {
	// This is an optimization (this is basically a fixed hash map,
	// with hash(key) = position(key)
	switch key {
	case API_SEARCH_KEY_SOURCES:
		return (*self)[0]
	case API_SEARCH_KEY_ASNS:
		return (*self)[1]
	case API_SEARCH_KEY_COMMUNITIES:
		return (*self)[2]
	case API_SEARCH_KEY_EXT_COMMUNITIES:
		return (*self)[3]
	case API_SEARCH_KEY_LARGE_COMMUNITIES:
		return (*self)[4]
	}
	return nil
}

func (self *ApiSearchFilterGroup) GetFilterByValue(value interface{}) *ApiSearchFilter {
	idx, ok := self.filtersIdx[value]
	if !ok {
		return nil // We don't have this particular filter
	}

	return self.Filters[idx]
}

func (self *ApiSearchFilterGroup) AddFilter(filter *ApiSearchFilter) {
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
	self.filtersIdx[filter.Value] = idx
}

/*
 Update filter struct to include route:
  - Extract ASN, source, bgp communites,
  - Find Filter in group, increment result count if required.
*/
func (self *ApiSearchFilters) UpdateFromRoute(route LookupRoute) {
	// Add source
	self.GetGroupByKey(API_SEARCH_KEY_SOURCES).AddFilter(&ApiSearchFilter{
		Name:  route.Routeserver.Name,
		Value: route.Routeserver.Id,
	})

	// Add ASN from neighbor
	self.GetGroupByKey(API_SEARCH_KEY_ASNS).AddFilter(&ApiSearchFilter{
		Name:  route.Neighbour.Description,
		Value: route.Neighbour.Asn,
	})

	// Add communities
	communities := self.GetGroupByKey(API_SEARCH_KEY_COMMUNITIES)
	for _, c := range route.Bgp.Communities {
		communities.AddFilter(&ApiSearchFilter{
			Name:  c.String(),
			Value: c,
		})
	}
	extCommunities := self.GetGroupByKey(API_SEARCH_KEY_EXT_COMMUNITIES)
	for _, c := range route.Bgp.ExtCommunities {
		extCommunities.AddFilter(&ApiSearchFilter{
			Name:  c.String(),
			Value: c,
		})
	}
	largeCommunities := self.GetGroupByKey(API_SEARCH_KEY_LARGE_COMMUNITIES)
	for _, c := range route.Bgp.LargeCommunities {
		largeCommunities.AddFilter(&ApiSearchFilter{
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
func (self *ApiSearchFilters) UpdateFromQuery(query string) {

}

func (self *ApiSearchFilters) MatchRoute(route LookupRoute) {

}
