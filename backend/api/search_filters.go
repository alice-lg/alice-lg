package api

import (
	"fmt"
	"log"
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
type Filterable interface {
	MatchSourceId(sourceId int) bool
	MatchAsn(asn int) bool
	MatchCommunity(community Community) bool
	MatchExtCommunity(community ExtCommunity) bool
	MatchLargeCommunity(community Community) bool
}

type FilterValue interface{}

type SearchFilter struct {
	Cardinality int         `json:"cardinality"`
	Name        string      `json:"name"`
	Value       FilterValue `json:"value"`
}

type SearchFilterCmpFunc func(a FilterValue, b FilterValue) bool

func searchFilterCmpInt(a FilterValue, b FilterValue) bool {
	return a.(int) == b.(int)
}

func searchFilterCmpCommunity(a FilterValue, b FilterValue) bool {
	ca := a.(Community)
	cb := b.(Community)

	if len(ca) != len(cb) {
		return false
	}

	// Compare components
	for i, _ := range ca {
		if ca[i] != cb[i] {
			return false
		}
	}

	return true
}

func searchFilterCmpExtCommunity(a FilterValue, b FilterValue) bool {
	ca := a.(ExtCommunity)
	cb := b.(ExtCommunity)

	if len(ca) != len(cb) || len(ca) != 3 || len(cb) != 3 {
		return false
	}

	return ca[0] == cb[0] && ca[1] == cb[1] && ca[2] == cb[2]
}

func (self *SearchFilter) Equal(other *SearchFilter) bool {
	var cmp SearchFilterCmpFunc
	switch other.Value.(type) {
	case Community:
		cmp = searchFilterCmpCommunity
		break
	case ExtCommunity:
		cmp = searchFilterCmpExtCommunity
		break
	case int:
		cmp = searchFilterCmpInt
	}

	if cmp == nil {
		log.Println("Unknown search filter value type")
		return false
	}

	return cmp(self.Value, other.Value)
}

/*
 Search Filter Groups
*/

type SearchFilterGroup struct {
	Key string `json:"key"`

	Filters    []*SearchFilter `json:"filters"`
	filtersIdx map[string]int
}

func (self *SearchFilterGroup) FindFilter(filter *SearchFilter) *SearchFilter {
	for _, f := range self.Filters {
		if f.Equal(filter) == true {
			return f
		}
	}
	return nil
}

func (self *SearchFilterGroup) Contains(filter *SearchFilter) bool {
	return self.FindFilter(filter) != nil
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

func (self *SearchFilterGroup) rebuildIndex() {
	self.filtersIdx = map[string]int{}

	for i, filter := range self.Filters {
		self.filtersIdx[fmt.Sprintf("%v", filter.Value)] = i
	}
}

/*
 Search comparators
*/
type SearchFilterComparator func(route Filterable, value interface{}) bool

func searchFilterMatchSource(route Filterable, value interface{}) bool {
	sourceId, ok := value.(int)
	if !ok {
		return false
	}
	return route.MatchSourceId(sourceId)
}

func searchFilterMatchAsn(route Filterable, value interface{}) bool {
	asn, ok := value.(int)
	if !ok {
		return false
	}

	return route.MatchAsn(asn)
}

func searchFilterMatchCommunity(route Filterable, value interface{}) bool {
	community, ok := value.(Community)
	if !ok {
		return false
	}
	return route.MatchCommunity(community)
}

func searchFilterMatchExtCommunity(route Filterable, value interface{}) bool {
	community, ok := value.(ExtCommunity)
	if !ok {
		return false
	}
	return route.MatchExtCommunity(community)
}

func searchFilterMatchLargeCommunity(route Filterable, value interface{}) bool {
	community, ok := value.(Community)
	if !ok {
		return false
	}
	return route.MatchLargeCommunity(community)
}

func selectCmpFuncByKey(key string) SearchFilterComparator {
	var cmp SearchFilterComparator
	switch key {
	case SEARCH_KEY_SOURCES:
		cmp = searchFilterMatchSource
		break
	case SEARCH_KEY_ASNS:
		cmp = searchFilterMatchAsn
		break
	case SEARCH_KEY_COMMUNITIES:
		cmp = searchFilterMatchCommunity
		break
	case SEARCH_KEY_EXT_COMMUNITIES:
		cmp = searchFilterMatchExtCommunity
		break
	case SEARCH_KEY_LARGE_COMMUNITIES:
		cmp = searchFilterMatchLargeCommunity
		break
	default:
		cmp = nil
	}

	return cmp
}

func (self *SearchFilterGroup) MatchAny(route Filterable) bool {
	// Check if we have any filter to match
	if len(self.Filters) == 0 {
		return true // no filter, everything matches
	}

	// Get comparator
	cmp := selectCmpFuncByKey(self.Key)
	if cmp == nil {
		return false // This should not have happened!
	}

	// Check if any of the given filters matches
	for _, filter := range self.Filters {
		if cmp(route, filter.Value) {
			return true
		}
	}

	return false
}

func (self *SearchFilterGroup) MatchAll(route Filterable) bool {
	// Check if we have any filter to match
	if len(self.Filters) == 0 {
		return true // no filter, everything matches. Like above.
	}

	// Get comparator
	cmp := selectCmpFuncByKey(self.Key)
	if cmp == nil {
		return false // This again should not have happened!
	}

	// Assert that all filters match.
	for _, filter := range self.Filters {
		if !cmp(route, filter.Value) {
			return false
		}
	}

	// Everythings fine.
	return true
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

/*
 Update filter struct to include route:
  - Extract ASN, source, bgp communites,
  - Find Filter in group, increment result count if required.
*/
func (self *SearchFilters) UpdateFromLookupRoute(route *LookupRoute) {
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

// This is the same as above, but only the communities
// are considered.
func (self *SearchFilters) UpdateFromRoute(route *Route) {

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
func (self *SearchFilters) MatchRoute(route Filterable) bool {
	sources := self.GetGroupByKey(SEARCH_KEY_SOURCES)
	if !sources.MatchAny(route) {
		return false
	}

	asns := self.GetGroupByKey(SEARCH_KEY_ASNS)
	if !asns.MatchAny(route) {
		return false
	}

	communities := self.GetGroupByKey(SEARCH_KEY_COMMUNITIES)
	if !communities.MatchAll(route) {
		return false
	}

	extCommunities := self.GetGroupByKey(SEARCH_KEY_EXT_COMMUNITIES)
	if !extCommunities.MatchAll(route) {
		return false
	}

	largeCommunities := self.GetGroupByKey(SEARCH_KEY_LARGE_COMMUNITIES)
	if !largeCommunities.MatchAll(route) {
		return false
	}

	return true
}

func (self *SearchFilters) Sub(other *SearchFilters) *SearchFilters {
	result := make(SearchFilters, len(*self))

	for id, group := range *self {
		otherGroup := (*other)[id]
		diff := &SearchFilterGroup{
			Key:     group.Key,
			Filters: []*SearchFilter{},
		}

		// Combine filters
		for _, f := range group.Filters {
			if otherGroup.Contains(f) {
				continue // Let's skip this
			}
			diff.Filters = append(diff.Filters, f)
		}

		diff.rebuildIndex()
		result[id] = diff
	}

	return &result
}
