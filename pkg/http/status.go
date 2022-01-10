package http

import (
	"context"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store"
)

// AppStatus contains application status information
type AppStatus struct {
	Version   string                   `json:"version"`
	Routes    *api.RoutesStoreStats    `json:"routes"`
	Neighbors *api.NeighborsStoreStats `json:"neighbors"`
}

// CollectAppStatus initializes the application
// status with stats gathered from the various
// application modules.
func CollectAppStatus(
	ctx context.Context,
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

	status := &AppStatus{
		Version:   config.Version,
		Routes:    routesStatus,
		Neighbors: neighborsStatus,
	}
	return status, nil

}
