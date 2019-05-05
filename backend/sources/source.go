package sources

import (
	"github.com/alice-lg/alice-lg/backend/api"
)

type Source interface {
	ExpireCaches() int
	Status() (*api.StatusResponse, error)
	Neighbours() (*api.NeighboursResponse, error)
	Routes(neighbourId string) (*api.RoutesResponse, error)
	RoutesReceived(neighbourId string) (*api.RoutesResponse, error)
	RoutesFiltered(neighbourId string) (*api.RoutesResponse, error)
	RoutesNotExported(neighbourId string) (*api.RoutesResponse, error)
	AllRoutes() (*api.RoutesResponse, error)
}
