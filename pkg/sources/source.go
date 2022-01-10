// Package sources provides the base interface for all
// route server data source implementations.
package sources

import (
	"errors"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Source Errors
var (
	// SourceNotFound indicates that a source could
	// not be resolved by an identifier.
	ErrSourceNotFound = errors.New("route server unknown")

	// ErrSourceBusy is returned when a refresh is
	// already in progress.
	ErrSourceBusy = errors.New("source is busy")
)

// Source is a generic datasource for alice.
// All route server adapters implement this interface.
type Source interface {
	ExpireCaches() int
	Status() (*api.StatusResponse, error)
	Neighbors() (*api.NeighborsResponse, error)
	NeighborsSummary() (*api.NeighborsResponse, error)
	NeighborsStatus() (*api.NeighborsStatusResponse, error)
	Routes(neighborID string) (*api.RoutesResponse, error)
	RoutesReceived(neighborID string) (*api.RoutesResponse, error)
	RoutesFiltered(neighborID string) (*api.RoutesResponse, error)
	RoutesNotExported(neighborID string) (*api.RoutesResponse, error)
	AllRoutes() (*api.RoutesResponse, error)
}
