package birdwatcher

import (
	"log"

	"github.com/ecix/alice-lg/backend/api"
)

type Birdwatcher struct {
	config Config
	client *Client
}

func NewBirdwatcher(config Config) *Birdwatcher {
	client := NewClient(config.Api)

	birdwatcher := &Birdwatcher{
		config: config,
		client: client,
	}
	return birdwatcher
}

func (self *Birdwatcher) Status() (api.StatusResponse, error) {
	bird, err := self.client.GetJson("/status")
	if err != nil {
		return api.StatusResponse{}, nil
	}

	birdStatus := bird["status"].(map[string]interface{})

	// Get special fields
	serverTime, _ := parseServerTime(
		birdStatus["current_server"],
		SERVER_TIME_SHORT,
		self.config.Timezone,
	)

	lastReboot, _ := parseServerTime(
		birdStatus["last_reboot"],
		SERVER_TIME_SHORT,
		self.config.Timezone,
	)

	lastReconfig, _ := parseServerTime(
		birdStatus["last_reconfig"],
		SERVER_TIME_EXT,
		self.config.Timezone,
	)

	// Make status response
	status := api.Status{
		ServerTime:   serverTime,
		LastReboot:   lastReboot,
		LastReconfig: lastReconfig,
		Backend:      "bird",
		Version:      birdStatus["version"].(string),
		Message:      birdStatus["message"].(string),
		RouterId:     birdStatus["router_id"].(string),
	}

	response := api.StatusResponse{
		Status: status,
	}

	return response, nil
}

func (self *Birdwatcher) Neighbours() (api.NeighboursResponse, error) {
	log.Println("Implement me: Neighbours()")

	return api.NeighboursResponse{}, nil
}

func (self *Birdwatcher) Routes(neighbourId string) (api.RoutesResponse, error) {
	log.Println("Implement me: Routes()")

	return api.RoutesResponse{}, nil
}
