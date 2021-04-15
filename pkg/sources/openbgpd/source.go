package openbgpd

import (
	"context"
	"net/http"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
)

const (
	// SourceVersion is currently fixed at 1.0
	SourceVersion = "1.0"
)

// Source implements the OpenBGPD source for Alice.
// It is intendet to consume structured bgpctl output
// queried over HTTP using a `openbgpd-state-server`.
type Source struct {
	// API is the http host and api prefix. For
	// example http://rs1.mgmt.ixp.example.net:29111/api
	API string
}

// ExpireCaches expires all cached data
func (s *Source) ExpireCaches() int {
	return 0 // Nothing to expire yet
}

// Status returns an API status response. In our case
// this is pretty much only that the service is available.
func (s *Source) Status() (*api.StatusResponse, error) {
	apiStatus := api.ApiStatus{
		Version:         SourceVersion,
		ResultFromCache: false,
		Ttl:             time.Now().UTC(),
	}

	// Make API request and read response
	req, err := StatusRequest(context.Background(), s)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	status := decodeAPIStatus(body)

	response := &api.StatusResponse{
		Api:    apiStatus,
		Status: status,
	}
	return response, nil
}

// Neighbours retrievs a full list of all neighbors
func (s *Source) Neighbours() (*api.NeighboursResponse, error) {
	// Retrieve neighbours
	apiStatus := api.ApiStatus{
		Version:         SourceVersion,
		ResultFromCache: false,
		Ttl:             time.Now().UTC(),
	}

	// Make API request and read response
	req, err := NeighborsRequest(context.Background(), s)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := decoders.ReadJSONResponse(res)
	if err != nil {
		return nil, err
	}

	nb, err := decodeNeighbors(body)
	if err != nil {
		return nil, err
	}

	response := &api.NeighboursResponse{
		Api:        apiStatus,
		Neighbours: nb,
	}
	return response, nil
}
