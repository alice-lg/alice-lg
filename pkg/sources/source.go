// Package sources provides the base interface for all
// route server data source implementations.
package sources

import (
	"context"
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

	Status(context.Context) (*api.StatusResponse, error)
	Neighbors(context.Context) (*api.NeighborsResponse, error)
	NeighborsSummary(context.Context) (*api.NeighborsResponse, error)
	NeighborsStatus(context.Context) (*api.NeighborsStatusResponse, error)
	Routes(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	RoutesReceived(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	RoutesFiltered(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	RoutesNotExported(ctx context.Context, neighborID string) (*api.RoutesResponse, error)
	AllRoutes(context.Context) (*api.RoutesResponse, error)
}
