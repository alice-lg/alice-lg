package store

import (
	"context"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
)

// TestSourcesStoreReconcileAddsAndRemoves proves that reconcile
// computes the correct add and remove sets, preserves the status
// of surviving sources, and reuses their SourceConfig instance.
func TestSourcesStoreReconcileAddsAndRemoves(t *testing.T) {
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src1", Name: "src1"},
			{ID: "src2", Name: "src2"},
		},
	}
	s := NewSourcesStore(cfg, time.Minute, 1)

	// Warm src2 so we can prove its status survives a reconcile.
	if err := s.LockSource("src2"); err != nil {
		t.Fatal(err)
	}
	if err := s.RefreshSuccess("src2"); err != nil {
		t.Fatal(err)
	}

	newCfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src2", Name: "src2"},
			{ID: "src3", Name: "src3"},
		},
	}
	added, changed, removed := s.reconcile(newCfg)

	if len(added) != 1 || added[0] != "src3" {
		t.Error("expected added [src3], got", added)
	}
	if len(changed) != 0 {
		t.Error("expected no changed sources, got", changed)
	}
	if len(removed) != 1 || removed[0] != "src1" {
		t.Error("expected removed [src1], got", removed)
	}
	if s.Exists("src1") {
		t.Error("src1 should have been removed")
	}
	if !s.Exists("src2") || !s.Exists("src3") {
		t.Error("src2 and src3 should be present")
	}

	// The surviving source keeps its original config instance.
	if s.sources["src2"] != cfg.Sources[1] {
		t.Error("src2 config pointer should be preserved across reconcile")
	}

	// The surviving source keeps its warm (initialized) status.
	st2, err := s.GetStatus("src2")
	if err != nil {
		t.Fatal(err)
	}
	if !st2.Initialized {
		t.Error("src2 status should remain initialized across reconcile")
	}

	// The new source starts fresh with the store's refresh interval.
	st3, err := s.GetStatus("src3")
	if err != nil {
		t.Fatal(err)
	}
	if st3.Initialized {
		t.Error("src3 should not be initialized yet")
	}
	if st3.RefreshInterval != time.Minute {
		t.Error("src3 should inherit the refresh interval")
	}
}

// TestSourcesStoreReconcileNoChange proves a reconcile against an
// identical configuration is a no-op.
func TestSourcesStoreReconcileNoChange(t *testing.T) {
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src1", Name: "src1"},
			{ID: "src2", Name: "src2"},
		},
	}
	s := NewSourcesStore(cfg, time.Minute, 1)

	added, changed, removed := s.reconcile(cfg)
	if len(added) != 0 || len(changed) != 0 || len(removed) != 0 {
		t.Error(
			"expected no changes, got added", added,
			"changed", changed, "removed", removed)
	}
}

// TestRoutesStoreReconcilePurgesRemovedSource proves that routes
// for a removed source are purged from the backend, so stale routes
// no longer surface in a prefix lookup.
func TestRoutesStoreReconcilePurgesRemovedSource(t *testing.T) {
	ctx := context.Background()
	s := makeTestRoutesStore() // imports rs1 routes

	before, err := s.LookupPrefix(ctx, "193.200.", api.NewSearchFilters())
	if err != nil {
		t.Fatal(err)
	}
	if len(before) == 0 {
		t.Fatal("expected rs1 routes before reconcile")
	}

	// Drop rs1, keep rs2, add rs3.
	cfg := &config.Config{
		Server: config.ServerConfig{RoutesStoreRefreshInterval: 1},
		Sources: []*config.SourceConfig{
			{ID: "rs2", Name: "rs2"},
			{ID: "rs3", Name: "rs3"},
		},
	}
	s.Reconcile(ctx, cfg)

	if s.sources.Exists("rs1") {
		t.Error("rs1 should have been removed from the source set")
	}
	if !s.sources.Exists("rs3") {
		t.Error("rs3 should have been added to the source set")
	}

	after, err := s.LookupPrefix(ctx, "193.200.", api.NewSearchFilters())
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != 0 {
		t.Error("expected rs1 routes to be purged, got", len(after))
	}
}

// TestNeighborsStoreReconcilePurgesRemovedSource proves that
// neighbors for a removed source are purged from the backend.
func TestNeighborsStoreReconcilePurgesRemovedSource(t *testing.T) {
	ctx := context.Background()
	s := makeTestNeighborsStore() // rs1 and rs2 populated

	if _, err := s.backend.GetNeighborsAt(ctx, "rs1"); err != nil {
		t.Fatal("expected rs1 neighbors before reconcile:", err)
	}

	// Drop rs1, keep rs2.
	cfg := &config.Config{
		Server: config.ServerConfig{NeighborsStoreRefreshInterval: 1},
		Sources: []*config.SourceConfig{
			{ID: "rs2", Name: "rs2"},
		},
	}
	s.Reconcile(ctx, cfg)

	if s.sources.Exists("rs1") {
		t.Error("rs1 should have been removed from the source set")
	}
	// Assert against the backend directly: the store level
	// GetNeighborsAt would report the source as uninitialized
	// regardless of whether the backend was purged.
	if _, err := s.backend.GetNeighborsAt(ctx, "rs1"); err == nil {
		t.Error("expected rs1 neighbors to be purged from the backend")
	}
	if _, err := s.backend.GetNeighborsAt(ctx, "rs2"); err != nil {
		t.Error("rs2 neighbors should remain:", err)
	}
}
