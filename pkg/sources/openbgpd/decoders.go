package openbgpd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
	"github.com/alice-lg/alice-lg/pkg/pools"
)

// Decode the api status response from the openbgpd
// state server.
func decodeAPIStatus(res map[string]interface{}) api.Status {
	now := time.Now().UTC()
	uptime := decoders.Duration(res["server_uptime"], 0)

	// This is an approximation and maybe wrong
	lastReboot := now.Add(-uptime)
	s := api.Status{
		ServerTime:   decoders.TimeUTC(res["server_time_utc"], time.Time{}),
		LastReboot:   lastReboot,
		LastReconfig: time.Time{},
		Message:      "openbgpd up and running",
		Version:      "",
		Backend:      "openbgpd",
	}
	return s
}

// decodeNeighbor decodes a single neighbor in a
// bgpctl response.
func decodeNeighbor(n interface{}) (*api.Neighbor, error) {
	nb, ok := n.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("decode neighbor failed, interface is not a map")
	}

	stats := decoders.MapGet(nb, "stats", map[string]interface{}{})
	prefixes := decoders.MapGet(stats, "prefixes", map[string]interface{}{})

	neighbor := &api.Neighbor{
		ID:             decoders.MapGetString(nb, "remote_addr", "invalid_id"),
		Address:        decoders.MapGetString(nb, "remote_addr", "invalid_address"),
		ASN:            decoders.IntFromString(decoders.MapGetString(nb, "remote_as", ""), -1),
		State:          decodeState(decoders.MapGetString(nb, "state", "unknown")),
		Description:    describeNeighbor(nb),
		RoutesReceived: int(decoders.MapGet(prefixes, "received", -1).(float64)),
		// TODO: RoutesFiltered
		RoutesExported: int(decoders.MapGet(prefixes, "sent", -1).(float64)),
		// TODO: RoutesPreferred
		// TODO: RoutesAccepted
		Uptime: decoders.DurationTimeframe(decoders.MapGet(nb, "last_updown", ""), 0),
	}
	return neighbor, nil
}

// describeNeighbor creates a neighbor description
func describeNeighbor(nb interface{}) string {
	desc := decoders.MapGetString(nb, "description", "")
	if desc != "" {
		return desc
	}

	addr := decoders.MapGetString(nb, "remote_addr", "invalid_address")
	asn := decoders.MapGetString(nb, "remote_as", "")
	return fmt.Sprintf("PEER AS%s %s", asn, addr)
}

// decodeNeighbors retrievs neighbors data from
// the bgpctl response.
func decodeNeighbors(res map[string]interface{}) (api.Neighbors, error) {
	nbs := decoders.MapGet(res, "neighbors", nil)
	if nbs == nil {
		return nil, fmt.Errorf("missing neighbors in response body")
	}
	neighbors, ok := nbs.([]interface{})
	if !ok {
		return nil, fmt.Errorf("no a list of neighbors")
	}
	all := make(api.Neighbors, 0, len(neighbors))
	for _, nb := range neighbors {
		n, err := decodeNeighbor(nb)
		if err != nil {
			return nil, err
		}
		all = append(all, n)
	}
	return all, nil
}

// decodeNeighborsStatus retrievs a neighbors summary
// and decodes the status.
func decodeNeighborsStatus(res map[string]interface{}) (api.NeighborsStatus, error) {
	nbs := decoders.MapGet(res, "neighbors", nil)
	if nbs == nil {
		return nil, fmt.Errorf("missing neighbors in response body")
	}
	neighbors, ok := nbs.([]interface{})
	if !ok {
		return nil, fmt.Errorf("no a list of interfaces")
	}

	all := make(api.NeighborsStatus, 0, len(neighbors))
	for _, nb := range neighbors {
		status := decodeNeighborStatus(nb)
		all = append(all, status)
	}

	return all, nil
}

// decodeNeighborStatus decodes a single status from a
// list of neighbor summaries.
func decodeNeighborStatus(nb interface{}) *api.NeighborStatus {
	id := decoders.MapGetString(nb, "bgpid", "undefined")
	state := decodeState(decoders.MapGetString(nb, "state", "Down"))
	uptime := decoders.DurationTimeframe(decoders.MapGet(nb, "last_updown", ""), 0)
	return &api.NeighborStatus{
		ID:    id,
		State: state,
		Since: uptime,
	}
}

// decodeRoutes decodes a response with a rib query.
// The toplevel element is expected to be "rib".
func decodeRoutes(res interface{}) (api.Routes, error) {
	r := decoders.MapGet(res, "rib", nil)
	if r == nil {
		// The response was a valid json but empty. So no
		// routes are present.
		return api.Routes{}, nil
	}
	rib, ok := r.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a list of interfaces")
	}
	routes := make(api.Routes, 0, len(rib))
	for _, details := range rib {
		route, err := decodeRoute(details.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// decodeRoute decodes a single route received from the source
func decodeRoute(details map[string]interface{}) (*api.Route, error) {
	prefix := decoders.MapGetString(details, "prefix", "")
	origin := decoders.MapGetString(details, "origin", "")
	neighbor := decoders.MapGet(details, "neighbor", nil)
	neighborID := "unknown"
	if neighbor != nil {
		neighborID = decoders.MapGetString(neighbor, "remote_addr", neighborID)
	}
	trueNextHop := decoders.MapGetString(details, "true_nexthop", "")
	lastUpdate := decoders.DurationTimeframe(
		decoders.MapGet(details, "last_update", nil), 0)

	asPath := decodeASPath(decoders.MapGetString(details, "aspath", ""))
	localPref := int(decoders.MapGet(details, "localpref", 0).(float64))

	// Decode BGP communities
	communities := decodeCommunities(
		decoders.MapGet(details, "communities", nil))
	largeCommunities := decodeCommunities(
		decoders.MapGet(details, "large_communities", nil))
	extendedCommunities := decodeExtendedCommunities(
		decoders.MapGet(details, "extended_communities", nil))

	// Is preferred route
	isPrimary := decoders.MapGetBool(details, "best", false)

	// Make bgp info
	bgpInfo := &api.BGPInfo{
		Origin:           pools.Origins.Acquire(origin),
		AsPath:           pools.ASPaths.Acquire(asPath),
		NextHop:          pools.Gateways4.Acquire(trueNextHop),
		Communities:      pools.Communities.Acquire(communities),
		ExtCommunities:   extendedCommunities,
		LargeCommunities: pools.LargeCommunities.Acquire(largeCommunities),
		LocalPref:        localPref,
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		log.Println("error while encoding details:", err)
	}
	rawDetails := json.RawMessage(detailsJSON)

	r := &api.Route{
		ID:         prefix,
		NeighborID: pools.Neighbors.Acquire(neighborID),
		Network:    pools.Networks4.Acquire(prefix),
		Gateway:    pools.Gateways4.Acquire(trueNextHop),
		BGP:        bgpInfo,
		Age:        lastUpdate,
		Type:       pools.Types.Acquire([]string{origin}),
		Primary:    isPrimary,
		Details:    &rawDetails,
	}
	return r, nil
}

// decodeState will decode the state into a canonical form
// used by the looking glass.
func decodeState(s string) string {
	s = strings.ToLower(s) // todo elaborate
	return s
}

// decodeASPath decodes a space separated list of
// string encoded ASNs into a list of integers.
func decodeASPath(path string) []int {
	tokens := strings.Split(path, " ")
	return decoders.IntListFromStrings(tokens)
}

// decodeCommunities decodes communities into a list of
// list of ints.
func decodeCommunities(c interface{}) api.Communities {
	details := decoders.StringList(c)
	comms := make(api.Communities, 0, len(details))
	for _, com := range details {
		tokens := strings.Split(com, ":")
		comms = append(comms, decoders.IntListFromStrings(tokens))
	}
	return comms
}

// decodeExtendedCommunities decodes extended communties
// into a list of (str, int, int).
func decodeExtendedCommunities(c interface{}) api.ExtCommunities {
	details := decoders.StringList(c)
	comms := make(api.ExtCommunities, 0, len(details))
	for _, com := range details {
		tokens := strings.SplitN(com, " ", 2)
		if len(tokens) != 2 {
			continue
		}
		nums := decoders.IntListFromStrings(
			strings.SplitN(tokens[1], ":", 2))
		if len(nums) != 2 {
			continue
		}
		comms = append(comms, []interface{}{tokens[0], nums[0], nums[1]})
	}
	return comms
}
