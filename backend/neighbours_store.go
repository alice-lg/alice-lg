package main

import (
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/backend/api"
)

var REGEX_MATCH_ASLOOKUP = regexp.MustCompile(`(?i)^AS(\d+)`)

type NeighboursIndex map[string]*api.Neighbour

type NeighboursStore struct {
	neighboursMap   map[string]NeighboursIndex
	configMap       map[string]*SourceConfig
	statusMap       map[string]StoreStatus
	refreshInterval time.Duration

	sync.RWMutex
}

func NewNeighboursStore(config *Config) *NeighboursStore {

	// Build source mapping
	neighboursMap := make(map[string]NeighboursIndex)
	configMap := make(map[string]*SourceConfig)
	statusMap := make(map[string]StoreStatus)

	for _, source := range config.Sources {
		sourceId := source.Id
		configMap[sourceId] = source
		statusMap[sourceId] = StoreStatus{
			State: STATE_INIT,
		}

		neighboursMap[sourceId] = make(NeighboursIndex)
	}

	// Set refresh interval, default to 5 minutes when
	// interval is set to 0
	refreshInterval := time.Duration(
		config.Server.NeighboursStoreRefreshInterval) * time.Minute
	if refreshInterval == 0 {
		refreshInterval = time.Duration(5) * time.Minute
	}

	store := &NeighboursStore{
		neighboursMap:   neighboursMap,
		statusMap:       statusMap,
		configMap:       configMap,
		refreshInterval: refreshInterval,
	}
	return store
}

func (self *NeighboursStore) Start() {
	log.Println("Starting local neighbours store")
	log.Println("Neighbours Store refresh interval set to:", self.refreshInterval)
	go self.init()
}

func (self *NeighboursStore) init() {
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

func (self *NeighboursStore) SourceStatus(sourceId string) StoreStatus {
	self.RLock()
	status := self.statusMap[sourceId]
	self.RUnlock()

	return status
}

// Get state by source Id
func (self *NeighboursStore) SourceState(sourceId string) int {
	status := self.SourceStatus(sourceId)
	return status.State
}

// Update all neighbors
func (self *NeighboursStore) update() {
	successCount := 0
	errorCount := 0
	t0 := time.Now()
	for sourceId, _ := range self.neighboursMap {
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

		neighboursRes, err := source.Neighbours()
		if err != nil {
			log.Println(
				"Refreshing the neighbors store failed for:",
				sourceConfig.Name, "(", sourceConfig.Id, ")",
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

		neighbours := neighboursRes.Neighbours

		// Update data
		// Make neighbours index
		index := make(NeighboursIndex)
		for _, neighbour := range neighbours {
			index[neighbour.Id] = neighbour
		}

		self.Lock()
		self.neighboursMap[sourceId] = index
		// Update state
		self.statusMap[sourceId] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
		self.Unlock()
		successCount++
	}

	refreshDuration := time.Since(t0)
	log.Println(
		"Refreshed neighbors store for", successCount, "of", successCount+errorCount,
		"sources with", errorCount, "error(s) in", refreshDuration,
	)
}

func (self *NeighboursStore) GetNeighborsAt(sourceId string) api.Neighbours {
	self.RLock()
	neighborsIdx := self.neighboursMap[sourceId]
	self.RUnlock()

	neighbors := make(api.Neighbours, 0, len(neighborsIdx))

	for _, neighbor := range neighborsIdx {
		neighbors = append(neighbors, neighbor)
	}

	return neighbors
}

func (self *NeighboursStore) GetNeighbourAt(
	sourceId string,
	id string,
) *api.Neighbour {
	// Lookup neighbour on RS
	self.RLock()
	neighborsIdx := self.neighboursMap[sourceId]
	self.RUnlock()

	return neighborsIdx[id]
}

func (self *NeighboursStore) LookupNeighboursAt(
	sourceId string,
	query string,
) api.Neighbours {
	results := api.Neighbours{}

	self.RLock()
	neighbours := self.neighboursMap[sourceId]
	self.RUnlock()

	asn := -1
	if REGEX_MATCH_ASLOOKUP.MatchString(query) {
		groups := REGEX_MATCH_ASLOOKUP.FindStringSubmatch(query)
		if a, err := strconv.Atoi(groups[1]); err == nil {
			asn = a
		}
	}

	for _, neighbour := range neighbours {
		if asn >= 0 && neighbour.Asn == asn { // only executed if valid AS query is detected
			results = append(results, neighbour)
		} else if ContainsCi(neighbour.Description, query) {
			results = append(results, neighbour)
		} else {
			continue
		}
	}

	return results
}

func (self *NeighboursStore) LookupNeighbours(
	query string,
) api.NeighboursLookupResults {
	// Create empty result set
	results := make(api.NeighboursLookupResults)

	for sourceId, _ := range self.neighboursMap {
		results[sourceId] = self.LookupNeighboursAt(sourceId, query)
	}

	return results
}

// Build some stats for monitoring
func (self *NeighboursStore) Stats() NeighboursStoreStats {
	totalNeighbours := 0
	rsStats := []RouteServerNeighboursStats{}

	self.RLock()
	for sourceId, neighbours := range self.neighboursMap {
		status := self.statusMap[sourceId]
		totalNeighbours += len(neighbours)
		serverStats := RouteServerNeighboursStats{
			Name:       self.configMap[sourceId].Name,
			State:      stateToString(status.State),
			Neighbours: len(neighbours),
			UpdatedAt:  status.LastRefresh,
		}
		rsStats = append(rsStats, serverStats)
	}
	self.RUnlock()

	storeStats := NeighboursStoreStats{
		TotalNeighbours: totalNeighbours,
		RouteServers:    rsStats,
	}
	return storeStats
}
