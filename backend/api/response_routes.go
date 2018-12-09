package api

import (
	"time"
)

// Prefixes
type Route struct {
	Id          string `json:"id"`
	NeighbourId string `json:"neighbour_id"`

	Network   string        `json:"network"`
	Interface string        `json:"interface"`
	Gateway   string        `json:"gateway"`
	Metric    int           `json:"metric"`
	Bgp       BgpInfo       `json:"bgp"`
	Age       time.Duration `json:"age"`
	Type      []string      `json:"type"` // [BGP, unicast, univ]
	Primary   bool          `json:"primary"`

	Details Details `json:"details"`
}

// Implement Filterable interface for routes
func (self *Route) MatchSourceId(id string) bool {
	return true // A route has no source info so we exclude this filter
}

func (self *Route) MatchAsn(asn int) bool {
	return true // Same here
}

// Only community filters are interesting at this point:
func (self *Route) MatchCommunity(community Community) bool {
	return self.Bgp.HasCommunity(community)
}

func (self *Route) MatchExtCommunity(community ExtCommunity) bool {
	return self.Bgp.HasExtCommunity(community)
}

func (self *Route) MatchLargeCommunity(community Community) bool {
	return self.Bgp.HasLargeCommunity(community)
}

type Routes []*Route

// Implement sorting interface for routes
func (routes Routes) Len() int {
	return len(routes)
}

func (routes Routes) Less(i, j int) bool {
	return routes[i].Network < routes[j].Network
}

func (routes Routes) Swap(i, j int) {
	routes[i], routes[j] = routes[j], routes[i]
}

type RoutesResponse struct {
	Api         ApiStatus `json:"api"`
	Imported    Routes    `json:"imported"`
	Filtered    Routes    `json:"filtered"`
	NotExported Routes    `json:"not_exported"`
}

func (self *RoutesResponse) CacheTtl() time.Duration {
	now := time.Now().UTC()
	return self.Api.Ttl.Sub(now)
}

type TimedResponse struct {
	RequestDuration float64 `json:"request_duration_ms"`
}

type Pagination struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type PaginatedResponse struct {
	Pagination Pagination `json:"pagination"`
}

type FilterableResponse struct {
	FiltersAvailable *SearchFilters `json:"filters_available"`
	FiltersApplied   *SearchFilters `json:"filters_applied"`
}

type PaginatedRoutesResponse struct {
	*RoutesResponse
	TimedResponse
	FilterableResponse
	Pagination Pagination `json:"pagination"`
}

// Lookup Prefixes
type LookupRoute struct {
	Id          string     `json:"id"`
	NeighbourId string     `json:"neighbour_id"`
	Neighbour   *Neighbour `json:"neighbour"`

	State string `json:"state"` // Filtered, Imported, ...

	Routeserver Routeserver `json:"routeserver"`

	Network   string        `json:"network"`
	Interface string        `json:"interface"`
	Gateway   string        `json:"gateway"`
	Metric    int           `json:"metric"`
	Bgp       BgpInfo       `json:"bgp"`
	Age       time.Duration `json:"age"`
	Type      []string      `json:"type"` // [BGP, unicast, univ]
	Primary   bool          `json:"primary"`

	Details Details `json:"details"`
}

// Implement Filterable interface for lookup routes
func (self *LookupRoute) MatchSourceId(id string) bool {
	return self.Routeserver.Id == id
}

func (self *LookupRoute) MatchAsn(asn int) bool {
	return self.Neighbour.Asn == asn
}

// Only community filters are interesting at this point:
func (self *LookupRoute) MatchCommunity(community Community) bool {
	return self.Bgp.HasCommunity(community)
}

func (self *LookupRoute) MatchExtCommunity(community ExtCommunity) bool {
	return self.Bgp.HasExtCommunity(community)
}

func (self *LookupRoute) MatchLargeCommunity(community Community) bool {
	return self.Bgp.HasLargeCommunity(community)
}

// Implement sorting interface for lookup routes
func (routes LookupRoutes) Len() int {
	return len(routes)
}

func (routes LookupRoutes) Less(i, j int) bool {
	return routes[i].Network < routes[j].Network
}

func (routes LookupRoutes) Swap(i, j int) {
	routes[i], routes[j] = routes[j], routes[i]
}

type LookupRoutes []*LookupRoute

// TODO: Naming is a bit yuck
type LookupRoutesResponse struct {
	*PaginatedResponse
	Routes LookupRoutes `json:"routes"`
}

// TODO: Refactor this (might be legacy)
type RoutesLookupResponse struct {
	Api    ApiStatus    `json:"api"`
	Routes LookupRoutes `json:"routes"`
}

type RoutesLookupResponseGlobal struct {
	Routes LookupRoutes `json:"routes"`

	// Pagination
	TotalRoutes int `json:"total_routes"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`

	// Meta
	Time float64 `json:"query_duration_ms"`
}

type PaginatedRoutesLookupResponse struct {
	TimedResponse
	FilterableResponse

	Api ApiStatus `json:"api"` // Add to provide cache status information

	Imported *LookupRoutesResponse `json:"imported"`
	Filtered *LookupRoutesResponse `json:"filtered"`
}
