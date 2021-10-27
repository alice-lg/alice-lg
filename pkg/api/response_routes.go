package api

import (
	"encoding/json"
	"time"
)

// Route is a prefix with BGP information.
type Route struct {
	ID         string `json:"id"`
	NeighborID string `json:"neighbour_id"`

	Network   string        `json:"network"`
	Interface string        `json:"interface"`
	Gateway   string        `json:"gateway"`
	Metric    int           `json:"metric"`
	BGP       *BGPInfo      `json:"bgp"`
	Age       time.Duration `json:"age"`
	Type      []string      `json:"type"` // [BGP, unicast, univ]
	Primary   bool          `json:"primary"`

	Details Details `json:"details"`
}

func (r *Route) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}

// MatchSourceID implements Filterable interface for routes
func (r *Route) MatchSourceID(id string) bool {
	return true // A route has no source info so we exclude this filter
}

// MatchASN is not defined
func (r *Route) MatchASN(asn int) bool {
	return true // Same here
}

// MatchCommunity checks for the presence of a BGP community
func (r *Route) MatchCommunity(community Community) bool {
	return r.BGP.HasCommunity(community)
}

// MatchExtCommunity checks for the presence of a BGP extended community
func (r *Route) MatchExtCommunity(community ExtCommunity) bool {
	return r.BGP.HasExtCommunity(community)
}

// MatchLargeCommunity checks for the presence of a large BGP community
func (r *Route) MatchLargeCommunity(community Community) bool {
	return r.BGP.HasLargeCommunity(community)
}

// Routes is a collection of routes
type Routes []*Route

func (routes Routes) Len() int {
	return len(routes)
}

func (routes Routes) Less(i, j int) bool {
	return routes[i].Network < routes[j].Network
}

func (routes Routes) Swap(i, j int) {
	routes[i], routes[j] = routes[j], routes[i]
}

// RoutesResponse contains all routes from a source
type RoutesResponse struct {
	Meta        *Meta  `json:api`
	Imported    Routes `json:"imported"`
	Filtered    Routes `json:"filtered"`
	NotExported Routes `json:"not_exported"`
}

// CacheTTL returns the cache ttl of the reponse
func (res *RoutesResponse) CacheTTL() time.Duration {
	now := time.Now().UTC()
	return res.Meta.TTL.Sub(now)
}

// Timed responses include the duration of the request
type Timed struct {
	RequestDuration float64 `json:"request_duration_ms"`
}

// Pagination strucutres information about the
// current page, total pages, page size, etc...
type Pagination struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// A Paginated response with pagination info
type Paginated struct {
	Pagination Pagination `json:"pagination"`
}

// Searchable responses include filters applied and available
type Searchable struct {
	FiltersAvailable *SearchFilters `json:"filters_available"`
	FiltersApplied   *SearchFilters `json:"filters_applied"`
}

// LookupRoute is a route with additional
// neighbor and state information
type LookupRoute struct {
	*Route

	State string `json:"state"` // Filtered, Imported, ...

	Neighbor    *Neighbor    `json:"neighbor"`
	RouteServer *RouteServer `json:"routeserver"`
}

// MatchSourceID implements filterable interface for lookup routes
func (r *LookupRoute) MatchSourceID(id string) bool {
	return r.RouteServer.ID == id
}

// MatchASN matches the neighbor's ASN
func (r *LookupRoute) MatchASN(asn int) bool {
	return r.Neighbor.MatchASN(asn)
}

// MatchCommunity checks for the presence of a BGP community.
func (r *LookupRoute) MatchCommunity(community Community) bool {
	return r.Route.BGP.HasCommunity(community)
}

// MatchExtCommunity matches an extended community
func (r *LookupRoute) MatchExtCommunity(community ExtCommunity) bool {
	return r.Route.BGP.HasExtCommunity(community)
}

// MatchLargeCommunity matches large communities.
func (r *LookupRoute) MatchLargeCommunity(community Community) bool {
	return r.Route.BGP.HasLargeCommunity(community)
}

// LookupRoutes is a collection of lookup routes.
type LookupRoutes []*LookupRoute

func (r LookupRoutes) Len() int {
	return len(r)
}

func (r LookupRoutes) Less(i, j int) bool {
	return r[i].Route.Network < r[j].Route.Network
}

func (r LookupRoutes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// RoutesLookupResponse is a PaginatedResponse with
// a set of lookup routes, as the result of a query of
// a specific route server.
type RoutesLookupResponse struct {
	Paginated
	Timed
	Searchable
	Routes LookupRoutes `json:"routes"`
	Meta   *Meta        `json:"api"`
}

// GlobalRoutesLookupResponse is the result of a routes
// query across all route servers.
type GlobalRoutesLookupResponse struct {
	Response
	Paginated
	Timed
	Searchable
	Routes LookupRoutes `json:"routes"`
}

// A PaginatedRoutesLookupResponse TODO
type PaginatedRoutesLookupResponse struct {
	Response
	Timed
	Searchable

	Imported *RoutesLookupResponse `json:"imported"`
	Filtered *RoutesLookupResponse `json:"filtered"`
}
