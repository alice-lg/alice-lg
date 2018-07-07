package main

import (
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/backend/api"
)

type NeighboursIndex map[string]*api.Neighbour

type NeighboursStore struct {
	neighboursMap   map[int]NeighboursIndex
	configMap       map[int]SourceConfig
	statusMap       map[int]StoreStatus
	refreshInterval time.Duration

	sync.RWMutex
}

func NewNeighboursStore(config *Config) *NeighboursStore {

	// Build source mapping
	neighboursMap := make(map[int]NeighboursIndex)
	configMap := make(map[int]SourceConfig)
	statusMap := make(map[int]StoreStatus)

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
	log.Println("Neighbours Store refresh interval set to:", self.refreshInterval, "minutes")
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

func (self *NeighboursStore) update() {
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

		source := self.configMap[sourceId].getInstance()

		neighboursRes, err := source.Neighbours()
		neighbours := neighboursRes.Neighbours
		if err != nil {
			// That's sad.
			self.Lock()
			self.statusMap[sourceId] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			self.Unlock()
			continue
		}

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
	}
}

func (self *NeighboursStore) GetNeighbourAt(
	sourceId int,
	id string,
) *api.Neighbour {
	// Lookup neighbour on RS
	self.RLock()
	neighbours := self.neighboursMap[sourceId]
	self.RUnlock()
	return neighbours[id]
}

func (self *NeighboursStore) LookupNeighboursAt(
	sourceId int,
	query string,
) api.Neighbours {
	results := api.Neighbours{}

	self.RLock()
	neighbours := self.neighboursMap[sourceId]
	self.RUnlock()

	asn := -1
	if regex := regexp.MustCompile(`(?i)^AS(\d+)`); regex.MatchString(query) {
		groups := regex.FindStringSubmatch(query)
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
