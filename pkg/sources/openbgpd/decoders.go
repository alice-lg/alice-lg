package openbgpd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
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
		Message:      "bgpd up and running",
		Version:      "1.0",
		Backend:      "openbgpd",
	}
	return s
}

// decodeNeighbor decodes a single neighbor in a
// bgpctl response.
func decodeNeighbor(n interface{}) (*api.Neighbour, error) {
	nb, ok := n.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("decode neighbor failed, interface is not a map")
	}
	neighbor := &api.Neighbour{
		Id:      decoders.String(nb, "remote_addr", "invalid_id"),
		Address: decoders.String(nb, "remote_addr", "invalid_address"),
	}
	return neighbor, nil
}

// decodenNeighbors retrievs neighbors data from
// the bgpctl response.
func decodeNeighbors(res map[string]interface{}) (api.Neighbours, error) {
	nbs := decoders.MapGet(res, "neighbors", nil)
	if nbs == nil {
		return nil, fmt.Errorf("missing neighbors in response body")
	}
	neighbors, ok := nbs.([]interface{})
	if !ok {
		return nil, fmt.Errorf("no a list of neighbors")
	}
	all := make(api.Neighbours, 0, len(neighbors))
	for _, nb := range neighbors {
		n, err := decodeNeighbor(nb)
		if err != nil {
			return nil, err
		}
		all = append(all, n)
	}
	return all, nil
}
