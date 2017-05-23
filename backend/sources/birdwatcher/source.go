package birdwatcher

import (
	"fmt"
	"github.com/ecix/alice-lg/backend/api"
)

type Birdwatcher struct {
	config Config
	client *Client
}

func NewBirdwatcher(config Config) *Birdwatcher {
	client := NewClient(config.Api)

	birdwatcher := &Birdwatcher{
		config: config,
		client: client,
	}
	return birdwatcher
}

func (self *Birdwatcher) Status() (api.StatusResponse, error) {
	bird, err := self.client.GetJson("/status")
	if err != nil {
		return api.StatusResponse{}, err
	}

	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return api.StatusResponse{}, err
	}

	birdStatus, err := parseBirdwatcherStatus(bird, self.config)
	if err != nil {
		return api.StatusResponse{}, err
	}

	response := api.StatusResponse{
		Api:    apiStatus,
		Status: birdStatus,
	}

	return response, nil
}

// Get bird BGP protocols
func (self *Birdwatcher) Neighbours() (api.NeighboursResponse, error) {
	bird, err := self.client.GetJson("/protocols/bgp")
	if err != nil {
		return api.NeighboursResponse{}, err
	}

	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return api.NeighboursResponse{}, err
	}

	neighbours, err := parseNeighbours(bird, self.config)
	if err != nil {
		return api.NeighboursResponse{}, err
	}

	return api.NeighboursResponse{
		Api:        apiStatus,
		Neighbours: neighbours,
	}, nil
}

// Get filtered and exported routes
func (self *Birdwatcher) Routes(neighbourId string) (api.RoutesResponse, error) {
	// Exported
	bird, err := self.client.GetJson("/routes/protocol/" + neighbourId)
	if err != nil {
		return api.RoutesResponse{}, err
	}

	// Use api status from first request
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return api.RoutesResponse{}, err
	}

	imported, err := parseRoutes(bird, self.config)
	if err != nil {
		return api.RoutesResponse{}, err
	}

	// Filtered
	bird, err = self.client.GetJson("/routes/filtered/" + neighbourId)
	if err != nil {
		return api.RoutesResponse{}, err
	}

	filtered, err := parseRoutes(bird, self.config)
	if err != nil {
		return api.RoutesResponse{}, err
	}

	// Optional: NoExport
	bird, _ = self.client.GetJson("/routes/noexport/" + neighbourId)
	noexport, err := parseRoutes(bird, self.config)

	return api.RoutesResponse{
		Api:         apiStatus,
		Imported:    imported,
		Filtered:    filtered,
		NotExported: noexport,
	}, nil
}

// Make routes lookup
func (self *Birdwatcher) LookupPrefix(prefix string) (api.LookupResponse, error) {
	return api.LookupResponse{}, fmt.Errorf("not implemented")
}
