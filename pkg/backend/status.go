package backend

// Version Alice (set during the build)
var Version = "unknown"

// Build is the current revision pointing at HEAD
var Build = "unknown"

// AppStatus contains application status information
type AppStatus struct {
	Version    string               `json:"version"`
	Routes     RoutesStoreStats     `json:"routes"`
	Neighbours NeighboursStoreStats `json:"neighbours"`
}

// NewAppStatus calculates the application status,
// and perform health checks on backends.
//
// TODO: Rename this.
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
		Version:    Version,
		Routes:     routesStatus,
		Neighbours: neighboursStatus,
	}
	return status, nil
}
