package birdwatcher

import (
	"fmt"

	"github.com/ecix/alice-lg/backend/api"
)

/*
Helper functions for dealing with birdwatcher API data
*/

// Get neighbour by protocol id
func getNeighbourById(neighbours api.Neighbours, id string) (api.Neighbour, error) {
	for _, n := range neighbours {
		if n.Id == id {
			return n, nil
		}
	}
	unknown := api.Neighbour{
		Id:          "unknown",
		Description: "Unknown neighbour",
	}
	return unknown, fmt.Errorf("Neighbour not found")
}
