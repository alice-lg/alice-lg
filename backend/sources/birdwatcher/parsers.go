package birdwatcher

// Parsers and helpers

import (
	"time"

	"github.com/ecix/alice-lg/backend/api"
)

const SERVER_TIME = time.RFC3339Nano
const SERVER_TIME_SHORT = "2006-01-02 15:04:05"
const SERVER_TIME_EXT = "Mon, 2 Jan 2006 15:04:05 +0000"

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
	return t, err
}

// Make api status from response:
// The api status is always included in a birdwatcher response
func parseApiStatus(bird ClientResponse, config Config) (api.ApiStatus, error) {
	birdApi := bird["api"].(map[string]interface{})

	ttl, err := parseServerTime(
		bird["ttl"],
		SERVER_TIME,
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

// Parse neighbour uptime
func parseRelativeServerTime(uptime interface{}, config Config) time.Duration {
	serverTime, _ := parseServerTime(uptime, SERVER_TIME_SHORT, config.Timezone)
	return time.Since(serverTime)
}

// Parse neighbours response
func parseNeighbours(bird ClientResponse, config Config) ([]api.Neighbour, error) {
	neighbours := []api.Neighbour{}
	protocols := bird["protocols"].(map[string]interface{})

	// Iterate over protocols map:
	for protocolId, proto := range protocols {
		protocol := proto.(map[string]interface{})
		routes := protocol["routes"].(map[string]interface{})

		uptime := parseRelativeServerTime(protocol["state_changed"], config)

		neighbour := api.Neighbour{
			Id: protocolId,

			Address:     protocol["neighbor_address"].(string),
			Asn:         int(protocol["neighbor_as"].(float64)),
			State:       protocol["state"].(string),
			Description: protocol["description"].(string),

			RoutesReceived:  int(routes["imported"].(float64)),
			RoutesExported:  int(routes["exported"].(float64)),
			RoutesFiltered:  int(routes["filtered"].(float64)),
			RoutesPreferred: int(routes["preferred"].(float64)),

			Uptime: uptime,

			Details: protocol,
		}

		neighbours = append(neighbours, neighbour)
	}

	return neighbours, nil
}

// Parse route bgp info
func parseRouteBgpInfo(data interface{}) api.BgpInfo {
	bgpData := data.(map[string]interface{})

	_ = bgpData

	bgp := api.BgpInfo{}
	return bgp
}

// Get route type information
func parseRouteType(data interface{}) []string {
	rtype := []string{}
	tdata := data.([]interface{})
	for _, t := range tdata {
		rtype = append(rtype, t.(string))
	}
	return rtype
}

// Parse routes response
func parseRoutes(bird ClientResponse, config Config) ([]api.Route, error) {
	routes := []api.Route{}
	birdRoutes := bird["routes"].([]interface{})

	for _, data := range birdRoutes {
		rdata := data.(map[string]interface{})

		age := parseRelativeServerTime(rdata["age"], config)
		rtype := parseRouteType(rdata["type"])
		bgpInfo := parseRouteBgpInfo(rdata["bgp"])

		route := api.Route{
			Id:          rdata["network"].(string),
			NeighbourId: rdata["from_protocol"].(string),

			Network:   rdata["network"].(string),
			Interface: rdata["interface"].(string),
			Gateway:   rdata["gateway"].(string),
			Metric:    int(rdata["metric"].(float64)),
			Age:       age,
			Type:      rtype,
			Bgp:       bgpInfo,

			Details: rdata,
		}

		routes = append(routes, route)
	}

	return routes, nil
}
