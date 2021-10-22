package store

import (
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
)

// ReMatchASLookup matches lookups with an 'AS' prefix
var ReMatchASLookup = regexp.MustCompile(`(?i)^AS(\d+)`)

// NeighborsIndex is a mapping from a string to a neighbor.
type NeighborsIndex map[string]*api.Neighbor

// NeighborsStore is queryable for neighbor information
type NeighborsStore struct {
	neighborsMap          map[string]NeighborsIndex
	cfgMap                map[string]*config.SourceConfig
	statusMap             map[string]StoreStatus
	refreshInterval       time.Duration
	refreshNeighborStatus bool
	lastRefresh           time.Time

	sync.RWMutex
}

// NewNeighborsStore creates a new store for neighbors
func NewNeighborsStore(cfg *config.Config) *NeighborsStore {
	// Build source mapping
	neighborsMap := make(map[string]NeighborsIndex)
	cfgMap := make(map[string]*config.SourceConfig)
	statusMap := make(map[string]StoreStatus)

	for _, source := range cfg.Sources {
		id := source.ID
		cfgMap[id] = source
		statusMap[id] = StoreStatus{
			State: STATE_INIT,
		}

		neighborsMap[id] = make(NeighborsIndex)
	}

	// Set refresh interval, default to 5 minutes when
	// interval is set to 0
	refreshInterval := time.Duration(
		cfg.Server.NeighborsStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}

	refreshNeighborStatus := cfg.Server.EnableNeighborsStatusRefresh

	store := &NeighborsStore{
		neighborsMap:          neighborsMap,
		statusMap:             statusMap,
		cfgMap:                cfgMap,
		refreshInterval:       refreshInterval,
		refreshNeighborStatus: refreshNeighborStatus,
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
	// Perform initial update
	s.update()

	// Initial logging
	s.Stats().Log()

	// Periodically update store
	for {
		time.Sleep(s.refreshInterval)
		s.update()
	}
}

// SourceStatus retrievs the status for a route server
// identified by sourceID.
func (s *NeighborsStore) SourceStatus(sourceID string) StoreStatus {
	s.RLock()
	defer s.RUnlock()
	return s.statusMap[sourceID]
}

// SourceState gets the state by source ID
func (s *NeighborsStore) SourceState(sourceID string) int {
	status := s.SourceStatus(sourceID)
	return status.State
}

// Update all neighbors
func (s *NeighborsStore) update() {
	successCount := 0
	errorCount := 0
	t0 := time.Now()
	for sourceID := range s.neighborsMap {
		// Get current state
		if s.statusMap[sourceID].State == STATE_UPDATING {
			continue // nothing to do here. really.
		}

		// Start updating
		s.Lock()
		s.statusMap[sourceID] = StoreStatus{
			State: STATE_UPDATING,
		}
		s.Unlock()

		sourceConfig := s.cfgMap[sourceID]
		source := sourceConfig.GetInstance()

		neighborsRes, err := source.Neighbors()
		if err != nil {
			log.Println(
				"Refreshing the neighbors store failed for:",
				sourceConfig.Name, "(", sourceConfig.ID, ")",
				"with:", err,
				"- NEXT STATE: ERROR",
			)
			// That's sad.
			s.Lock()
			s.statusMap[sourceID] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			s.Unlock()

			errorCount++
			continue
		}

		neighbors := neighborsRes.Neighbors

		// Update data
		// Make neighbors index
		index := make(NeighborsIndex)
		for _, neighbor := range neighbors {
			index[neighbor.ID] = neighbor
		}

		s.Lock()
		s.neighborsMap[sourceID] = index
		// Update state
		s.statusMap[sourceID] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
		s.lastRefresh = time.Now().UTC()
		s.Unlock()
		successCount++
	}

	refreshDuration := time.Since(t0)
	log.Println(
		"Refreshed neighbors store for", successCount, "of", successCount+errorCount,
		"sources with", errorCount, "error(s) in", refreshDuration,
	)
}

// GetNeighborsAt gets all neighbors from a routeserver
func (s *NeighborsStore) GetNeighborsAt(sourceID string) api.Neighbors {
	s.RLock()
	neighborsIDx := s.neighborsMap[sourceID]
	s.RUnlock()

	var neighborsStatus map[string]api.NeighborStatus
	if s.refreshNeighborStatus {
		sourceConfig := s.cfgMap[sourceID]
		source := sourceConfig.GetInstance()

		neighborsStatusData, err := source.NeighborsStatus()
		if err == nil {
			neighborsStatus = make(
				map[string]api.NeighborStatus,
				len(neighborsStatusData.Neighbors))

			for _, neighbor := range neighborsStatusData.Neighbors {
				neighborsStatus[neighbor.ID] = *neighbor
			}
		}
	}

	neighbors := make(api.Neighbors, 0, len(neighborsIDx))
	for _, neighbor := range neighborsIDx {
		if s.refreshNeighborStatus {
			if _, ok := neighborsStatus[neighbor.ID]; ok {
				s.Lock()
				neighbor.State = neighborsStatus[neighbor.ID].State
				s.Unlock()
			}
		}
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

// GetNeighborAt looks up a neighbor on a RS by ID.
func (s *NeighborsStore) GetNeighborAt(
	sourceID string,
	id string,
) *api.Neighbor {
	s.RLock()
	defer s.RUnlock()
	neighborsIDx := s.neighborsMap[sourceID]
	return neighborsIDx[id]
}

// LookupNeighborsAt filters for neighbors at a route
// server matching a given query string.
func (s *NeighborsStore) LookupNeighborsAt(
	sourceID string,
	query string,
) api.Neighbors {
	results := api.Neighbors{}

	s.RLock()
	neighbors := s.neighborsMap[sourceID]
	s.RUnlock()

	asn := -1
	if REGEX_MATCH_ASLOOKUP.MatchString(query) {
		groups := REGEX_MATCH_ASLOOKUP.FindStringSubmatch(query)
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
func (s *NeighborsStore) Stats() *NeighborsStoreStats {
	totalNeighbors := 0
	rsStats := []RouteServerNeighborsStats{}

	s.RLock()
	for sourceID, neighbors := range s.neighborsMap {
		status := s.statusMap[sourceID]
		totalNeighbors += len(neighbors)
		serverStats := RouteServerNeighborsStats{
			Name:      s.cfgMap[sourceID].Name,
			State:     stateToString(status.State),
			Neighbors: len(neighbors),
			UpdatedAt: status.LastRefresh,
		}
		rsStats = append(rsStats, serverStats)
	}
	s.RUnlock()

	storeStats := &NeighborsStoreStats{
		TotalNeighbors: totalNeighbors,
		RouteServers:   rsStats,
	}
	return storeStats
}

// CachedAt returns the last time the store content
// was refreshed.
func (s *NeighborsStore) CachedAt() time.Time {
	return s.lastRefresh
}

// CacheTTL returns the next time when a refresh
// will be started.
func (s *NeighborsStore) CacheTTL() time.Time {
	return s.lastRefresh.Add(s.refreshInterval)
}
