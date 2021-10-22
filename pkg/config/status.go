package config

// Version Alice (set during the build)
var Version = "unknown"

// Build is the current revision pointing at HEAD
var Build = "unknown"

// AppStatus contains application status information
type AppStatus struct {
	Version   string              `json:"version"`
	Routes    RoutesStoreStats    `json:"routes"`
	Neighbors NeighborsStoreStats `json:"neighbors"`
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

	neighborsStatus := NeighborsStoreStats{}
	if AliceRoutesStore != nil {
		neighborsStatus = AliceNeighborsStore.Stats()
	}

	status := &AppStatus{
		Version:   Version,
		Routes:    routesStatus,
		Neighbors: neighborsStatus,
	}
	return status, nil
}
