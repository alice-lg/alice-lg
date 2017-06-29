package main

import (
	"github.com/ecix/alice-lg/backend/api"
	"testing"
)

// Make a store and populate it with data
func makeNeighboursStore() *NeighboursStore {

	// Populate neighbours
	rs1 := NeighboursIndex{
		"ID2233_AS2342": api.Neighbour{
			Id:          "ID2233_AS2342",
			Description: "PEER AS2342 192.9.23.42 Customer Peer 1",
		},
		"ID2233_AS2343": api.Neighbour{
			Id:          "ID2233_AS2343",
			Description: "PEER AS2343 192.9.23.43 Different Peer 1",
		},
		"ID2233_AS2344": api.Neighbour{
			Id:          "ID2233_AS2344",
			Description: "PEER AS2344 192.9.23.44 3rd Peer from the sun",
		},
	}

	rs2 := NeighboursIndex{
		"ID2233_AS2342": api.Neighbour{
			Id:          "ID2233_AS2342",
			Description: "PEER AS2342 192.9.23.42 Customer Peer 1",
		},
		"ID2233_AS4223": api.Neighbour{
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
	}

	return store
}

func TestGetNeighbourAt(t *testing.T) {
	store := makeNeighboursStore()

	neighbour := store.GetNeighbourAt(1, "ID2233_AS2343")
	if neighbour.Id != "ID2233_AS2343" {
		t.Error("Expected another peer in GetNeighbourAt")
	}

}

func TestNeighbourLookupAt(t *testing.T) {
	store := makeNeighboursStore()
	_ = store
}
