package main

import (
	"log"
	"os"
	"strings"
	"testing"

	"encoding/json"
	"io/ioutil"

	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/alice-lg/alice-lg/backend/sources/birdwatcher"
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
	routesMap := map[int]api.RoutesResponse{
		1: rs1RoutesResponse,
	}

	configMap := map[int]SourceConfig{
		1: SourceConfig{
			Id:   1,
			Name: "rs1.test",
			Type: SOURCE_BIRDWATCHER,

			Birdwatcher: birdwatcher.Config{
				Api:             "http://localhost:2342",
				Timezone:        "UTC",
				ServerTime:      "2006-01-02T15:04:05",
				ServerTimeShort: "2006-01-02",
				ServerTimeExt:   "Mon, 02 Jan 2006 15:04: 05 -0700",
			},
		},
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

	if stats.TotalRoutes.Filtered != 1 {
		t.Error(
			"expected 1 filtered route, got:",
			stats.TotalRoutes.Filtered,
		)
	}
}

func TestLookupPrefix(t *testing.T) {
	startTestNeighboursStore()
	store := makeTestRoutesStore()
	query := "193.200."

	results := store.LookupPrefix(query)

	if len(results) == 0 {
		t.Error("Expected lookup results. None present.")
		return
	}

	// Check results
	for _, route := range results {
		if strings.HasPrefix(route.Network, query) == false {
			t.Error(
				"All network addresses should start with the",
				"queried prefix",
			)
		}
	}

}
