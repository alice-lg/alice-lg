package main

var version = "unknown"

// Gather application status information
type AppStatus struct {
	Version string `json:"version"`
}

// Get application status, perform health checks
// on backends.
func NewAppStatus() (*AppStatus, error) {
	status := &AppStatus{
		Version: version,
	}
	return status, nil
}
