package store

import (
	"sort"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)


/*
 Make a store and populate it with data
*/
func makeTestNeighborsStore() *NeighborsStore {

	// Populate neighbors
	rs1 := NeighborsIndex{
		"ID2233_AS2342": &api.Neighbor{
			ID:            "ID2233_AS2342",
			ASN:           2342,
			Description:   "PEER AS2342 192.9.23.42 Customer Peer 1",
			RouteServerID: "rs1",
		},
		"ID2233_AS2343": &api.Neighbor{
			ID:            "ID2233_AS2343",
			ASN:           2343,
			Description:   "PEER AS2343 192.9.23.43 Different Peer 1",
			RouteServerID: "rs1",
		},
		"ID2233_AS2344": &api.Neighbor{
			ID:            "ID2233_AS2344",
			ASN:           2344,
			Description:   "PEER AS2344 192.9.23.44 3rd Peer from the sun",
			RouteServerID: "rs1",
		},
	}

	rs2 := NeighborsIndex{
		"ID2233_AS2342": &api.Neighbor{
			ID:            "ID2233_AS2342",
			ASN:           2342,
			Description:   "PEER AS2342 192.9.23.42 Customer Peer 1",
			RouteServerID: "rs2",
		},
		"ID2233_AS4223": &api.Neighbor{
			ID:            "ID2233_AS4223",
			ASN:           4223,
			Description:   "PEER AS4223 192.9.42.23 Cloudfoo Inc.",
			RouteServerID: "rs2",
		},
	}

	// Create store
	store := &NeighborsStore{
		neighborsMap: map[string]NeighborsIndex{
			"rs1": rs1,
			"rs2": rs2,
		},
		statusMap: map[string]StoreStatus{
			"rs1": StoreStatus{
				State: STATE_READY,
			},
			"rs2": StoreStatus{
				State: STATE_INIT,
			},
		},
	}

	return store
}

func TestGetSourceState(t *testing.T) {
	store := makeTestNeighborsStore()

	if store.SourceState("rs1") != STATE_READY {
		t.Error("Expected Source(1) to be STATE_READY")
	}

	if store.SourceState("rs2") == STATE_READY {
		t.Error("Expected Source(2) to be NOT STATE_READY")
	}
}

func TestGetNeighborAt(t *testing.T) {
	store := makeTestNeighborsStore()

	neighbor := store.GetNeighborAt("rs1", "ID2233_AS2343")
	if neighbor.ID != "ID2233_AS2343" {
		t.Error("Expected another peer in GetNeighborAt")
	}
}

func TestGetNeighbors(t *testing.T) {
	store := makeTestNeighborsStore()
	neighbors := store.GetNeighborsAt("rs2")

	if len(neighbors) != 2 {
		t.Error("Expected 2 neighbors, got:", len(neighbors))
	}

	sort.Sort(neighbors)

	if neighbors[0].ID != "ID2233_AS2342" {
		t.Error("Expected neighbor: ID2233_AS2342, got:",
			neighbors[0])
	}

	neighbors = store.GetNeighborsAt("rs3")
	if len(neighbors) != 0 {
		t.Error("Unknown source should have yielded zero results")
	}

}

func TestNeighborLookupAt(t *testing.T) {
	store := makeTestNeighborsStore()

	expected := []string{
		"ID2233_AS2342",
		"ID2233_AS2343",
	}

	neighbors := store.LookupNeighborsAt("rs1", "peer 1")

	// Make index
	index := NeighborsIndex{}
	for _, n := range neighbors {
		index[n.ID] = n
	}

	for _, id := range expected {
		_, ok := index[id]
		if !ok {
			t.Error("Expected", id, "to be in result set")
		}
	}
}

func TestNeighborLookup(t *testing.T) {
	store := makeTestNeighborsStore()

	// First result set: "Peer 1"
	_ = store

	results := store.LookupNeighbors("Cloudfoo")

	// Peer should be present at RS2
	neighbors, ok := results["rs2"]
	if !ok {
		t.Error("Lookup on rs2 unsuccessful.")
	}

	if len(neighbors) > 1 {
		t.Error("Lookup should match exact 1 peer.")
	}

	n := neighbors[0]
	if n.ID != "ID2233_AS4223" {
		t.Error("Wrong peer in lookup response")
	}
}

func TestNeighborFilter(t *testing.T) {
	store := makeTestNeighborsStore()
	filter := api.NeighborFilterFromQueryString("asn=2342")
	neighbors := store.FilterNeighbors(filter)
	if len(neighbors) != 2 {
		t.Error("Expected two results")
	}

	filter = api.NeighborFilterFromQueryString("")
	neighbors = store.FilterNeighbors(filter)
	if len(neighbors) != 0 {
		t.Error("Expected empty result set")
	}

}
