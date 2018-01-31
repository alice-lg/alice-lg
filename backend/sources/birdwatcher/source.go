package birdwatcher

import (
	"github.com/ecix/alice-lg/backend/api"

	"log"
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

	// Optional: Filtered
	bird, _ = self.client.GetJson("/routes/filtered/" + neighbourId)
	filtered, err := parseRoutes(bird, self.config)
	if err != nil {
		log.Println("WARNING Could not retrieve filtered routes:", err)
		log.Println("Is the 'routes_filtered' module active in birdwatcher?")

		filtered = api.Routes{}
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
func (self *Birdwatcher) LookupPrefix(prefix string) (api.RoutesLookupResponse, error) {
	// Get RS info
	rs := api.Routeserver{
		Id:   self.config.Id,
		Name: self.config.Name,
	}

	// Query prefix on RS
	bird, err := self.client.GetJson("/routes/prefix?prefix=" + prefix)
	if err != nil {
		return api.RoutesLookupResponse{}, err
	}

	// Parse API status
	apiStatus, err := parseApiStatus(bird, self.config)
	if err != nil {
		return api.RoutesLookupResponse{}, err
	}

	// Parse routes
	routes, err := parseRoutes(bird, self.config)

	// Add corresponding neighbour and source rs to result
	results := []api.LookupRoute{}
	for _, src := range routes {
		// Okay. This is actually really hacky.
		// A less bruteforce approach would be highly appreciated
		route := api.LookupRoute{
			Id: src.Id,

			Routeserver: rs,

			NeighbourId: src.NeighbourId,

			Network:   src.Network,
			Interface: src.Interface,
			Gateway:   src.Gateway,
			Metric:    src.Metric,
			Bgp:       src.Bgp,
			Age:       src.Age,
			Type:      src.Type,

			Details: src.Details,
		}
		results = append(results, route)
	}

	// Make result
	response := api.RoutesLookupResponse{
		Api:    apiStatus,
		Routes: results,
	}
	return response, nil
}

func (self *Birdwatcher) AllRoutes() (api.RoutesResponse, error) {
	bird, err := self.client.GetJson("/routes/dump")
	if err != nil {
		return api.RoutesResponse{}, err
	}
	result, err := parseRoutesDump(bird, self.config)
	return result, err
}
