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
	Error string `json:"error"`
}

// Config
type ConfigResponse struct {
	Rejection     Rejection      `json:"rejection"`
	RejectReasons map[int]string `json:"reject_reasons"`

	Noexport        Noexport       `json:"noexport"`
	NoexportReasons map[int]string `json:"noexport_reasons"`

	RoutesColumns map[string]string `json:"routes_columns"`
}

type Rejection struct {
	Asn      int `json:"asn"`
	RejectId int `json:"reject_id"`
}

type Noexport struct {
	Asn        int `json:"asn"`
	NoexportId int `json:"noexport_id"`
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
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RouteserversResponse struct {
	Routeservers []Routeserver `json:"routeservers"`
}

// Neighbours
type Neighbours []Neighbour

type Neighbour struct {
	Id string `json:"id"`

	// Mandatory fields
	Address         string        `json:"address"`
	Asn             int           `json:"asn"`
	State           string        `json:"state"`
	Description     string        `json:"description"`
	RoutesReceived  int           `json:"routes_received"`
	RoutesFiltered  int           `json:"routes_filtered"`
	RoutesExported  int           `json:"routes_exported"`
	RoutesPreferred int           `json:"routes_preferred"`
	Uptime          time.Duration `json:"uptime"`
	LastError       string        `json:"last_error"`

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

type NeighboursLookupResults map[int][]Neighbour

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

	Details Details `json:"details"`
}

// Lookup Prefixes
type LookupRoute struct {
	Id          string    `json:"id"`
	NeighbourId string    `json:"neighbour_id"`
	Neighbour   Neighbour `json:"neighbour"`

	State string `json:"state"` // Filtered, Imported, ...

	Routeserver Routeserver `json:"routeserver"`

	Network   string        `json:"network"`
	Interface string        `json:"interface"`
	Gateway   string        `json:"gateway"`
	Metric    int           `json:"metric"`
	Bgp       BgpInfo       `json:"bgp"`
	Age       time.Duration `json:"age"`
	Type      []string      `json:"type"` // [BGP, unicast, univ]

	Details Details `json:"details"`
}

type Routes []Route

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
	Imported    []Route   `json:"imported"`
	Filtered    []Route   `json:"filtered"`
	NotExported []Route   `json:"not_exported"`
}

type RoutesLookupResponse struct {
	Api    ApiStatus     `json:"api"`
	Routes []LookupRoute `json:"routes"`
}

type RoutesLookupResponseGlobal struct {
	Routes []LookupRoute `json:"routes"`

	// Pagination
	TotalRoutes int `json:"total_routes"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`

	// Meta
	Time float64 `json:"query_duration_ms"`
}
