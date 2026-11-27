package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store/backends/memory"
)

func TestRoutesStoreUpdateSourceWithoutNeighbors(t *testing.T) {
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "rs1", Name: "rs1"},
		},
	}
	ns := NewNeighborsStore(cfg, memory.NewNeighborsBackend())
	rs := NewRoutesStore(ns, cfg, memory.NewRoutesBackend())

	// The neighbors store never starts, so rs1 has no neighbors. The
	// timeout only bounds a refresh that waits instead of returning.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := rs.updateSource(ctx, cfg.Sources[0])
	if !errors.Is(err, ErrSourceNotInitialized) {
		t.Errorf("expected ErrSourceNotInitialized, got %v", err)
	}
}
