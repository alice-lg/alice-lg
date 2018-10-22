package main

import (
	"github.com/alice-lg/alice-lg/backend/api"

	"sort"
	"testing"
)

/*
 Start the global neighbours store,
 because the route store in the tests have
 this as a dependency.
*/
func startTestNeighboursStore() {
	store := makeTestNeighboursStore()
	AliceNeighboursStore = store
}

/*
 Make a store and populate it with data
*/
func makeTestNeighboursStore() *NeighboursStore {

	// Populate neighbours
	rs1 := NeighboursIndex{
		"ID2233_AS2342": &api.Neighbour{
			Id:          "ID2233_AS2342",
			Description: "PEER AS2342 192.9.23.42 Customer Peer 1",
		},
		"ID2233_AS2343": &api.Neighbour{
			Id:          "ID2233_AS2343",
			Description: "PEER AS2343 192.9.23.43 Different Peer 1",
		},
		"ID2233_AS2344": &api.Neighbour{
			Id:          "ID2233_AS2344",
			Description: "PEER AS2344 192.9.23.44 3rd Peer from the sun",
		},
	}

	rs2 := NeighboursIndex{
		"ID2233_AS2342": &api.Neighbour{
			Id:          "ID2233_AS2342",
			Description: "PEER AS2342 192.9.23.42 Customer Peer 1",
		},
		"ID2233_AS4223": &api.Neighbour{
			Id:          "ID2233_AS4223",
			Description: "PEER AS4223 192.9.42.23 Cloudfoo Inc.",
		},
	}

	// Create store
	store := &NeighboursStore{
		neighboursMap: map[int]NeighboursIndex{
			1: rs1,
			2: rs2,
		},
		statusMap: map[int]StoreStatus{
			1: StoreStatus{
				State: STATE_READY,
			},
			2: StoreStatus{
				State: STATE_INIT,
			},
		},
	}

	return store
}

func TestGetSourceState(t *testing.T) {
	store := makeTestNeighboursStore()

	if store.SourceState(1) != STATE_READY {
		t.Error("Expected Source(1) to be STATE_READY")
	}

	if store.SourceState(2) == STATE_READY {
		t.Error("Expected Source(2) to be NOT STATE_READY")
	}
}

func TestGetNeighbourAt(t *testing.T) {
	store := makeTestNeighboursStore()

	neighbour := store.GetNeighbourAt(1, "ID2233_AS2343")
	if neighbour.Id != "ID2233_AS2343" {
		t.Error("Expected another peer in GetNeighbourAt")
	}
}

func TestGetNeighbors(t *testing.T) {
	store := makeTestNeighboursStore()
	neighbors := store.GetNeighborsAt(2)

	if len(neighbors) != 2 {
		t.Error("Expected 2 neighbors, got:", len(neighbors))
	}

	sort.Sort(neighbors)

	if neighbors[0].Id != "ID2233_AS2342" {
		t.Error("Expected neighbor: ID2233_AS2342, got:",
			neighbors[0])
	}

	neighbors = store.GetNeighborsAt(3)
	if len(neighbors) != 0 {
		t.Error("Unknown source should have yielded zero results")
	}

}

func TestNeighbourLookupAt(t *testing.T) {
	store := makeTestNeighboursStore()

	expected := []string{
		"ID2233_AS2342",
		"ID2233_AS2343",
	}

	neighbours := store.LookupNeighboursAt(1, "peer 1")

	// Make index
	index := NeighboursIndex{}
	for _, n := range neighbours {
		index[n.Id] = n
	}

	for _, id := range expected {
		_, ok := index[id]
		if !ok {
			t.Error("Expected", id, "to be in result set")
		}
	}
}

func TestNeighbourLookup(t *testing.T) {
	store := makeTestNeighboursStore()

	// First result set: "Peer 1"
	_ = store

	results := store.LookupNeighbours("Cloudfoo")

	// Peer should be present at RS2
	neighbours, ok := results[2]
	if !ok {
		t.Error("Lookup on rs2 unsuccessful.")
	}

	if len(neighbours) > 1 {
		t.Error("Lookup should match exact 1 peer.")
	}

	n := neighbours[0]
	if n.Id != "ID2233_AS4223" {
		t.Error("Wrong peer in lookup response")
	}
}
