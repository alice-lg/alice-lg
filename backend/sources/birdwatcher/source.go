package birdwatcher

import (
	"log"

	"github.com/ecix/alice-lg/backend/sources"
)

type Birdwatcher struct {
	config *Config
}

func NewBirdwatcher(config *Config) *Birdwatcher {
	birdwatcher := &Birdwatcher{
		config: config,
	}
	return birdwatcher
}

func (self *Birdwatcher) Status() {
	log.Println("Implement me: Status()")
}

func (self *Birdwatcher) Neighbours() {
	log.Println("Implement me: Neighbours()")
}

func (self *Birdwatcher) Routes() {
	log.Println("Implement me: Routes()")
}
