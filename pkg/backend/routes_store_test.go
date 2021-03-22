package backend

import (
	"log"
	"os"
	"strings"
	"testing"

	"encoding/json"
	"io/ioutil"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/sources/birdwatcher"
)

//
// Api Tets Helpers
//
func loadTestRoutesResponse() *api.RoutesResponse {
	file, err := os.Open("../../testdata/api/routes_response.json")
	if err != nil {
		log.Panic("could not load test data:", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panic("could not read test data:", err)
	}

	response := &api.RoutesResponse{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Panic("could not unmarshal response test data:", err)
	}

	return response
}

/*
 Check for presence of network in result set
*/
func testCheckPrefixesPresence(prefixes, resultset []string, t *testing.T) {
	// Check prefixes
	presence := map[string]bool{}
	for _, prefix := range prefixes {
		presence[prefix] = false
	}

	for _, prefix := range resultset {
		// Check if prefixes are all accounted for
		for net := range presence {
			if prefix == net {
				presence[net] = true
			}
		}
	}

	for net, present := range presence {
		if present == false {
			t.Error(net, "not found in result set")
		}
	}
}

//
// Route Store Tests
//

func makeTestRoutesStore() *RoutesStore {
	rs1RoutesResponse := loadTestRoutesResponse()

	// Build mapping based on source instances:
	//   rs : <response>
	statusMap := make(map[string]StoreStatus)
	routesMap := map[string]*api.RoutesResponse{
		"rs1": rs1RoutesResponse,
	}

	configMap := map[string]*SourceConfig{
		"rs1": &SourceConfig{
			Id:   "rs1",
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

func TestLookupPrefixAt(t *testing.T) {
	startTestNeighboursStore()
	store := makeTestRoutesStore()

	query := "193.200."
	results := store.LookupPrefixAt("rs1", query)

	prefixes := <-results

	// Check results
	for _, prefix := range prefixes {
		if strings.HasPrefix(prefix.Network, query) == false {
			t.Error(
				"All network addresses should start with the",
				"queried prefix",
			)
		}
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
	for _, prefix := range results {
		if strings.HasPrefix(prefix.Network, query) == false {
			t.Error(
				"All network addresses should start with the",
				"queried prefix",
			)
		}
	}
}

func TestLookupNeighboursPrefixesAt(t *testing.T) {
	startTestNeighboursStore()
	store := makeTestRoutesStore()

	// Query
	results := store.LookupNeighboursPrefixesAt("rs1", []string{
		"ID163_AS31078",
	})

	// Check prefixes
	presence := []string{
		"193.200.230.0/24", "193.34.24.0/22", "31.220.136.0/21",
	}

	resultset := []string{}
	for _, prefix := range <-results {
		resultset = append(resultset, prefix.Network)
	}

	testCheckPrefixesPresence(presence, resultset, t)
}

func TestLookupPrefixForNeighbours(t *testing.T) {
	// Construct a neighbours lookup result
	neighbours := api.NeighboursLookupResults{
		"rs1": api.Neighbours{
			&api.Neighbour{
				Id: "ID163_AS31078",
			},
		},
	}

	startTestNeighboursStore()
	store := makeTestRoutesStore()

	// Query
	results := store.LookupPrefixForNeighbours(neighbours)

	// We should have retrived 8 prefixes,
	if len(results) != 8 {
		t.Error("Expected result lenght: 8, got:", len(results))
	}

	presence := []string{
		"193.200.230.0/24", "193.34.24.0/22", "31.220.136.0/21",
	}

	resultset := []string{}
	for _, prefix := range results {
		resultset = append(resultset, prefix.Network)
	}

	testCheckPrefixesPresence(presence, resultset, t)
}
