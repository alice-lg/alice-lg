package store

import (
	"context"
	"testing"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store/backends/memory"
	"github.com/alice-lg/alice-lg/pkg/store/testdata"
)

// A blockingRoutesBackend pauses inside SetRoutes so a test can
// land a reconcile between the source check and the write.
type blockingRoutesBackend struct {
	RoutesStoreBackend
	entered chan struct{}
	release chan struct{}
}

func (b *blockingRoutesBackend) SetRoutes(
	ctx context.Context,
	sourceID string,
	routes api.LookupRoutes,
) error {
	close(b.entered)
	<-b.release
	return b.RoutesStoreBackend.SetRoutes(ctx, sourceID, routes)
}

// A blockingNeighborsBackend pauses inside SetNeighbors for the
// same reason.
type blockingNeighborsBackend struct {
	NeighborsStoreBackend
	entered chan struct{}
	release chan struct{}
}

func (b *blockingNeighborsBackend) SetNeighbors(
	ctx context.Context,
	sourceID string,
	neighbors api.Neighbors,
) error {
	close(b.entered)
	<-b.release
	return b.NeighborsStoreBackend.SetNeighbors(ctx, sourceID, neighbors)
}

// TestSourcesStoreAccessorsAfterRemovalDoNotPanic proves the
// accessors a refresh uses after its jitter sleep tolerate a
// source that a reload removed while it slept.
func TestSourcesStoreAccessorsAfterRemovalDoNotPanic(t *testing.T) {
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src1", Name: "src1"},
		},
	}
	s := NewSourcesStore(cfg, time.Minute, 1)

	// A refresh locks the source, then sleeps for jitter.
	if err := s.LockSource("src1"); err != nil {
		t.Fatal(err)
	}

	// The reload lands during that sleep.
	s.reconcile(&config.Config{})

	if src := s.Get("src1"); src != nil {
		t.Error("expected nil source config, got", src)
	}
	if name := s.GetName("src1"); name != "" {
		t.Error("expected empty name, got", name)
	}
	if inst := s.GetInstance("src1"); inst != nil {
		t.Error("expected nil instance, got", inst)
	}
}

// TestRoutesStoreRefreshSurvivesSourceRemoval proves the routes
// refresh exits cleanly when its source disappears mid refresh.
func TestRoutesStoreRefreshSurvivesSourceRemoval(t *testing.T) {
	ctx := context.Background()
	s := makeTestRoutesStore()

	if err := s.sources.LockSource("rs1"); err != nil {
		t.Fatal(err)
	}
	s.Reconcile(ctx, &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "rs2", Name: "rs2"},
		},
	})

	// This is the sequence safeUpdateSource runs after its
	// jitter sleep. None of it may panic.
	src := s.sources.Get("rs1")
	name := s.sources.GetName("rs1")
	if src != nil || name != "" {
		t.Error("expected rs1 to be gone, got", src, name)
	}
	s.sources.RefreshError("rs1", "refresh failed")
}

// TestAwaitNeighborStoreReturnsWhenSourceRemoved proves the
// routes refresh stops waiting on a neighbor store that will
// never initialize a source which no longer exists.
func TestAwaitNeighborStoreReturnsWhenSourceRemoved(t *testing.T) {
	ctx := context.Background()
	s := makeTestRoutesStore()

	// rs3 is registered but its neighbors never initialize.
	s.Reconcile(ctx, &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "rs3", Name: "rs3"},
		},
	})

	done := make(chan error, 1)
	go func() {
		done <- s.awaitNeighborStore(ctx, "rs3")
	}()

	// Remove it again; the wait must end.
	s.Reconcile(ctx, &config.Config{})

	select {
	case err := <-done:
		if err == nil {
			t.Error("expected an error for the removed source")
		}
	case <-time.After(5 * time.Second):
		t.Error("awaitNeighborStore did not return for a removed source")
	}
}

// TestRoutesStoreReconcileBeatsInFlightWrite proves a refresh
// that started before a reload cannot write routes back for a
// source the reload purged.
func TestRoutesStoreReconcileBeatsInFlightWrite(t *testing.T) {
	ctx := context.Background()
	s := makeTestRoutesStore()

	blocking := &blockingRoutesBackend{
		RoutesStoreBackend: s.backend,
		entered:            make(chan struct{}),
		release:            make(chan struct{}),
	}
	s.backend = blocking

	routes := testdata.LoadTestLookupRoutes("rs1", "rs1")
	src := s.sources.Get("rs1")
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- s.setRoutes(ctx, src, routes)
	}()
	<-blocking.entered // the write is past the source check

	reconcileDone := make(chan struct{})
	go func() {
		s.Reconcile(ctx, &config.Config{
			Sources: []*config.SourceConfig{
				{ID: "rs2", Name: "rs2"},
			},
		})
		close(reconcileDone)
	}()

	// Give the reload a chance to get ahead of the write.
	time.Sleep(100 * time.Millisecond)
	close(blocking.release)

	<-writeDone
	<-reconcileDone

	found, err := s.LookupPrefix(ctx, "193.200.", api.NewSearchFilters())
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Error(
			"a refresh in flight restored routes for a purged source, got",
			len(found))
	}
}

// TestNeighborsStoreReconcileBeatsInFlightWrite proves the same
// for the neighbors store.
func TestNeighborsStoreReconcileBeatsInFlightWrite(t *testing.T) {
	ctx := context.Background()
	s := makeTestNeighborsStore()

	blocking := &blockingNeighborsBackend{
		NeighborsStoreBackend: s.backend,
		entered:               make(chan struct{}),
		release:               make(chan struct{}),
	}
	s.backend = blocking

	neighbors := api.Neighbors{
		&api.Neighbor{ID: "ID2233_AS2342", ASN: 2342},
	}
	inst := s.sources.GetInstance("rs1")
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- s.setNeighbors(ctx, inst, "rs1", neighbors)
	}()
	<-blocking.entered

	reconcileDone := make(chan struct{})
	go func() {
		s.Reconcile(ctx, &config.Config{
			Sources: []*config.SourceConfig{
				{ID: "rs2", Name: "rs2"},
			},
		})
		close(reconcileDone)
	}()

	time.Sleep(100 * time.Millisecond)
	close(blocking.release)

	<-writeDone
	<-reconcileDone

	if _, err := s.backend.GetNeighborsAt(ctx, "rs1"); err == nil {
		t.Error("a refresh in flight restored neighbors for a purged source")
	}
}

// TestSourcesStoreReconcileDetectsChangedSource proves an edited
// source definition is applied on reload rather than silently
// kept at its old settings.
func TestSourcesStoreReconcileDetectsChangedSource(t *testing.T) {
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src1", Name: "src1"},
		},
	}
	s := NewSourcesStore(cfg, time.Minute, 1)
	if err := s.LockSource("src1"); err != nil {
		t.Fatal(err)
	}
	if err := s.RefreshSuccess("src1"); err != nil {
		t.Fatal(err)
	}

	// Same ID, different API endpoint.
	newSrc := &config.SourceConfig{ID: "src1", Name: "src1 renamed"}
	newSrc.Birdwatcher.API = "http://192.0.2.1:29184"
	added, changed, removed := s.reconcile(&config.Config{
		Sources: []*config.SourceConfig{newSrc},
	})

	if len(added) != 0 || len(removed) != 0 {
		t.Error("expected no adds or removes, got", added, removed)
	}
	if len(changed) != 1 || changed[0] != "src1" {
		t.Error("expected changed [src1], got", changed)
	}
	if s.Get("src1") != newSrc {
		t.Error("changed source should adopt the new configuration")
	}
	if s.GetName("src1") != "src1 renamed" {
		t.Error("changed source should serve the new name")
	}

	// A changed source has to re-warm against its new settings.
	st, err := s.GetStatus("src1")
	if err != nil {
		t.Fatal(err)
	}
	if st.Initialized {
		t.Error("changed source status should be reset")
	}
}

// TestSourcesStoreReconcileIgnoresEquivalentConfig proves a
// reload of an unchanged file keeps the memoized instance and
// reports nothing.
func TestSourcesStoreReconcileIgnoresEquivalentConfig(t *testing.T) {
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src1", Name: "src1", Blackholes: []string{"192.0.2.1"}},
		},
	}
	s := NewSourcesStore(cfg, time.Minute, 1)
	held := s.Get("src1")

	// A freshly parsed but identical configuration.
	added, changed, removed := s.reconcile(&config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src1", Name: "src1", Blackholes: []string{"192.0.2.1"}},
		},
	})

	if len(added) != 0 || len(changed) != 0 || len(removed) != 0 {
		t.Error("expected no changes, got", added, changed, removed)
	}
	if s.Get("src1") != held {
		t.Error("an unchanged source should keep its config instance")
	}
}

// TestRoutesStoreRejectsWriteFromReplacedSource proves a refresh
// that started against the previous definition of an edited
// source cannot write its result once the reload has landed.
func TestRoutesStoreRejectsWriteFromReplacedSource(t *testing.T) {
	ctx := context.Background()
	s := makeTestRoutesStore()

	stale := s.sources.Get("rs1")
	routes := testdata.LoadTestLookupRoutes("rs1", "rs1")

	// The operator points rs1 somewhere else and reloads.
	edited := &config.SourceConfig{ID: "rs1", Name: "rs1"}
	edited.Birdwatcher.API = "http://192.0.2.9:29184"
	s.Reconcile(ctx, &config.Config{
		Sources: []*config.SourceConfig{edited},
	})

	if err := s.setRoutes(ctx, stale, routes); err == nil {
		t.Error("a write from the replaced definition was accepted")
	}

	found, err := s.LookupPrefix(ctx, "193.200.", api.NewSearchFilters())
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Error("expected no routes for the re-warming source, got", len(found))
	}
}

// TestSourcesStoreReconcileTreatsMovedSourceAsUnchanged proves that
// inserting a source ahead of others does not purge and refetch them.
func TestSourcesStoreReconcileTreatsMovedSourceAsUnchanged(t *testing.T) {
	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src1", Name: "src1", Order: 0},
			{ID: "src2", Name: "src2", Order: 1},
		},
	}
	s := NewSourcesStore(cfg, time.Minute, 1)

	newCfg := &config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src0", Name: "src0", Order: 0},
			{ID: "src1", Name: "src1", Order: 1},
			{ID: "src2", Name: "src2", Order: 2},
		},
	}
	added, changed, removed := s.reconcile(newCfg)

	if len(added) != 1 || added[0] != "src0" {
		t.Error("expected added [src0], got", added)
	}
	if len(changed) != 0 {
		t.Error("expected moved sources to be unchanged, got", changed)
	}
	if len(removed) != 0 {
		t.Error("expected no removed sources, got", removed)
	}
}

// TestSourcesRegistryReconcileMovesSourceAndKeepsInstance proves the
// registry takes a moved source's new position and keeps its client.
func TestSourcesRegistryReconcileMovesSourceAndKeepsInstance(t *testing.T) {
	src := &config.SourceConfig{
		ID:      "src1",
		Name:    "src1",
		Order:   0,
		Backend: config.SourceBackendBirdwatcher,
	}
	r := NewSourcesRegistry(&config.Config{Sources: []*config.SourceConfig{src}})
	instance := r.SourceInstanceByID("src1")

	r.Reconcile(&config.Config{
		Sources: []*config.SourceConfig{
			{ID: "src0", Name: "src0", Order: 0, Backend: config.SourceBackendBirdwatcher},
			{ID: "src1", Name: "src1", Order: 1, Backend: config.SourceBackendBirdwatcher},
		},
	})

	if got := r.SourceByID("src1").Order; got != 1 {
		t.Errorf("expected src1 at order 1, got %d", got)
	}
	if r.SourceInstanceByID("src1") != instance {
		t.Error("the moved source should keep its instance")
	}
}

// TestRoutesStoreRefreshSuccessIgnoresReplacedSource proves a refresh
// that finishes after a reload replaced its source does not mark the
// empty replacement as initialized.
func TestRoutesStoreRefreshSuccessIgnoresReplacedSource(t *testing.T) {
	ctx := context.Background()
	s := makeTestRoutesStore()
	stale := s.sources.Get("rs1")

	edited := &config.SourceConfig{ID: "rs1", Name: "rs1"}
	edited.Birdwatcher.API = "http://192.0.2.9:29184"
	s.Reconcile(ctx, &config.Config{
		Sources: []*config.SourceConfig{edited},
	})

	if err := s.refreshSuccess(stale); err != nil {
		t.Fatal(err)
	}
	if initialized, _ := s.sources.IsInitialized("rs1"); initialized {
		t.Error("a stale refresh marked the replaced source initialized")
	}
}

// TestNeighborsStoreRefreshSuccessIgnoresReplacedSource is the same
// check for the neighbors store, which compares source instances.
func TestNeighborsStoreRefreshSuccessIgnoresReplacedSource(t *testing.T) {
	ctx := context.Background()
	source := func(api string) *config.SourceConfig {
		src := &config.SourceConfig{
			ID:      "rs1",
			Name:    "rs1",
			Backend: config.SourceBackendBirdwatcher,
		}
		src.Birdwatcher.API = api
		src.Birdwatcher.Type = "multi_table"
		return src
	}
	s := NewNeighborsStore(&config.Config{
		Sources: []*config.SourceConfig{source("http://192.0.2.1:29184")},
	}, memory.NewNeighborsBackend())
	stale := s.sources.GetInstance("rs1")
	if stale == nil {
		t.Fatal("the test source has no instance")
	}

	s.Reconcile(ctx, &config.Config{
		Sources: []*config.SourceConfig{source("http://192.0.2.9:29184")},
	})

	if err := s.refreshSuccess(stale, "rs1"); err != nil {
		t.Fatal(err)
	}
	if initialized, _ := s.sources.IsInitialized("rs1"); initialized {
		t.Error("a stale refresh marked the replaced source initialized")
	}
}
