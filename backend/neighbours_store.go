package main

import (
	"github.com/ecix/alice-lg/backend/api"
	"github.com/ecix/alice-lg/backend/sources"

	"log"
	"time"
)

type NeighboursIndex map[string]api.Neighbour

type NeighboursStore struct {
	neighboursMap map[sources.Source]NeighboursIndex
	configMap     map[sources.Source]SourceConfig
	statusMap     map[sources.Source]StoreStatus
}

func NewNeighboursStore(config *Config) *NeighboursStore {

	// Build source mapping
	neighboursMap := make(map[sources.Source]NeighboursIndex)
	configMap := make(map[sources.Source]SourceConfig)
	statusMap := make(map[sources.Source]StoreStatus)

	for _, source := range config.Sources {
		instance := source.getInstance()
		configMap[instance] = source
		statusMap[instance] = StoreStatus{
			State: STATE_INIT,
		}

		neighboursMap[instance] = make(NeighboursIndex)
	}

	store := &NeighboursStore{
		neighboursMap: neighboursMap,
		statusMap:     statusMap,
		configMap:     configMap,
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
	for source, _ := range self.neighboursMap {
		// Get current state
		if self.statusMap[source].State == STATE_UPDATING {
			continue // nothing to do here. really.
		}

		// Start updating
		self.statusMap[source] = StoreStatus{
			State: STATE_UPDATING,
		}

		neighboursRes, err := source.Neighbours()
		neighbours := neighboursRes.Neighbours
		if err != nil {
			// That's sad.
			self.statusMap[source] = StoreStatus{
				State:       STATE_ERROR,
				LastError:   err,
				LastRefresh: time.Now(),
			}
			continue
		}

		// Update data
		// Make neighbours index
		index := make(NeighboursIndex)
		for _, neighbour := range neighbours {
			index[neighbour.Id] = neighbour
		}

		self.neighboursMap[source] = index
		// Update state
		self.statusMap[source] = StoreStatus{
			LastRefresh: time.Now(),
			State:       STATE_READY,
		}
	}
}

func (self *NeighboursStore) GetNeighbourAt(
	source sources.Source,
	id string,
) api.Neighbour {
	// Lookup neighbour on RS
	neighbours := self.neighboursMap[source]
	log.Println("Fetching neighbour:", id)
	log.Println("neighbour:", neighbours[id])
	log.Println("neighbours:", neighbours)
	return neighbours[id]
}

// Build some stats for monitoring
func (self *NeighboursStore) Stats() NeighboursStoreStats {
	totalNeighbours := 0
	rsStats := []RouteServerNeighboursStats{}

	for source, neighbours := range self.neighboursMap {
		status := self.statusMap[source]
		totalNeighbours += len(neighbours)
		serverStats := RouteServerNeighboursStats{
			Name:       self.configMap[source].Name,
			State:      stateToString(status.State),
			Neighbours: len(neighbours),
			UpdatedAt:  status.LastRefresh,
		}
		rsStats = append(rsStats, serverStats)
	}

	storeStats := NeighboursStoreStats{
		TotalNeighbours: totalNeighbours,
		RouteServers:    rsStats,
	}
	return storeStats
}
