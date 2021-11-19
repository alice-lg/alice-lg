package store

import (
	"context"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
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
	)
}

// NeighborsStore is queryable for neighbor information
type NeighborsStore struct {
	backend NeighborsBackend
	sources *SourcesStore

	refreshInterval      time.Duration
	forceNeighborRefresh bool

	sync.RWMutex
}

// NewNeighborsStore creates a new store for neighbors
func NewNeighborsStore(
	cfg *config.Config,
	sources *SourcesStore,
	backend NeighborsBackend,
) *NeighborsStore {
	// Set refresh interval, default to 5 minutes when
	// interval is set to 0
	refreshInterval := time.Duration(
		cfg.Server.NeighborsStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}

	// Neighbors will be refreshed on every GetNeighborsAt
	// invocation. Why? I (Annika) don't know. I have to ask Patrick.
	// TODO: This feels wrong here. Figure out reason why it
	//       was added and refactor.
	// At least now the variable name is a bit more honest.
	forceNeighborRefresh := config.Server.EnableNeighborsStatusRefresh

	store := &NeighborsStore{
		backend:              backend,
		sources:              sources,
		refreshInterval:      refreshInterval,
		forceNeighborRefresh: forceNeighborRefresh,
	}
	return store
}

// Start the store's housekeeping.
func (s *NeighborsStore) Start() {
	log.Println("Starting local neighbors store")
	log.Println("Neighbors Store refresh interval set to:", s.refreshInterval)
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

// SourceStatus retrievs the status for a route server
// identified by sourceID.
func (s *NeighborsStore) SourceStatus(sourceID string) Status {
	return s.sources.GetStatus(sourceID)
}

// SourceState gets the state by source ID
func (s *NeighborsStore) SourceState(sourceID string) int {
	status := s.SourceStatus(sourceID)
	return status.State
}

// updateSource will update a single source. This
// function may crash or return errors.
func (s *NeighborsStore) updateSource(src sources.Source) error {
	// Get neighbors form source instance and update backend
	res, err := source.Neighbors()
	if err != nil {
		return err
	}
	if err := s.backend.SetNeighbors(id, res); err != nil {
		return err
	}
	return s.sources.RefreshNeighborsSuccess(id)
}

// safeUpdateSource will try to update a source but
// will recover from a panic if something goes wrong.
// In that case, the LastError and State will be updated.
func (s *NeighborsStore) safeUpdateSource(id string) {
	if !s.sources.ShouldRefreshNeighbors(id, s.refreshInterval) {
		return // Nothing to do here
	}

	if err := s.sources.LockSource(id); err != nil {
		log.Println("Cloud not start neighbor refresh:", err)
		return
	}

	// Apply jitter so, we do not hit everything at once.
	// TODO: Make configurable
	time.Sleep(time.Duration(rand.Intn(30) * time.Second))

	src := s.sources.GetSource(id)
	srcName := s.sources.GetName(id)

	// Prepare for impact.
	defer func() {
		if err := recover(); err != nil {
			log.Println(
				"Recovering after failed neighbors refresh of",
				srcName, "from:", err)
			s.sources.NeighborsRefreshError(id, err)
		}
	}()

	if err := s.updateSource(id); err != nil {
		log.Println(
			"Refeshing neighbors of", srcName, "failed:", err)
		s.sources.NeighborsRefreshError(id, err)
	}
}

// Update all neighbors from all sources, where the
// sources last neighbor refresh is longer ago
// than the configured refresh period.
func (s *NeighborsStore) update() {
	for _, id := range s.sources.GetSourceIDs() {
		go safeUpdateSource(id)
	}
}

// GetNeighborsAt gets all neighbors from a routeserver
func (s *NeighborsStore) GetNeighborsAt(sourceID string) (api.Neighbors, error) {
	ctx := context.Background()

	if s.forceNeighborRefresh {
		if err := s.updateSource(sourceID); err != nil {
			return nil, err
		}
	}

	return s.backend.GetNeighborsAt(ctx, sourceID)
}

// GetNeighborAt looks up a neighbor on a RS by ID.
func (s *NeighborsStore) GetNeighborAt(
	sourceID string, neighborID string,
) (*api.Neighbor, error) {
	return s.backend.GetNeighborAt(ctx, sourceID, neighborID)
}

// LookupNeighborsAt filters for neighbors at a route
// server matching a given query string.
func (s *NeighborsStore) LookupNeighborsAt(
	sourceID string,
	query string,
) api.Neighbors {
	results := api.Neighbors{}
	neighbors := s.backend.GetNeighborsAt(sourceID)

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

	return results
}

// LookupNeighbors filters for neighbors matching a query
// on all route servers.
func (s *NeighborsStore) LookupNeighbors(
	query string,
) api.NeighborsLookupResults {
	// Create empty result set
	results := make(api.NeighborsLookupResults)
	for sourceID := range s.neighborsMap {
		results[sourceID] = s.LookupNeighborsAt(sourceID, query)
	}
	return results
}

// FilterNeighborsAt filters neighbors from a single route server.
func (s *NeighborsStore) FilterNeighborsAt(
	sourceID string,
	filter *api.NeighborFilter,
) api.Neighbors {
	results := []*api.Neighbor{}
	s.RLock()
	neighbors := s.neighborsMap[sourceID]
	s.RUnlock()

	// Apply filters
	for _, neighbor := range neighbors {
		if filter.Match(neighbor) {
			results = append(results, neighbor)
		}
	}
	return results
}

// FilterNeighbors retrieves neighbors by name or by ASN
// from all route servers.
func (s *NeighborsStore) FilterNeighbors(
	filter *api.NeighborFilter,
) api.Neighbors {
	results := []*api.Neighbor{}
	// Get neighbors from all routeservers
	for sourceID := range s.neighborsMap {
		rsResults := s.FilterNeighborsAt(sourceID, filter)
		results = append(results, rsResults...)
	}
	return results
}

// Stats exports some statistics for monitoring.
func (s *NeighborsStore) Stats() *api.NeighborsStoreStats {
	totalNeighbors := 0
	rsStats := []api.RouteServerNeighborsStats{}

	s.RLock()
	for sourceID, neighbors := range s.neighborsMap {
		status := s.statusMap[sourceID]
		totalNeighbors += len(neighbors)
		serverStats := api.RouteServerNeighborsStats{
			Name:      s.cfgMap[sourceID].Name,
			State:     stateToString(status.State),
			Neighbors: len(neighbors),
			UpdatedAt: status.LastRefresh,
		}
		rsStats = append(rsStats, serverStats)
	}
	s.RUnlock()

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
	return status.LastNeighborsRefresh
}

// SourceCacheTTL returns the next time when a refresh
// will be started.
func (s *NeighborsStore) SourceCacheTTL(sourceID string) time.Time {
	lastRefresh := s.SourceCachedAt(sourceID)
	return lastRefresh.Add(s.refreshInterval)
}
