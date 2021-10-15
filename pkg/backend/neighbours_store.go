package backend

import (
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
)

var REGEX_MATCH_ASLOOKUP = regexp.MustCompile(`(?i)^AS(\d+)`)

type NeighborsIndex map[string]*api.Neighbor

type NeighborsStore struct {
	neighborsMap         map[string]NeighborsIndex
	configMap             map[string]*SourceConfig
	statusMap             map[string]StoreStatus
	refreshInterval       time.Duration
	refreshNeighborStatus bool
	lastRefresh           time.Time

	sync.RWMutex
}

// NewNeighborsStore creates a new store for neighbors
func NewNeighborsStore(config *Config) *NeighborsStore {

	// Build source mapping
	neighborsMap := make(map[string]NeighborsIndex)
	configMap := make(map[string]*SourceConfig)
	statusMap := make(map[string]StoreStatus)

	for _, source := range config.Sources {
		id := source.ID
		configMap[id] = source
		statusMap[id] = StoreStatus{
			State: STATE_INIT,
		}

		neighborsMap[id] = make(NeighborsIndex)
	}

	// Set refresh interval, default to 5 minutes when
	// interval is set to 0
	refreshInterval := time.Duration(
		config.Server.NeighborsStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}

	refreshNeighborStatus := config.Server.EnableNeighborsStatusRefresh

	store := &NeighborsStore{
		neighborsMap:         neighborsMap,
		statusMap:             statusMap,
		configMap:             configMap,
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

func (self *NeighborsStore) SourceStatus(sourceId string) StoreStatus {
	self.RLock()
	status := self.statusMap[sourceId]
	self.RUnlock()

	return status
}

// Get state by source Id
func (self *NeighborsStore) SourceState(sourceId string) int {
	status := self.SourceStatus(sourceId)
	return status.State
}

// Update all neighbors
func (self *NeighborsStore) update() {
	successCount := 0
	errorCount := 0
	t0 := time.Now()
	for sourceId, _ := range self.neighborsMap {
		// Get current state
		if self.statusMap[sourceId].State == STATE_UPDATING {
			continue // nothing to do here. really.
		}

		// Start updating
		self.Lock()
		self.statusMap[sourceId] = StoreStatus{
			State: STATE_UPDATING,
		}
		self.Unlock()

		sourceConfig := self.configMap[sourceId]
		source := sourceConfig.getInstance()

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
			self.statusMap[sourceId] = StoreStatus{
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
			index[neighbor.Id] = neighbor
		}

		self.Lock()
		self.neighborsMap[sourceId] = index
		// Update state
		self.statusMap[sourceId] = StoreStatus{
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

func (self *NeighborsStore) GetNeighborsAt(sourceId string) api.Neighbors {
	self.RLock()
	neighborsIdx := self.neighborsMap[sourceId]
	self.RUnlock()

	var neighborsStatus map[string]api.NeighborStatus
	if self.refreshNeighborStatus {
		sourceConfig := self.configMap[sourceId]
		source := sourceConfig.getInstance()

		neighborsStatusData, err := source.NeighborsStatus()
		if err == nil {
			neighborsStatus = make(map[string]api.NeighborStatus, len(neighborsStatusData.Neighbors))

			for _, neighbor := range neighborsStatusData.Neighbors {
				neighborsStatus[neighbor.Id] = *neighbor
			}
		}
	}

	neighbors := make(api.Neighbors, 0, len(neighborsIdx))

	for _, neighbor := range neighborsIdx {
		if self.refreshNeighborStatus {
			if _, ok := neighborsStatus[neighbor.Id]; ok {
				self.Lock()
				neighbor.State = neighborsStatus[neighbor.Id].State
				self.Unlock()
			}
		}

		neighbors = append(neighbors, neighbor)
	}

	return neighbors
}

func (self *NeighborsStore) GetNeighborAt(
	sourceId string,
	id string,
) *api.Neighbor {
	// Lookup neighbor on RS
	self.RLock()
	neighborsIdx := self.neighborsMap[sourceId]
	self.RUnlock()

	return neighborsIdx[id]
}

func (self *NeighborsStore) LookupNeighborsAt(
	sourceId string,
	query string,
) api.Neighbors {
	results := api.Neighbors{}

	self.RLock()
	neighbors := self.neighborsMap[sourceId]
	self.RUnlock()

	asn := -1
	if REGEX_MATCH_ASLOOKUP.MatchString(query) {
		groups := REGEX_MATCH_ASLOOKUP.FindStringSubmatch(query)
		if a, err := strconv.Atoi(groups[1]); err == nil {
			asn = a
		}
	}

	for _, neighbor := range neighbors {
		if asn >= 0 && neighbor.Asn == asn { // only executed if valid AS query is detected
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

	for sourceId, _ := range self.neighborsMap {
		results[sourceId] = self.LookupNeighborsAt(sourceId, query)
	}

	return results
}

/*
 Filter neighbors from a single route server.
*/
func (self *NeighborsStore) FilterNeighborsAt(
	sourceId string,
	filter *api.NeighborFilter,
) api.Neighbors {
	results := []*api.Neighbor{}

	self.RLock()
	neighbors := self.neighborsMap[sourceId]
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
	for sourceId, _ := range self.neighborsMap {
		rsResults := self.FilterNeighborsAt(sourceId, filter)
		results = append(results, rsResults...)
	}

	return results
}

// Build some stats for monitoring
func (self *NeighborsStore) Stats() NeighborsStoreStats {
	totalNeighbors := 0
	rsStats := []RouteServerNeighborsStats{}

	self.RLock()
	for sourceId, neighbors := range self.neighborsMap {
		status := self.statusMap[sourceId]
		totalNeighbors += len(neighbors)
		serverStats := RouteServerNeighborsStats{
			Name:       self.configMap[sourceId].Name,
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
