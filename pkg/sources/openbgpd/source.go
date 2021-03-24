package openbgpd

import (
	"github.com/alice-lg/alice-lg/pkg/api"
)

// Source implements the OpenBGPD source for Alice.
// It is intendet to consume structured bgpctl output
// queried over HTTP using a `openbgpd-state-server`.
type Source struct {
	// API is the http host and api prefix. For
	// example http://rs1.mgmt.ixp.example.net:29111/api/v1
	API string
}

// ExpireCaches expires all cached data
func (s *Source) ExpireCaches() int {
	return 0 // Nothing to expire yet
}

// Status returns an API status response. In our case
// this is pretty much only that the service is available.
func (s *Source) Status() (*api.StatusResponse, error) {
	return nil, nil
}
