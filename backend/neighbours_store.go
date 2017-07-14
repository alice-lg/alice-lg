package main

import (
	"log"
	"sync"
	"time"

	"github.com/ecix/alice-lg/backend/api"
)

type NeighboursIndex map[string]api.Neighbour

type NeighboursStore struct {
	neighboursMap map[int]NeighboursIndex
	configMap     map[int]SourceConfig
	statusMap     map[int]StoreStatus

	rwlock *sync.RWMutex
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

	store := &NeighboursStore{
		neighboursMap: neighboursMap,
		statusMap:     statusMap,
		configMap:     configMap,

		rwlock: &sync.RWMutex{},
	}
	return store
}

func (self *NeighboursStore) Start() {
	log.Println("Starting local neighbours store")
	go self.init()
}

func (self *NeighboursStore) init() {
	// Perform initial update
	self.update()

	// Initial logging
	self.Stats().Log()

	// Periodically update store
	for {
		time.Sleep(5 * time.Minute)
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
		self.rwlock.Lock()
		self.statusMap[sourceId] = StoreStatus{
			State: STATE_UPDATING,
		}
		self.rwlock.Unlock()

		source := self.configMap[sourceId].getInstance()

		neighboursRes, err := source.Neighbours()
		neighbours := neighboursRes.Neighbours
		if err != nil {
			// That's sad.
			self.rwlock.Lock()
			self.statusMap[sourceId] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			self.rwlock.Unlock()
			continue
		}

		// Update data
		// Make neighbours index
		index := make(NeighboursIndex)
		for _, neighbour := range neighbours {
			index[neighbour.Id] = neighbour
		}

		self.rwlock.Lock()
		self.neighboursMap[sourceId] = index
		// Update state
		self.statusMap[sourceId] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
		self.rwlock.Unlock()
	}
}

func (self *NeighboursStore) GetNeighbourAt(
	sourceId int,
	id string,
) api.Neighbour {
	// Lookup neighbour on RS
	self.rwlock.RLock()
	neighbours := self.neighboursMap[sourceId]
	self.rwlock.RUnlock()
	return neighbours[id]
}

func (self *NeighboursStore) LookupNeighboursAt(
	sourceId int,
	query string,
) []api.Neighbour {
	results := []api.Neighbour{}

	self.rwlock.RLock()
	neighbours := self.neighboursMap[sourceId]
	self.rwlock.RUnlock()

	for _, neighbour := range neighbours {
		if !ContainsCi(neighbour.Description, query) {
			continue
		}

		results = append(results, neighbour)
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

	self.rwlock.RLock()
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
	self.rwlock.RUnlock()

	storeStats := NeighboursStoreStats{
		TotalNeighbours: totalNeighbours,
		RouteServers:    rsStats,
	}
	return storeStats
}
