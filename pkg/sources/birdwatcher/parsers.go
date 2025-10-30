package birdwatcher

// Parsers and helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
	"github.com/alice-lg/alice-lg/pkg/pools"
)

// Convert server time string to time
func parseServerTime(value any, layout, timezone string) (time.Time, error) {
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
func parseAPIStatus(bird ClientResponse, config Config) (*api.Meta, error) {
	birdAPI, ok := bird["api"].(map[string]any)
	if !ok {
		// Try to retrieve the real error from server
		birdErr, ok := bird["error"].(string)
		if !ok {
			// Unknown error
			return nil, fmt.Errorf("invalid API response received from server")
		}
		return nil, fmt.Errorf("bird error: %s", birdErr)
	}

	// Parse TTL
	ttl, err := parseServerTime(
		bird["ttl"],
		config.ServerTime,
		config.Timezone,
	)
	if err != nil {
		return nil, err
	}

	// Parse Cache Status
	cacheStatus, _ := parseCacheStatus(birdAPI, config)

	status := &api.Meta{
		Version:         birdAPI["Version"].(string),
		ResultFromCache: birdAPI["result_from_cache"].(bool),
		TTL:             ttl,
		CacheStatus:     cacheStatus,
	}

	return status, nil
}

// Parse cache status from api response
func parseCacheStatus(
	cacheStatus map[string]any,
	config Config,
) (api.CacheStatus, error) {
	cache, ok := cacheStatus["cache_status"].(map[string]any)
	if !ok {
		return api.CacheStatus{}, fmt.Errorf("invalid Cache Status")
	}

	cachedAt, ok := cache["cached_at"].(map[string]any)
	if !ok {
		return api.CacheStatus{}, fmt.Errorf("invalid Cache Status")
	}

	cachedAtTime, err := parseServerTime(
		cachedAt["date"], config.ServerTime, config.Timezone)
	if err != nil {
		return api.CacheStatus{}, err
	}

	status := api.CacheStatus{
		CachedAt: cachedAtTime,
		// We omit OrigTTL for now...
	}

	return status, nil
}

// Parse birdwatcher status
func parseBirdwatcherStatus(bird ClientResponse, config Config) (api.Status, error) {
	birdStatus := bird["status"].(map[string]any)

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

	if !config.ShowLastReboot {
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
		RouterID:     decoders.String(birdStatus["router_id"], "unknown"),
	}

	return status, nil
}

// Parse neighbor uptime
func parseRelativeServerTime(uptime any, config Config) time.Duration {
	serverTime, _ := parseServerTime(uptime, config.ServerTimeShort, config.Timezone)
	return time.Since(serverTime)
}

// parseRoutesChannels extracts per-channel route counts from protocol data
func parseRoutesChannels(protocol map[string]any) map[string]*api.RoutesChannel {
	routesChannels := make(map[string]*api.RoutesChannel)

	channels, ok := protocol["channels"].(map[string]any)
	if !ok {
		return routesChannels
	}

	for channelName, channelData := range channels {
		channel, ok := channelData.(map[string]any)
		if !ok {
			continue
		}

		routesCounts, ok := channel["routes_count"].(map[string]any)
		if !ok {
			continue
		}

		if len(routesCounts) == 0 {
			continue // No channel info
		}

		routesChannels[channelName] = &api.RoutesChannel{
			RoutesReceived:  decoders.Int(routesCounts["imported"], 0) + decoders.Int(routesCounts["filtered"], 0),
			RoutesAccepted:  decoders.Int(routesCounts["imported"], 0),
			RoutesFiltered:  decoders.Int(routesCounts["filtered"], 0),
			RoutesExported:  decoders.Int(routesCounts["exported"], 0),
			RoutesPreferred: decoders.Int(routesCounts["preferred"], 0),
		}
	}

	return routesChannels
}

// Parse neighbors response
func parseNeighbors(bird ClientResponse, config Config) (api.Neighbors, error) {
	rsID := config.ID
	neighbors := api.Neighbors{}
	protocols := bird["protocols"].(map[string]any)

	// Iterate over protocols map:
	for protocolID, proto := range protocols {
		protocol := proto.(map[string]any)
		routes := protocol["routes"].(map[string]any)

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

		neighbor := &api.Neighbor{
			ID: protocolID,

			Address: decoders.String(protocol["neighbor_address"], "error"),
			ASN:     decoders.Int(protocol["neighbor_as"], 0),
			State: strings.ToLower(
				decoders.String(protocol["state"], "unknown")),
			Description: decoders.String(protocol["description"], "no description"),

			RoutesReceived:  decoders.Int(routesReceived, 0),
			RoutesAccepted:  decoders.Int(routes["imported"], 0),
			RoutesFiltered:  decoders.Int(routes["filtered"], 0),
			RoutesExported:  decoders.Int(routes["exported"], 0), //TODO protocol_exported?
			RoutesPreferred: decoders.Int(routes["preferred"], 0),

			RoutesChannels: parseRoutesChannels(protocol),

			Uptime:    uptime,
			LastError: lastError,

			RouteServerID: rsID,

			Details: protocol,
		}

		neighbors = append(neighbors, neighbor)
	}

	sort.Sort(neighbors)

	return neighbors, nil
}

// Parse neighbors response
func parseNeighborsShort(bird ClientResponse, config Config) (api.NeighborsStatus, error) {
	neighbors := api.NeighborsStatus{}
	protocols := bird["protocols"].(map[string]any)

	// Iterate over protocols map:
	for protocolID, proto := range protocols {
		protocol := proto.(map[string]any)

		uptime := parseRelativeServerTime(protocol["since"], config)

		neighbor := &api.NeighborStatus{
			ID:    protocolID,
			State: decoders.String(protocol["state"], "unknown"),
			Since: uptime,
		}

		neighbors = append(neighbors, neighbor)
	}

	sort.Sort(neighbors)

	return neighbors, nil
}

// Parse route bgp info
func parseRouteBgpInfo(data any) *api.BGPInfo {
	gwpool := pools.Gateways4 // Let's see

	bgpData, ok := data.(map[string]any)
	if !ok {
		// Info is missing
		return &api.BGPInfo{}
	}

	asPath := decoders.IntList(bgpData["as_path"])
	communities := parseBgpCommunities(bgpData["communities"])
	largeCommunities := parseBgpCommunities(bgpData["large_communities"])
	extCommunities := parseExtBgpCommunities(bgpData["ext_communities"])

	localPref, _ := strconv.Atoi(decoders.String(bgpData["local_pref"], "0"))
	med, _ := strconv.Atoi(decoders.String(bgpData["med"], "0"))

	var otc *int
	otcv, err := strconv.Atoi(decoders.String(bgpData["otc"], "-"))
	if err == nil {
		otc = &otcv
	}

	// Using pools has a bit of a performance impact. While parsing
	// ~600000 routes without deduplication, this takes roughly 14 seconds.
	// With pools this is now 19 seconds.
	bgp := &api.BGPInfo{
		Origin: pools.Origins.Acquire(
			decoders.String(bgpData["origin"], "unknown")),
		AsPath: pools.ASPaths.Acquire(asPath),
		NextHop: gwpool.Acquire(
			decoders.String(bgpData["next_hop"], "unknown")),
		LocalPref:        localPref,
		Med:              med,
		OTC:              otc,
		Communities:      pools.CommunitiesSets.Acquire(communities),
		ExtCommunities:   pools.ExtCommunitiesSets.Acquire(extCommunities),
		LargeCommunities: pools.LargeCommunitiesSets.Acquire(largeCommunities),
	}
	return bgp
}

// Extract bgp communities from response
func parseBgpCommunities(data any) []api.Community {
	communities := []api.Community{}

	ldata, ok := data.([]any)
	if !ok { // We don't have any
		return []api.Community{}
	}

	for _, c := range ldata {
		cdata := c.([]any)
		community := api.Community{}
		for _, cinfo := range cdata {
			community = append(community, int(cinfo.(float64)))
		}
		communities = append(communities, community)
	}

	return communities
}

// Extract extended communities
func parseExtBgpCommunities(data any) []api.ExtCommunity {
	communities := []api.ExtCommunity{}
	ldata, ok := data.([]any)
	if !ok { // We don't have any
		return communities
	}

	for _, c := range ldata {
		cdata := c.([]any)
		if len(cdata) != 3 {
			log.Println("Ignoring malformed ext community:", cdata)
			continue
		}
		val1, _ := strconv.Atoi(cdata[1].(string))
		val2, _ := strconv.Atoi(cdata[2].(string))
		communities = append(communities, api.ExtCommunity{
			cdata[0],
			val1,
			val2,
		})
	}

	return communities
}

// Parse partial route
func parseRouteData(
	rdata map[string]any,
	config Config,
	keepDetails bool,
) *api.Route {
	gwpool := pools.Gateways4 // Let's see

	age := parseRelativeServerTime(rdata["age"], config)
	rtype := decoders.StringList(rdata["type"])
	bgpInfo := parseRouteBgpInfo(rdata["bgp"])

	// Precompute details as raw json message
	var details json.RawMessage = nil
	if keepDetails {
		detailsJSON, err := json.Marshal(rdata)
		if err != nil {
			log.Println("error while encoding details:", err)
		}
		details = json.RawMessage(detailsJSON)
	}

	gateway := decoders.String(rdata["gateway"], "unknown gateway")
	learntFrom := decoders.String(rdata["learnt_from"], "")
	if learntFrom == "" {
		learntFrom = gateway
	}

	network := decoders.String(rdata["network"], "unknown net")
	var addrFamily uint8 = 1 // Default to IPv4
	if strings.Contains(network, ":") {
		addrFamily = 2 // IPv6
	}

	route := &api.Route{
		// ID: decoders.String(rdata["network"], "unknown"),

		NeighborID: pools.Neighbors.Acquire(
			decoders.String(rdata["from_protocol"], "unknown neighbor")),
		Network: network,
		Interface: pools.Interfaces.Acquire(
			decoders.String(rdata["interface"], "unknown interface")),
		Metric:     decoders.Int(rdata["metric"], -1),
		Primary:    decoders.Bool(rdata["primary"], false),
		LearntFrom: gwpool.Acquire(learntFrom),
		Gateway:    gwpool.Acquire(gateway),
		Age:        age,
		Type:       pools.Types.Acquire(rtype),
		BGP:        bgpInfo,
		AddrFamily: addrFamily,

		Details: &details,
	}
	return route
}

// Parse partial routes response
func parseRoutesData(
	birdRoutes []any,
	config Config,
	keepDetails bool,
) api.Routes {
	routes := api.Routes{}

	for _, data := range birdRoutes {
		rdata := data.(map[string]any)
		route := parseRouteData(rdata, config, keepDetails)
		routes = append(routes, route)
	}
	return routes
}

// Parse routes response
func parseRoutes(
	bird ClientResponse,
	config Config,
	keepDetails bool,
) (api.Routes, error) {
	birdRoutes, ok := bird["routes"].([]any)
	if !ok {
		return api.Routes{}, fmt.Errorf("routes response missing")
	}

	routes := parseRoutesData(birdRoutes, config, keepDetails)

	// Sort routes
	sort.Sort(routes)
	return routes, nil
}

/*

Linter says parseRoutesDump is dead code.
So for now this is removed...


func parseRoutesDump(bird ClientResponse, config Config) (*api.RoutesResponse, error) {
	result := &api.RoutesResponse{}

	apiStatus, err := parseAPIStatus(bird, config)
	if err != nil {
		return result, err
	}
	result.Meta = apiStatus

	// Fetch imported routes
	importedRoutes, ok := bird["imported"].([]interface{})
	if !ok {
		return result, fmt.Errorf("imported routes missing")
	}

	// Sort routes by network for faster querying
	imported := parseRoutesData(importedRoutes, config)
	sort.Sort(imported)
	result.Imported = imported

	// Fetch filtered routes
	filteredRoutes, ok := bird["filtered"].([]interface{})
	if !ok {
		return result, fmt.Errorf("filtered routes missing")
	}
	filtered := parseRoutesData(filteredRoutes, config)
	sort.Sort(filtered)
	result.Filtered = filtered

	return result, nil
}
*/
