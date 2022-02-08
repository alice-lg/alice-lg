package birdwatcher

// Parsers and helpers

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
)

// Convert server time string to time
func parseServerTime(value interface{}, layout, timezone string) (time.Time, error) {
	svalue, ok := value.(string)
	if !ok {
		return time.Time{}, nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation(layout, svalue, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t.UTC(), nil
}

// Make api status from response:
// The api status is always included in a birdwatcher response
func parseApiStatus(bird ClientResponse, config Config) (api.ApiStatus, error) {
	birdApi, ok := bird["api"].(map[string]interface{})
	if !ok {
		// Define error status
		status := api.ApiStatus{
			Version:         "unknown / error",
			ResultFromCache: false,
			Ttl:             time.Now(),
		}

		// Try to retrieve the real error from server
		birdErr, ok := bird["error"].(string)
		if !ok {
			// Unknown error
			return status, fmt.Errorf("Invalid API response received from server")
		}

		return status, fmt.Errorf(birdErr)
	}

	// Parse TTL
	ttl, err := parseServerTime(
		bird["ttl"],
		config.ServerTime,
		config.Timezone,
	)
	if err != nil {
		return api.ApiStatus{}, err
	}

	// Parse Cache Status
	cacheStatus, _ := parseCacheStatus(birdApi, config)

	status := api.ApiStatus{
		Version:         birdApi["Version"].(string),
		ResultFromCache: birdApi["result_from_cache"].(bool),
		Ttl:             ttl,
		CacheStatus:     cacheStatus,
	}

	return status, nil
}

// Parse cache status from api response
func parseCacheStatus(cacheStatus map[string]interface{}, config Config) (api.CacheStatus, error) {
	cache, ok := cacheStatus["cache_status"].(map[string]interface{})
	if !ok {
		return api.CacheStatus{}, fmt.Errorf("Invalid Cache Status")
	}

	cachedAt, ok := cache["cached_at"].(map[string]interface{})
	if !ok {
		return api.CacheStatus{}, fmt.Errorf("Invalid Cache Status")
	}

	cachedAtTime, err := parseServerTime(cachedAt["date"], config.ServerTime, config.Timezone)
	if err != nil {
		return api.CacheStatus{}, err
	}

	status := api.CacheStatus{
		CachedAt: cachedAtTime,
		// We ommit OrigTTL for now...
	}

	return status, nil
}

// Parse birdwatcher status
func parseBirdwatcherStatus(bird ClientResponse, config Config) (api.Status, error) {
	birdStatus := bird["status"].(map[string]interface{})

	// Get special fields
	serverTime, _ := parseServerTime(
		birdStatus["current_server"],
		config.ServerTimeShort,
		config.Timezone,
	)

	lastReboot, _ := parseServerTime(
		birdStatus["last_reboot"],
		config.ServerTimeShort,
		config.Timezone,
	)

	if config.ShowLastReboot == false {
		lastReboot = time.Time{}
	}

	lastReconfig, _ := parseServerTime(
		birdStatus["last_reconfig"],
		config.ServerTimeExt,
		config.Timezone,
	)

	// Make status response
	status := api.Status{
		ServerTime:   serverTime,
		LastReboot:   lastReboot,
		LastReconfig: lastReconfig,
		Backend:      "bird",
		Version:      decoders.String(birdStatus["version"], "unknown"),
		Message:      decoders.String(birdStatus["message"], "unknown"),
		RouterId:     decoders.String(birdStatus["router_id"], "unknown"),
	}

	return status, nil
}

// Parse neighbour uptime
func parseRelativeServerTime(uptime interface{}, config Config) time.Duration {
	serverTime, _ := parseServerTime(uptime, config.ServerTimeShort, config.Timezone)
	return time.Since(serverTime)
}

// Parse neighbours response
func parseNeighbours(bird ClientResponse, config Config) (api.Neighbours, error) {
	rsId := config.ID
	neighbours := api.Neighbours{}
	protocols := bird["protocols"].(map[string]interface{})

	// Iterate over protocols map:
	for protocolId, proto := range protocols {
		protocol := proto.(map[string]interface{})
		routes := protocol["routes"].(map[string]interface{})

		uptime := parseRelativeServerTime(protocol["state_changed"], config)
		lastError := decoders.String(protocol["last_error"], "")

		routesReceived := float64(0)
		if routes != nil {
			if _, ok := routes["imported"]; ok {
				routesReceived = routesReceived + routes["imported"].(float64)
			}
			if _, ok := routes["filtered"]; ok {
				routesReceived = routesReceived + routes["filtered"].(float64)
			}
		}

		neighbour := &api.Neighbour{
			Id: protocolId,

			Address: decoders.String(protocol["neighbor_address"], "error"),
			Asn:     decoders.Int(protocol["neighbor_as"], 0),
			State: strings.ToLower(
				decoders.String(protocol["state"], "unknown")),
			Description: decoders.String(protocol["description"], "no description"),

			RoutesReceived:  decoders.Int(routesReceived, 0),
			RoutesAccepted:  decoders.Int(routes["imported"], 0),
			RoutesFiltered:  decoders.Int(routes["filtered"], 0),
			RoutesExported:  decoders.Int(routes["exported"], 0), //TODO protocol_exported?
			RoutesPreferred: decoders.Int(routes["preferred"], 0),

			Uptime:    uptime,
			LastError: lastError,

			RouteServerId: rsId,

			Details: protocol,
		}

		neighbours = append(neighbours, neighbour)
	}

	sort.Sort(neighbours)

	return neighbours, nil
}

// Parse neighbours response
func parseNeighboursShort(bird ClientResponse, config Config) (api.NeighboursStatus, error) {
	neighbours := api.NeighboursStatus{}
	protocols := bird["protocols"].(map[string]interface{})

	// Iterate over protocols map:
	for protocolId, proto := range protocols {
		protocol := proto.(map[string]interface{})

		uptime := parseRelativeServerTime(protocol["since"], config)

		neighbour := &api.NeighbourStatus{
			Id:    protocolId,
			State: decoders.String(protocol["state"], "unknown"),
			Since: uptime,
		}

		neighbours = append(neighbours, neighbour)
	}

	sort.Sort(neighbours)

	return neighbours, nil
}

// Parse route bgp info
func parseRouteBgpInfo(data interface{}) api.BgpInfo {
	bgpData, ok := data.(map[string]interface{})
	if !ok {
		// Info is missing
		return api.BgpInfo{}
	}

	asPath := decoders.IntList(bgpData["as_path"])
	communities := parseBgpCommunities(bgpData["communities"])
	largeCommunities := parseBgpCommunities(bgpData["large_communities"])
	extCommunities := parseExtBgpCommunities(bgpData["ext_communities"])

	localPref, _ := strconv.Atoi(decoders.String(bgpData["local_pref"], "0"))
	med, _ := strconv.Atoi(decoders.String(bgpData["med"], "0"))

	bgp := api.BgpInfo{
		Origin:           decoders.String(bgpData["origin"], "unknown"),
		AsPath:           asPath,
		NextHop:          decoders.String(bgpData["next_hop"], "unknown"),
		LocalPref:        localPref,
		Med:              med,
		Communities:      communities,
		ExtCommunities:   extCommunities,
		LargeCommunities: largeCommunities,
	}
	return bgp
}

// Extract bgp communities from response
func parseBgpCommunities(data interface{}) []api.Community {
	communities := []api.Community{}

	ldata, ok := data.([]interface{})
	if !ok { // We don't have any
		return []api.Community{}
	}

	for _, c := range ldata {
		cdata := c.([]interface{})
		community := api.Community{}
		for _, cinfo := range cdata {
			community = append(community, int(cinfo.(float64)))
		}
		communities = append(communities, community)
	}

	return communities
}

// Extract extended communtieis
func parseExtBgpCommunities(data interface{}) []api.ExtCommunity {
	communities := []api.ExtCommunity{}
	ldata, ok := data.([]interface{})
	if !ok { // We don't have any
		return communities
	}

	for _, c := range ldata {
		cdata := c.([]interface{})
		if len(cdata) != 3 {
			log.Println("Ignoring malformed ext community:", cdata)
			continue
		}
		communities = append(communities, api.ExtCommunity{
			cdata[0],
			cdata[1],
			cdata[2],
		})
	}

	return communities
}

// Parse partial routes response
func parseRoutesData(birdRoutes []interface{}, config Config) api.Routes {
	routes := api.Routes{}

	for _, data := range birdRoutes {
		rdata := data.(map[string]interface{})

		age := parseRelativeServerTime(rdata["age"], config)
		rtype := decoders.StringList(rdata["type"])
		bgpInfo := parseRouteBgpInfo(rdata["bgp"])

		route := &api.Route{
			Id:          decoders.String(rdata["network"], "unknown"),
			NeighbourId: decoders.String(rdata["from_protocol"], "unknown neighbour"),

			Network:   decoders.String(rdata["network"], "unknown net"),
			Interface: decoders.String(rdata["interface"], "unknown interface"),
			Gateway:   decoders.String(rdata["gateway"], "unknown gateway"),
			Metric:    decoders.Int(rdata["metric"], -1),
			Primary:   decoders.Bool(rdata["primary"], false),
			Age:       age,
			Type:      rtype,
			Bgp:       bgpInfo,

			Details: rdata,
		}

		routes = append(routes, route)
	}
	return routes
}

// Parse routes response
func parseRoutes(bird ClientResponse, config Config) (api.Routes, error) {
	birdRoutes, ok := bird["routes"].([]interface{})
	if !ok {
		return api.Routes{}, fmt.Errorf("Routes response missing")
	}

	routes := parseRoutesData(birdRoutes, config)

	// Sort routes
	sort.Sort(routes)
	return routes, nil
}

func parseRoutesDump(bird ClientResponse, config Config) (*api.RoutesResponse, error) {
	result := &api.RoutesResponse{}

	apiStatus, err := parseApiStatus(bird, config)
	if err != nil {
		return result, err
	}
	result.Api = apiStatus

	// Fetch imported routes
	importedRoutes, ok := bird["imported"].([]interface{})
	if !ok {
		return result, fmt.Errorf("Imported routes missing")
	}

	// Sort routes by network for faster querying
	imported := parseRoutesData(importedRoutes, config)
	sort.Sort(imported)
	result.Imported = imported

	// Fetch filtered routes
	filteredRoutes, ok := bird["filtered"].([]interface{})
	if !ok {
		return result, fmt.Errorf("Filtered routes missing")
	}
	filtered := parseRoutesData(filteredRoutes, config)
	sort.Sort(filtered)
	result.Filtered = filtered

	return result, nil
}
