package openbgpd

import (
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
