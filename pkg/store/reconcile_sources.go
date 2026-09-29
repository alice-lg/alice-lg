package store

import (
	"context"
	"errors"

	"github.com/alice-lg/alice-lg/pkg/config"
)

// ErrNoSources is returned when a reload would leave the looking
// glass with no sources; the running configuration is kept.
var ErrNoSources = errors.New("configuration contains no sources")

// ReconcileSources applies a reloaded configuration to the source
// registry and the stores, refusing one with no sources.
func ReconcileSources(
	ctx context.Context,
	cfg *config.Config,
	sources *SourcesRegistry,
	routesStore *RoutesStore,
	neighborsStore *NeighborsStore,
) error {
	if len(cfg.Sources) == 0 {
		return ErrNoSources
	}

	// The registry is reconciled even when the stores are
	// disabled, so the API always serves the configured sources.
	sources.Reconcile(cfg)
	if neighborsStore != nil {
		neighborsStore.Reconcile(ctx, cfg)
	}
	if routesStore != nil {
		routesStore.Reconcile(ctx, cfg)
	}
	return nil
}
