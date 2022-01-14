package http

import (
	"context"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store"
	"github.com/alice-lg/alice-lg/pkg/store/backends/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

// AppStatus contains application status information
type AppStatus struct {
	Version   string                   `json:"version"`
	Routes    *api.RoutesStoreStats    `json:"routes"`
	Neighbors *api.NeighborsStoreStats `json:"neighbors"`
	Postgres  *postgres.Status         `json:"postgres"`
}

// CollectAppStatus initializes the application
// status with stats gathered from the various
// application modules.
func CollectAppStatus(
	ctx context.Context,
	pool *pgxpool.Pool,
	routesStore *store.RoutesStore,
	neighborsStore *store.NeighborsStore,
) (*AppStatus, error) {
	routesStatus := &api.RoutesStoreStats{}
	if routesStore != nil {
		routesStatus = routesStore.Stats()
	}

	neighborsStatus := &api.NeighborsStoreStats{}
	if neighborsStore != nil {
		neighborsStatus = neighborsStore.Stats(ctx)
	}

	var pgStatus *postgres.Status
	if pool != nil {
		pgStatus = postgres.NewManager(pool).Status(ctx)
	}

	status := &AppStatus{
		Version:   config.Version,
		Routes:    routesStatus,
		Neighbors: neighborsStatus,
		Postgres:  pgStatus,
	}

	return status, nil
}
