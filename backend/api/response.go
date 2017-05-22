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

	// Original response
	Details map[string]interface{} `json:"details"`
}

type NeighboursResponse struct {
	Api        ApiStatus   `json:"api"`
	Neighbours []Neighbour `json:"neighbours"`
}

// BGP
type Community []int

type BgpInfo struct {
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

	Filtered       bool   `json:"filtered"`
	FilteredReason string `json:"filtered_reason"`

	Details Details `json:"details"`
}

type RoutesResponse struct {
	Api      ApiStatus `json:"api"`
	Exported []Route   `json:"exported"`
	Filtered []Route   `json:"filtered"`
}
