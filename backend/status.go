package main

import (
	"strings"

	"github.com/GeertJohan/go.rice"
)

// Gather application status information
type AppStatus struct {
	Version string `json:"version"`
}

// Get application status, perform health checks
// on backends.
func NewAppStatus() (*AppStatus, error) {
	status := &AppStatus{
		Version: statusGetVersion(),
	}
	return status, nil
}

// Get application version
func statusGetVersion() string {
	meta, err := rice.FindBox("../")
	if err != nil {
		return "unknown"
	}
	version, err := meta.String("VERSION")
	if err != nil {
		return "unknown"
	}
	return strings.Trim(version, "\n")
}
