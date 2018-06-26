package main

import (
	"testing"

	"github.com/alice-lg/alice-lg/backend/api"
)

//
// Api Mock Helpers
//

//
// Route Store Tests
//

func makeTestRoutesStore() *RoutesStore {
	rs1RoutesResponse := makeTestRoutesResponse()

	// Build mapping based on source instances:
	//   rs : <response>
	routesMap := make(map[int]api.RoutesResponse)
	statusMap := make(map[int]StoreStatus)
	configMap := make(map[int]SourceConfig)

	store := &RoutesStore{
		routesMap: routesMap,
		statusMap: statusMap,
		configMap: configMap,
	}

	return store
}

func TestRoutesStoreStats(t *testing.T) {

	store := makeTestRoutesStore()
	stats := store.Stats()

	// Check total routes

	// Filtered
	// Imported

}
