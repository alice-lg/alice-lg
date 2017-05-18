package birdwatcher

import (
	"log"

	"github.com/ecix/alice-lg/backend/api"
)

type Birdwatcher struct {
	config Config
}

func NewBirdwatcher(config Config) *Birdwatcher {
	birdwatcher := &Birdwatcher{
		config: config,
	}
	return birdwatcher
}

func (self *Birdwatcher) Status() (api.StatusResponse, error) {
	log.Println("Implement me: Status()")

	return api.StatusResponse{}, nil
}

func (self *Birdwatcher) Neighbours() (api.NeighboursResponse, error) {
	log.Println("Implement me: Neighbours()")

	return api.NeighboursResponse{}, nil
}

func (self *Birdwatcher) Routes(neighbourId string) (api.RoutesResponse, error) {
	log.Println("Implement me: Routes()")

	return api.RoutesResponse{}, nil
}
