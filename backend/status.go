package main

var version = "unknown"

// Gather application status information
type AppStatus struct {
	Version    string               `json:"version"`
	Routes     RoutesStoreStats     `json:"routes"`
	Neighbours NeighboursStoreStats `json:"neighbours"`
}

// Get application status, perform health checks
// on backends.
func NewAppStatus() (*AppStatus, error) {
	routesStatus := RoutesStoreStats{}
	if AliceRoutesStore != nil {
		routesStatus = AliceRoutesStore.Stats()
	}

	neighboursStatus := NeighboursStoreStats{}
	if AliceRoutesStore != nil {
		neighboursStatus = AliceNeighboursStore.Stats()
	}

	status := &AppStatus{
		Version:    version,
		Routes:     routesStatus,
		Neighbours: neighboursStatus,
	}
	return status, nil
}
