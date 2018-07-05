package birdwatcher

// Parsers and helpers

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/alice-lg/alice-lg/backend/api"
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

	ttl, err := parseServerTime(
		bird["ttl"],
		config.ServerTime,
		config.Timezone,
	)
	if err != nil {
		return api.ApiStatus{}, err
	}

	status := api.ApiStatus{
		Version:         birdApi["Version"].(string),
		ResultFromCache: birdApi["result_from_cache"].(bool),
		Ttl:             ttl,
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
		Version:      mustString(birdStatus["version"], "unknown"),
		Message:      mustString(birdStatus["message"], "unknown"),
		RouterId:     mustString(birdStatus["router_id"], "unknown"),
	}

	return status, nil
}

// Parse neighbour uptime
func parseRelativeServerTime(uptime interface{}, config Config) time.Duration {
	serverTime, _ := parseServerTime(uptime, config.ServerTimeShort, config.Timezone)
	return time.Since(serverTime)
}

// Parse neighbours response
func parseNeighbours(bird ClientResponse, config Config) ([]api.Neighbour, error) {
	neighbours := api.Neighbours{}
	protocols := bird["protocols"].(map[string]interface{})

	// Iterate over protocols map:
	for protocolId, proto := range protocols {
		protocol := proto.(map[string]interface{})
		routes := protocol["routes"].(map[string]interface{})

		uptime := parseRelativeServerTime(protocol["state_changed"], config)
		lastError := mustString(protocol["last_error"], "")

		neighbour := api.Neighbour{
			Id: protocolId,

			Address:     mustString(protocol["neighbor_address"], "error"),
			Asn:         mustInt(protocol["neighbor_as"], 0),
			State:       mustString(protocol["state"], "unknown"),
			Description: mustString(protocol["description"], "no description"),
			//TODO make these changes configurable
			RoutesReceived:     mustInt(routes["imported"], 0),
			RoutesExported:     mustInt(routes["exported"], 0), //TODO protocol_exported?
			RoutesFiltered:     mustInt(routes["filtered"], 0),
			RoutesPreferred:    mustInt(routes["preferred"], 0),
			RoutesAccepted:     mustInt(routes["pipe_imported"], mustInt(routes["imported"], 0)),
			RoutesPipeFiltered: mustInt(routes["pipe_filtered"], mustInt(routes["filtered"], 0)),

			Uptime:    uptime,
			LastError: lastError,

			Details: protocol,
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

	asPath := parseIntList(bgpData["as_path"])
	communities := parseBgpCommunities(bgpData["communities"])
	largeCommunities := parseBgpCommunities(bgpData["large_communities"])

	localPref, _ := strconv.Atoi(mustString(bgpData["local_pref"], "0"))
	med, _ := strconv.Atoi(mustString(bgpData["med"], "0"))

	bgp := api.BgpInfo{
		Origin:           mustString(bgpData["origin"], "unknown"),
		AsPath:           asPath,
		NextHop:          mustString(bgpData["next_hop"], "unknown"),
		LocalPref:        localPref,
		Med:              med,
		Communities:      communities,
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

// Assert string, provide default
func mustString(value interface{}, fallback string) string {
	sval, ok := value.(string)
	if !ok {
		return fallback
	}
	return sval
}

// Assert list of strings
func mustStringList(data interface{}) []string {
	list := []string{}
	ldata, ok := data.([]interface{})
	if !ok {
		return []string{}
	}
	for _, e := range ldata {
		s, ok := e.(string)
		if ok {
			list = append(list, s)
		}
	}
	return list
}

// Convert list of strings to int
func parseIntList(data interface{}) []int {
	list := []int{}
	sdata := mustStringList(data)
	for _, e := range sdata {
		val, _ := strconv.Atoi(e)
		list = append(list, val)
	}
	return list
}

func mustInt(value interface{}, fallback int) int {
	fval, ok := value.(float64)
	if !ok {
		return fallback
	}
	return int(fval)
}

// Parse partial routes response
func parseRoutesData(birdRoutes []interface{}, config Config) api.Routes {
	routes := api.Routes{}

	for _, data := range birdRoutes {
		rdata := data.(map[string]interface{})

		age := parseRelativeServerTime(rdata["age"], config)
		rtype := mustStringList(rdata["type"])
		bgpInfo := parseRouteBgpInfo(rdata["bgp"])

		route := api.Route{
			Id:          mustString(rdata["network"], "unknown"),
			NeighbourId: mustString(rdata["from_protocol"], "unknown neighbour"),

			Network:   mustString(rdata["network"], "unknown net"),
			Interface: mustString(rdata["interface"], "unknown interface"),
			Gateway:   mustString(rdata["gateway"], "unknown gateway"),
			Metric:    mustInt(rdata["metric"], -1),
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
func parseRoutes(bird ClientResponse, config Config) ([]api.Route, error) {
	birdRoutes, ok := bird["routes"].([]interface{})
	if !ok {
		return []api.Route{}, fmt.Errorf("Routes response missing")
	}

	routes := parseRoutesData(birdRoutes, config)

	// Sort routes
	sort.Sort(routes)
	return routes, nil
}

func parseRoutesDump(bird ClientResponse, config Config) (api.RoutesResponse, error) {
	result := api.RoutesResponse{}

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
