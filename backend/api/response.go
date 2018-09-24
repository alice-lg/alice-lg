package api

import (
	"time"
)

// General api response
type Response interface{}

// Details, usually the original backend response
type Details map[string]interface{}

// Error Handling
type ErrorResponse struct {
	Message       string `json:"message"`
	Code          int    `json:"code"`
	Tag           string `json:"tag"`
	RouteserverId int    `json:"routeserver_id"`
}

// Cache aware api response
type CacheableResponse interface {
	CacheTtl() time.Duration
}

// Config
type ConfigResponse struct {
	Rejection     Rejection         `json:"rejection"`
	RejectReasons map[string]string `json:"reject_reasons"`

	Noexport        Noexport          `json:"noexport"`
	NoexportReasons map[string]string `json:"noexport_reasons"`

	Rpki Rpki `json:"rpki"`

	BgpCommunities map[string]interface{} `json:"bgp_communities"`

	NeighboursColumns      map[string]string `json:"neighbours_columns"`
	NeighboursColumnsOrder []string          `json:"neighbours_columns_order"`

	RoutesColumns      map[string]string `json:"routes_columns"`
	RoutesColumnsOrder []string          `json:"routes_columns_order"`

	LookupColumns      map[string]string `json:"lookup_columns"`
	LookupColumnsOrder []string          `json:"lookup_columns_order"`

	PrefixLookupEnabled bool `json:"prefix_lookup_enabled"`
}

type Rejection struct {
	Asn      int `json:"asn"`
	RejectId int `json:"reject_id"`
}

type Noexport struct {
	Asn          int  `json:"asn"`
	NoexportId   int  `json:"noexport_id"`
	LoadOnDemand bool `json:"load_on_demand"`
}

type Rpki struct {
	Enabled    bool     `json:"enabled"`
	Valid      []string `json:"valid"`
	Unknown    []string `json:"unknown"`
	NotChecked []string `json:"not_checked"`
	Invalid    []string `json:"invalid"`
}

// Status
type ApiStatus struct {
	Version         string      `json:"version"`
	CacheStatus     CacheStatus `json:"cache_status"`
	ResultFromCache bool        `json:"result_from_cache"`
	Ttl             time.Time   `json:"ttl"`
}

type CacheStatus struct {
	CachedAt time.Time `json:"cached_at"`
	OrigTtl  int       `json:"orig_ttl"`
}

type Status struct {
	ServerTime   time.Time `json:"server_time"`
	LastReboot   time.Time `json:"last_reboot"`
	LastReconfig time.Time `json:"last_reconfig"`
	Message      string    `json:"message"`
	RouterId     string    `json:"router_id"`
	Version      string    `json:"version"`
	Backend      string    `json:"backend"`
}

type StatusResponse struct {
	Api    ApiStatus `json:"api"`
	Status Status    `json:"status"`
}

// Routeservers
type Routeserver struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Asn        int      `json:"asn"`
	Blackholes []string `json:"blackholes"`
}

type RouteserversResponse struct {
	Routeservers []Routeserver `json:"routeservers"`
}

// Neighbours
type Neighbours []*Neighbour

type Neighbour struct {
	Id string `json:"id"`

	// Mandatory fields
	Address            string        `json:"address"`
	Asn                int           `json:"asn"`
	State              string        `json:"state"`
	Description        string        `json:"description"`
	RoutesReceived     int           `json:"routes_received"`
	RoutesFiltered     int           `json:"routes_filtered"`
	RoutesExported     int           `json:"routes_exported"`
	RoutesPreferred    int           `json:"routes_preferred"`
	RoutesAccepted     int           `json:"routes_accepted"`
	RoutesPipeFiltered int           `json:"routes_pipe_filtered"`
	Uptime             time.Duration `json:"uptime"`
	LastError          string        `json:"last_error"`

	// Original response
	Details map[string]interface{} `json:"details"`
}

// Implement sorting interface for routes
func (neighbours Neighbours) Len() int {
	return len(neighbours)
}

func (neighbours Neighbours) Less(i, j int) bool {
	return neighbours[i].Asn < neighbours[j].Asn
}

func (neighbours Neighbours) Swap(i, j int) {
	neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
}

type NeighboursResponse struct {
	Api        ApiStatus  `json:"api"`
	Neighbours Neighbours `json:"neighbours"`
}

// Neighbours response is cacheable
func (self *NeighboursResponse) CacheTtl() time.Duration {
	now := time.Now().UTC()
	return self.Api.Ttl.Sub(now)
}

type NeighboursLookupResults map[int]Neighbours

// BGP
type Community []int

type BgpInfo struct {
	Origin           string      `json:"origin"`
	AsPath           []int       `json:"as_path"`
	NextHop          string      `json:"next_hop"`
	Communities      []Community `json:"communities"`
	LargeCommunities []Community `json:"large_communities"`
	LocalPref        int         `json:"local_pref"`
	Med              int         `json:"med"`
}

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

type PaginatedRoutesResponse struct {
	*RoutesResponse
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
	*TimedResponse

	Imported *LookupRoutesResponse `json:"imported"`
	Filtered *LookupRoutesResponse `json:"filtered"`
}
