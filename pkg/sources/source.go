package sources

import (
	"github.com/alice-lg/alice-lg/pkg/api"
)

// Source is a generic datasource for alice.
// All route server adapters implement this interface.
type Source interface {
	ExpireCaches() int
	Status() (*api.StatusResponse, error)
	Neighbours() (*api.NeighboursResponse, error)
	NeighboursStatus() (*api.NeighboursStatusResponse, error)
	Routes(neighborID string) (*api.RoutesResponse, error)
	RoutesReceived(neighborID string) (*api.RoutesResponse, error)
	RoutesFiltered(neighborID string) (*api.RoutesResponse, error)
	RoutesNotExported(neighborID string) (*api.RoutesResponse, error)
	AllRoutes() (*api.RoutesResponse, error)
}
