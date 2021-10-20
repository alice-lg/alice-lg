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

var REGEX_MATCH_ASLOOKUP = regexp.MustCompile(`(?i)^AS(\d+)`)

type NeighborsIndex map[string]*api.Neighbor

type NeighborsStore struct {
	neighborsMap         map[string]NeighborsIndex
	cfgMap             map[string]*config.SourceConfig
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
		neighborsMap:         neighborsMap,
		statusMap:             statusMap,
		cfgMap:             cfgMap,
		refreshInterval:       refreshInterval,
		refreshNeighborStatus: refreshNeighborStatus,
	}
	return store
}

// Start the store's housekeeping.
func (self *NeighborsStore) Start() {
	log.Println("Starting local neighbors store")
	log.Println("Neighbors Store refresh interval set to:", self.refreshInterval)
	go self.init()
}

func (self *NeighborsStore) init() {
	// Perform initial update
	self.update()

	// Initial logging
	self.Stats().Log()

	// Periodically update store
	for {
		time.Sleep(self.refreshInterval)
		self.update()
	}
}

func (self *NeighborsStore) SourceStatus(sourceID string) StoreStatus {
	self.RLock()
	status := self.statusMap[sourceID]
	self.RUnlock()

	return status
}

// Get state by source ID
func (self *NeighborsStore) SourceState(sourceID string) int {
	status := self.SourceStatus(sourceID)
	return status.State
}

// Update all neighbors
func (self *NeighborsStore) update() {
	successCount := 0
	errorCount := 0
	t0 := time.Now()
	for sourceID, _ := range self.neighborsMap {
		// Get current state
		if self.statusMap[sourceID].State == STATE_UPDATING {
			continue // nothing to do here. really.
		}

		// Start updating
		self.Lock()
		self.statusMap[sourceID] = StoreStatus{
			State: STATE_UPDATING,
		}
		self.Unlock()

		sourceConfig := self.cfgMap[sourceID]
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
			self.Lock()
			self.statusMap[sourceID] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			self.Unlock()

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

		self.Lock()
		self.neighborsMap[sourceID] = index
		// Update state
		self.statusMap[sourceID] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
		self.lastRefresh = time.Now().UTC()
		self.Unlock()
		successCount++
	}

	refreshDuration := time.Since(t0)
	log.Println(
		"Refreshed neighbors store for", successCount, "of", successCount+errorCount,
		"sources with", errorCount, "error(s) in", refreshDuration,
	)
}

func (self *NeighborsStore) GetNeighborsAt(sourceID string) api.Neighbors {
	self.RLock()
	neighborsIDx := self.neighborsMap[sourceID]
	self.RUnlock()

	var neighborsStatus map[string]api.NeighborStatus
	if self.refreshNeighborStatus {
		sourceConfig := self.cfgMap[sourceID]
		source := sourceConfig.GetInstance()

		neighborsStatusData, err := source.NeighborsStatus()
		if err == nil {
			neighborsStatus = make(map[string]api.NeighborStatus, len(neighborsStatusData.Neighbors))

			for _, neighbor := range neighborsStatusData.Neighbors {
				neighborsStatus[neighbor.ID] = *neighbor
			}
		}
	}

	neighbors := make(api.Neighbors, 0, len(neighborsIDx))

	for _, neighbor := range neighborsIDx {
		if self.refreshNeighborStatus {
			if _, ok := neighborsStatus[neighbor.ID]; ok {
				self.Lock()
				neighbor.State = neighborsStatus[neighbor.ID].State
				self.Unlock()
			}
		}

		neighbors = append(neighbors, neighbor)
	}

	return neighbors
}

func (self *NeighborsStore) GetNeighborAt(
	sourceID string,
	id string,
) *api.Neighbor {
	// Lookup neighbor on RS
	self.RLock()
	neighborsIDx := self.neighborsMap[sourceID]
	self.RUnlock()

	return neighborsIDx[id]
}

func (self *NeighborsStore) LookupNeighborsAt(
	sourceID string,
	query string,
) api.Neighbors {
	results := api.Neighbors{}

	self.RLock()
	neighbors := self.neighborsMap[sourceID]
	self.RUnlock()

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

func (self *NeighborsStore) LookupNeighbors(
	query string,
) api.NeighborsLookupResults {
	// Create empty result set
	results := make(api.NeighborsLookupResults)

	for sourceID, _ := range self.neighborsMap {
		results[sourceID] = self.LookupNeighborsAt(sourceID, query)
	}

	return results
}

/*
 Filter neighbors from a single route server.
*/
func (self *NeighborsStore) FilterNeighborsAt(
	sourceID string,
	filter *api.NeighborFilter,
) api.Neighbors {
	results := []*api.Neighbor{}

	self.RLock()
	neighbors := self.neighborsMap[sourceID]
	self.RUnlock()

	// Apply filters
	for _, neighbor := range neighbors {
		if filter.Match(neighbor) {
			results = append(results, neighbor)
		}
	}
	return results
}

/*
 Filter neighbors by name or by ASN.
 Collect results from all routeservers.
*/
func (self *NeighborsStore) FilterNeighbors(
	filter *api.NeighborFilter,
) api.Neighbors {
	results := []*api.Neighbor{}

	// Get neighbors from all routeservers
	for sourceID, _ := range self.neighborsMap {
		rsResults := self.FilterNeighborsAt(sourceID, filter)
		results = append(results, rsResults...)
	}

	return results
}

// Build some stats for monitoring
func (self *NeighborsStore) Stats() NeighborsStoreStats {
	totalNeighbors := 0
	rsStats := []RouteServerNeighborsStats{}

	self.RLock()
	for sourceID, neighbors := range self.neighborsMap {
		status := self.statusMap[sourceID]
		totalNeighbors += len(neighbors)
		serverStats := RouteServerNeighborsStats{
			Name:       self.cfgMap[sourceID].Name,
			State:      stateToString(status.State),
			Neighbors: len(neighbors),
			UpdatedAt:  status.LastRefresh,
		}
		rsStats = append(rsStats, serverStats)
	}
	self.RUnlock()

	storeStats := NeighborsStoreStats{
		TotalNeighbors: totalNeighbors,
		RouteServers:    rsStats,
	}
	return storeStats
}

func (self *NeighborsStore) CachedAt() time.Time {
	return self.lastRefresh
}

func (self *NeighborsStore) CacheTTL() time.Time {
	return self.lastRefresh.Add(self.refreshInterval)
}
