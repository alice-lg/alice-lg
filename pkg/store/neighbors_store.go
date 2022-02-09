package store

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// ReMatchASLookup matches lookups with an 'AS' prefix
var ReMatchASLookup = regexp.MustCompile(`(?i)^AS(\d+)`)

// NeighborsStoreBackend interface
type NeighborsStoreBackend interface {
	// SetNeighbors replaces all neighbors for a given
	// route server identified by sourceID.
	SetNeighbors(
		ctx context.Context,
		sourceID string,
		neighbors api.Neighbors,
	) error

	// GetNeighborsAt retrieves all neighbors associated
	// with a route server (source).
	GetNeighborsAt(
		ctx context.Context,
		sourceID string,
	) (api.Neighbors, error)

	// GetNeighborsMapAt retrieve a map of neighbor ids
	// to the neighbor for a given route server for quick
	// consecutive lookup.
	GetNeighborsMapAt(
		ctx context.Context,
		sourceID string,
	) (map[string]*api.Neighbor, error)

	// CountNeighborsAt retrieves the current number of
	// stored neighbors.
	CountNeighborsAt(
		ctx context.Context,
		sourceID string,
	) (int, error)
}

// NeighborsStore is queryable for neighbor information
type NeighborsStore struct {
	backend NeighborsStoreBackend
	sources *SourcesStore

	forceNeighborRefresh bool
}

// NewNeighborsStore creates a new store for neighbors
func NewNeighborsStore(
	cfg *config.Config,
	backend NeighborsStoreBackend,
) *NeighborsStore {
	// Set refresh interval, default to 5 minutes when
	// interval is set to 0
	refreshInterval := time.Duration(
		cfg.Server.NeighborsStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}
	refreshParallelism := cfg.Server.NeighborsStoreRefreshParallelism
	if refreshParallelism <= 0 {
		refreshParallelism = 1
	}

	log.Println("Neighbors refresh interval set to:", refreshInterval)
	log.Println("Neighbors refresh parallelism:", refreshParallelism)

	// Store refresh information per store
	sources := NewSourcesStore(cfg, refreshInterval, refreshParallelism)

	// Neighbors will be refreshed on every GetNeighborsAt
	// invocation. Why? I (Annika) don't know. I have to ask Patrick.
	// TODO: This feels wrong here. Figure out reason why it
	//       was added and refactor.
	// At least now the variable name is a bit more honest.
	forceNeighborRefresh := cfg.Server.EnableNeighborsStatusRefresh

	store := &NeighborsStore{
		backend:              backend,
		sources:              sources,
		forceNeighborRefresh: forceNeighborRefresh,
	}
	return store
}

// Start the store's housekeeping.
func (s *NeighborsStore) Start() {
	log.Println("Starting local neighbors store")
	go s.init()
}

func (s *NeighborsStore) init() {
	// Periodically trigger updates. Sources with an
	// LastNeighborsRefresh < refreshInterval or currently
	// updating will be skipped.
	for {
		s.update()
		time.Sleep(time.Second)
	}
}

// GetStatus retrievs the status for a route server
// identified by sourceID.
func (s *NeighborsStore) GetStatus(sourceID string) (*Status, error) {
	return s.sources.GetStatus(sourceID)
}

// IsInitialized retrieves the status for a route server
// and checks if it is ready.
func (s *NeighborsStore) IsInitialized(sourceID string) bool {
	rdy, _ := s.sources.IsInitialized(sourceID)
	return rdy
}

// updateSource will update a single source. This
// function may crash or return errors.
func (s *NeighborsStore) updateSource(
	ctx context.Context,
	src sources.Source,
	srcID string,
) error {
	// Get neighbors form source instance and update backend
	res, err := src.Neighbors()
	if err != nil {
		return err
	}

	if err = s.backend.SetNeighbors(ctx, srcID, res.Neighbors); err != nil {
		return err
	}

	return s.sources.RefreshSuccess(srcID)
}

// safeUpdateSource will try to update a source but
// will recover from a panic if something goes wrong.
// In that case, the LastError and State will be updated.
func (s *NeighborsStore) safeUpdateSource(id string) {
	ctx := context.TODO()

	if !s.sources.ShouldRefresh(id) {
		return // Nothing to do here
	}

	if err := s.sources.LockSource(id); err != nil {
		log.Println("Cloud not start neighbor refresh:", err)
		return
	}

	// Apply jitter so, we do not hit everything at once.
	// TODO: Make configurable
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	src := s.sources.GetInstance(id)
	srcName := s.sources.GetName(id)

	// Prepare for impact.
	defer func() {
		if err := recover(); err != nil {
			log.Println(
				"Recovering after failed neighbors refresh of",
				srcName, "from:", err)
			s.sources.RefreshError(id, err)
		}
	}()

	if err := s.updateSource(ctx, src, id); err != nil {
		log.Println(
			"Refeshing neighbors of", srcName, "failed:", err)
		s.sources.RefreshError(id, err)
	}
}

// Update all neighbors from all sources, where the
// sources last neighbor refresh is longer ago
// than the configured refresh period.
func (s *NeighborsStore) update() {
	for _, id := range s.sources.GetSourceIDsForRefresh() {
		go s.safeUpdateSource(id)
	}
}

// CachedAt returns the time of the oldest partial
// refresh of the dataset.
func (s *NeighborsStore) CachedAt(
	ctx context.Context,
) time.Time {
	return s.sources.CachedAt(ctx)
}

// CacheTTL returns the TTL time
func (s *NeighborsStore) CacheTTL(
	ctx context.Context,
) time.Time {
	return s.sources.NextRefresh(ctx)
}

// GetNeighborsAt gets all neighbors from a routeserver
func (s *NeighborsStore) GetNeighborsAt(
	ctx context.Context,
	sourceID string,
) (api.Neighbors, error) {
	if s.forceNeighborRefresh {
		src := s.sources.GetInstance(sourceID)
		if src == nil {
			return nil, sources.ErrSourceNotFound
		}
		if err := s.updateSource(ctx, src, sourceID); err != nil {
			return nil, err
		}
	}
	return s.backend.GetNeighborsAt(ctx, sourceID)
}

// GetNeighborsMapAt looks up a neighbor on a RS by ID.
func (s *NeighborsStore) GetNeighborsMapAt(
	ctx context.Context,
	sourceID string,
) (map[string]*api.Neighbor, error) {
	return s.backend.GetNeighborsMapAt(ctx, sourceID)
}

// lookupNeighborsAt filters for neighbors at a route
// server matching a given query string.
func (s *NeighborsStore) lookupNeighborsAt(
	ctx context.Context,
	sourceID string,
	query string,
) (api.Neighbors, error) {

	results := api.Neighbors{}
	neighbors, err := s.backend.GetNeighborsAt(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	asn := -1
	if ReMatchASLookup.MatchString(query) {
		groups := ReMatchASLookup.FindStringSubmatch(query)
		if a, err := strconv.Atoi(groups[1]); err == nil {
			asn = a
		}
	}

	for _, neighbor := range neighbors {
		if asn >= 0 && neighbor.ASN == asn { // only executed if valid AS query is detected
			results = append(results, neighbor)
		} else if ContainsCi(neighbor.Description, query) {
			results = append(results, neighbor)
		} else {
			continue
		}
	}

	return results, nil
}

// LookupNeighbors filters for neighbors matching a query
// on all route servers.
func (s *NeighborsStore) LookupNeighbors(
	ctx context.Context,
	query string,
) (api.NeighborsLookupResults, error) {
	// Create empty result set
	results := make(api.NeighborsLookupResults)
	for _, sourceID := range s.sources.GetSourceIDs() {
		neighbors, err := s.lookupNeighborsAt(ctx, sourceID, query)
		if err != nil {
			return nil, err
		}
		results[sourceID] = neighbors
	}
	return results, nil
}

// FilterNeighbors retrieves neighbors by name or by ASN
// from all route servers.
func (s *NeighborsStore) FilterNeighbors(
	ctx context.Context,
	filter *api.NeighborFilter,
) (api.Neighbors, error) {
	results := []*api.Neighbor{}
	// Get neighbors from all routeservers
	for _, sourceID := range s.sources.GetSourceIDs() {
		neighbors, err := s.backend.GetNeighborsAt(ctx, sourceID)
		if err != nil {
			return nil, err
		}
		// Apply filters
		for _, neighbor := range neighbors {
			if filter.Match(neighbor) {
				results = append(results, neighbor)
			}
		}
	}
	return results, nil
}

// Stats exports some statistics for monitoring.
func (s *NeighborsStore) Stats(
	ctx context.Context,
) *api.NeighborsStoreStats {
	totalNeighbors := 0
	rsStats := []api.RouteServerNeighborsStats{}

	for _, sourceID := range s.sources.GetSourceIDs() {
		status, _ := s.sources.GetStatus(sourceID)
		ncount, err := s.backend.CountNeighborsAt(ctx, sourceID)
		if err != nil {
			if !errors.Is(err, sources.ErrSourceNotFound) {
				log.Println("error during neighbor count:", err)
			}
		}
		totalNeighbors += ncount
		serverStats := api.RouteServerNeighborsStats{
			Name:      s.sources.GetName(sourceID),
			State:     status.State.String(),
			Neighbors: ncount,
			UpdatedAt: s.SourceCachedAt(sourceID),
		}
		rsStats = append(rsStats, serverStats)
	}
	storeStats := &api.NeighborsStoreStats{
		TotalNeighbors: totalNeighbors,
		RouteServers:   rsStats,
	}
	return storeStats
}

// SourceCachedAt returns the last time the store content
// was refreshed.
func (s *NeighborsStore) SourceCachedAt(sourceID string) time.Time {
	status, err := s.sources.GetStatus(sourceID)
	if err != nil {
		log.Println("error while getting source cached at:", err)
		return time.Time{}
	}
	return status.LastRefresh
}

// SourceCacheTTL returns the next time when a refresh
// will be started.
func (s *NeighborsStore) SourceCacheTTL(
	ctx context.Context,
	sourceID string,
) time.Time {
	return s.sources.NextRefresh(ctx)
}
