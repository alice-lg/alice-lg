package store

import (
	"context"
	"errors"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
)

// TestReconcileSourcesRefusesEmptyConfiguration proves a reload
// with no sources, such as a truncated file, keeps every live
// source and its routes.
func TestReconcileSourcesRefusesEmptyConfiguration(t *testing.T) {
	ctx := context.Background()
	s := makeTestRoutesStore()
	registry := NewSourcesRegistry(&config.Config{
		Sources: []*config.SourceConfig{s.sources.Get("rs1"), s.sources.Get("rs2")},
	})

	err := ReconcileSources(ctx, &config.Config{}, registry, s, s.neighbors)
	if !errors.Is(err, ErrNoSources) {
		t.Fatalf("expected ErrNoSources, got %v", err)
	}
	if len(registry.Sources()) != 2 || !s.sources.Exists("rs1") {
		t.Error("an empty reload changed the live sources")
	}
	found, err := s.LookupPrefix(ctx, "193.200.", api.NewSearchFilters())
	if err != nil {
		t.Fatal(err)
	}
	if len(found) == 0 {
		t.Error("an empty reload purged the warmed routes")
	}

	err = ReconcileSources(ctx, &config.Config{
		Sources: []*config.SourceConfig{s.sources.Get("rs2")},
	}, registry, s, s.neighbors)
	if err != nil {
		t.Fatal(err)
	}
	if registry.SourceByID("rs1") != nil || s.sources.Exists("rs1") {
		t.Error("a non-empty reload should remove rs1")
	}
}
