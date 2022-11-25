package api

import (
	"encoding/json"
	"log"
	"time"
)

// Route is a prefix with BGP information.
type Route struct {
	ID         string  `json:"id"`
	NeighborID *string `json:"neighbor_id"`

	Network    string        `json:"network"`
	Interface  *string       `json:"interface"`
	Gateway    *string       `json:"gateway"`
	Metric     int           `json:"metric"`
	BGP        *BGPInfo      `json:"bgp"`
	Age        time.Duration `json:"age"`
	Type       []string      `json:"type"` // [BGP, unicast, univ]
	Primary    bool          `json:"primary"`
	LearntFrom *string       `json:"learnt_from"`

	Details *json.RawMessage `json:"details"`
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

// ToLookupRoutes prepares routes for lookup
func (routes Routes) ToLookupRoutes(
	state string,
	rs *RouteServer,
	neighbors map[string]*Neighbor,
) LookupRoutes {
	lookupRoutes := make(LookupRoutes, 0, len(routes))
	for _, route := range routes {
		neighbor, ok := neighbors[*route.NeighborID]
		if !ok {
			log.Println("prepare route, neighbor not found:", route.NeighborID)
			continue
		}
		lr := &LookupRoute{
			Route:       route,
			State:       state,
			Neighbor:    neighbor,
			RouteServer: rs,
		}
		lr.Route.Details = nil
		lookupRoutes = append(lookupRoutes, lr)
	}
	return lookupRoutes
}

// RoutesResponse contains all routes from a source
type RoutesResponse struct {
	Response
	Imported    Routes `json:"imported"`
	Filtered    Routes `json:"filtered"`
	NotExported Routes `json:"not_exported"`
}

// CacheTTL returns the cache ttl of the response
func (res *RoutesResponse) CacheTTL() time.Duration {
	now := time.Now().UTC()
	return res.Response.Meta.TTL.Sub(now)
}

// TimedResponse include the duration of the request
type TimedResponse struct {
	RequestDuration float64 `json:"request_duration_ms"`
}

// Pagination information, including the
// current page, total pages, page size, etc...
type Pagination struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// A PaginatedResponse with pagination info
type PaginatedResponse struct {
	Pagination Pagination `json:"pagination"`
}

// FilteredResponse includes filters applied and available
type FilteredResponse struct {
	FiltersAvailable *SearchFilters `json:"filters_available"`
	FiltersApplied   *SearchFilters `json:"filters_applied"`
}

const (
	// RouteStateFiltered indicates that the route
	// was not accepted by the route server.
	RouteStateFiltered = "filtered"
	// RouteStateImported indicates that the route was
	// imported by the route server.
	RouteStateImported = "imported"
)

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

// RoutesLookup contains routes and pagination info
type RoutesLookup struct {
	Routes     LookupRoutes `json:"routes"`
	Pagination Pagination   `json:"pagination"`
}

// RoutesLookupResponse is a PaginatedResponse with
// a set of lookup routes, as the result of a query of
// a specific route server.
type RoutesLookupResponse struct {
	Response
	PaginatedResponse
	TimedResponse
	FilteredResponse
	Routes LookupRoutes `json:"routes"`
}

// GlobalRoutesLookupResponse is the result of a routes
// query across all route servers.
type GlobalRoutesLookupResponse struct {
	Response
	PaginatedResponse
	TimedResponse
	FilteredResponse
	Routes LookupRoutes `json:"routes"`
}

// A PaginatedRoutesResponse includes routes and pagination
// information form a single route server
type PaginatedRoutesResponse struct {
	Response
	PaginatedResponse
	TimedResponse
	FilteredResponse
	RoutesResponse
}

// A PaginatedRoutesLookupResponse TODO
type PaginatedRoutesLookupResponse struct {
	Response
	TimedResponse
	FilteredResponse

	Imported *RoutesLookup `json:"imported"`
	Filtered *RoutesLookup `json:"filtered"`

	Status *StoreStatusMeta `json:"status"`
}
