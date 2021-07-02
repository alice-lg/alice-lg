package api

import (
	"fmt"
	"time"
)

// Response is a general API response
type Response interface{}

// Details are usually the original backend response
type Details map[string]interface{}

// ErrorResponse encodes an error message and code
type ErrorResponse struct {
	Message       string `json:"message"`
	Code          int    `json:"code"`
	Tag           string `json:"tag"`
	RouteserverId string `json:"routeserver_id"`
}

// CacheableResponse is a cache aware API response
type CacheableResponse interface {
	CacheTTL() time.Duration
}

// ConfigResponse is a response with client runtime configuration
type ConfigResponse struct {
	Asn int `json:"asn"`

	RejectReasons map[string]interface{} `json:"reject_reasons"`

	Noexport        Noexport               `json:"noexport"`
	NoexportReasons map[string]interface{} `json:"noexport_reasons"`

	RejectCandidates RejectCandidates `json:"reject_candidates"`

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

// Noexport options
type Noexport struct {
	LoadOnDemand bool `json:"load_on_demand"`
}

type RejectCandidates struct {
	Communities map[string]interface{} `json:"communities"`
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

// A Routeserver is a datasource with attributes.
type Routeserver struct {
	Id         string   `json:"id"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Group      string   `json:"group"`
	Blackholes []string `json:"blackholes"`

	Order int `json:"-"`
}

type Routeservers []Routeserver

// Implement sorting interface for routeservers
func (rs Routeservers) Len() int {
	return len(rs)
}

func (rs Routeservers) Less(i, j int) bool {
	return rs[i].Order < rs[j].Order
}

func (rs Routeservers) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

type RouteserversResponse struct {
	Routeservers []Routeserver `json:"routeservers"`
}

// BGP
type Community []int

func (com Community) String() string {
	res := ""
	for _, v := range com {
		res += fmt.Sprintf(":%d", v)
	}
	return res[1:]
}

type Communities []Community

/*
Deduplicate communities
*/
func (communities Communities) Unique() Communities {
	seen := map[string]bool{}
	result := make(Communities, 0, len(communities))

	for _, com := range communities {
		key := com.String()
		if _, ok := seen[key]; !ok {
			// We have not seen this community yet
			result = append(result, com)
			seen[key] = true
		}
	}

	return result
}

type ExtCommunity []interface{}

func (com ExtCommunity) String() string {
	res := ""
	for _, v := range com {
		res += fmt.Sprintf(":%v", v)
	}
	return res[1:]
}

type ExtCommunities []ExtCommunity

func (communities ExtCommunities) Unique() ExtCommunities {
	seen := map[string]bool{}
	result := make(ExtCommunities, 0, len(communities))

	for _, com := range communities {
		key := com.String()
		if _, ok := seen[key]; !ok {
			// We have not seen this community yet
			result = append(result, com)
			seen[key] = true
		}
	}

	return result
}

type BgpInfo struct {
	Origin           string         `json:"origin"`
	AsPath           []int          `json:"as_path"`
	NextHop          string         `json:"next_hop"`
	Communities      Communities    `json:"communities"`
	LargeCommunities Communities    `json:"large_communities"`
	ExtCommunities   ExtCommunities `json:"ext_communities"`
	LocalPref        int            `json:"local_pref"`
	Med              int            `json:"med"`
}

func (bgp BgpInfo) HasCommunity(community Community) bool {
	if len(community) != 2 {
		return false // This can never match.
	}

	for _, com := range bgp.Communities {
		if len(com) != len(community) {
			continue // This can't match.
		}

		if com[0] == community[0] &&
			com[1] == community[1] {
			return true
		}
	}

	return false
}

func (bgp BgpInfo) HasExtCommunity(community ExtCommunity) bool {
	if len(community) != 3 {
		return false // This can never match.
	}

	for _, com := range bgp.ExtCommunities {
		if len(com) != len(community) {
			continue // This can't match.
		}

		if com[0] == community[0] &&
			com[1] == community[1] &&
			com[2] == community[2] {
			return true
		}
	}

	return false
}

func (bgp BgpInfo) HasLargeCommunity(community Community) bool {
	// TODO: This is an almost 1:1 match to the function above.
	if len(community) != 3 {
		return false // This can never match.
	}

	for _, com := range bgp.LargeCommunities {
		if len(com) != len(community) {
			continue // This can't match.
		}

		if com[0] == community[0] &&
			com[1] == community[1] &&
			com[2] == community[2] {
			return true
		}
	}

	return false
}
