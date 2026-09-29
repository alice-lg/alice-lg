package http

import (
	"context"
	"sort"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store"
)

// Collect the route server IDs from a list response, sorted so
// the assertion does not depend on the response ordering.
func routeServerIDs(res response) []string {
	rs, ok := res.(api.RouteServersResponse)
	if !ok {
		return nil
	}
	ids := make([]string, 0, len(rs.RouteServers))
	for _, s := range rs.RouteServers {
		ids = append(ids, s.ID)
	}
	sort.Strings(ids)
	return ids
}

// TestApiRouteServersListFollowsReload proves the route server
// list tracks a configuration reload, so a source added on
// SIGHUP is served and a removed one disappears, without a
// restart.
func TestApiRouteServersListFollowsReload(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "rs1", Name: "rs1"},
			{ID: "rs2", Name: "rs2"},
		},
	}
	s := &Server{cfg: cfg, sources: store.NewSourcesRegistry(cfg)}

	res, err := s.apiRouteServersList(ctx, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	ids := routeServerIDs(res)
	if len(ids) != 2 || ids[0] != "rs1" || ids[1] != "rs2" {
		t.Fatal("expected rs1 and rs2, got", ids)
	}

	// Drop rs1, keep rs2, add rs3.
	s.sources.Reconcile(&config.Config{
		Sources: []*config.SourceConfig{
			{ID: "rs2", Name: "rs2"},
			{ID: "rs3", Name: "rs3"},
		},
	})

	res, err = s.apiRouteServersList(ctx, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	ids = routeServerIDs(res)
	if len(ids) != 2 || ids[0] != "rs2" || ids[1] != "rs3" {
		t.Error("expected rs2 and rs3 after the reload, got", ids)
	}

	// The configuration loaded at startup is never rewritten,
	// which is why the API cannot be served from it.
	if len(cfg.Sources) != 2 || cfg.Sources[0].ID != "rs1" {
		t.Error("startup config should be untouched, got", cfg.Sources)
	}
}

// TestApiSourceLookupFollowsReload proves a source removed by a
// reload no longer resolves, so its endpoints stop serving it.
func TestApiSourceLookupFollowsReload(t *testing.T) {
	src := &config.SourceConfig{
		ID:      "rs1",
		Name:    "rs1",
		Type:    config.SourceTypeBird,
		Backend: config.SourceBackendBirdwatcher,
	}
	// A birdwatcher only builds an instance for a known table
	// layout.
	src.Birdwatcher.Type = "single_table"
	src.Birdwatcher.API = "http://192.0.2.1:29184"

	cfg := &config.Config{
		Sources: []*config.SourceConfig{src},
	}
	s := &Server{cfg: cfg, sources: store.NewSourcesRegistry(cfg)}

	if src := s.sources.SourceByID("rs1"); src == nil {
		t.Fatal("rs1 should resolve before the reload")
	}
	if inst := s.sources.SourceInstanceByID("rs1"); inst == nil {
		t.Fatal("rs1 should have an instance before the reload")
	}

	s.sources.Reconcile(&config.Config{})

	if src := s.sources.SourceByID("rs1"); src != nil {
		t.Error("rs1 should not resolve after it was removed")
	}
	if inst := s.sources.SourceInstanceByID("rs1"); inst != nil {
		t.Error("rs1 should have no instance after it was removed")
	}
}
