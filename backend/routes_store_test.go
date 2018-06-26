package main

import (
	"log"
	"os"
	"testing"

	"encoding/json"
	"io/ioutil"

	"github.com/alice-lg/alice-lg/backend/api"
)

//
// Api Tets Helpers
//
func loadTestRoutesResponse() api.RoutesResponse {
	file, err := os.Open("testdata/api/routes_response.json")
	if err != nil {
		log.Panic("could not load test data:", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panic("could not read test data:", err)
	}

	response := api.RoutesResponse{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Panic("could not unmarshal response test data:", err)
	}

	return response
}

//
// Route Store Tests
//

func makeTestRoutesStore() *RoutesStore {
	rs1RoutesResponse := loadTestRoutesResponse()

	// Build mapping based on source instances:
	//   rs : <response>
	statusMap := make(map[int]StoreStatus)
	configMap := make(map[int]SourceConfig)
	routesMap := map[int]api.RoutesResponse{
		1: rs1RoutesResponse,
	}

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
	// There should be 8 imported, and 1 filtered route
	if stats.TotalRoutes.Imported != 8 {
		t.Error(
			"expected 8 imported routes, got:",
			stats.TotalRoutes.Imported,
		)
	}

	if stats.TotalRoutes.Filtered {
		t.Error(
			"expected 1 filtered route, got:",
			stats.TotalRoutes.Filtered,
		)
	}
}
