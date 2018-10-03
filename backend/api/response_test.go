package api

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStatusResponseSerialization(t *testing.T) {

	// Make status
	response := StatusResponse{
		Api: ApiStatus{
			Version:         "2.0.0",
			CacheStatus:     CacheStatus{},
			ResultFromCache: false,
			Ttl:             time.Now(),
		},
		Status: Status{
			Message:  "Server is up and running",
			RouterId: "testrouter",
			Version:  "1.6.3",
			Backend:  "birdwatcher",
		},
	}

	result, err := json.Marshal(response)
	if err != nil {
		t.Error(err)
	}

	_ = result
}

func TestNeighbourSerialization(t *testing.T) {

	// Original backend response
	payload := `{
        "action": "restart", 
        "bgp_state": "Established", 
        "bird_protocol": "BGP", 
        "connection": "Established"
    }`

	details := make(map[string]interface{})
	err := json.Unmarshal([]byte(payload), &details)

	if err != nil {
		t.Error(err)
	}

	// Make neighbour
	neighbour := Neighbour{
		Id:          "PROTOCOL_23_42_",
		State:       "Established",
		Description: "Some peer",
		Address:     "10.10.10.1",
		Details:     details,
	}

	result, err := json.Marshal(neighbour)
	if err != nil {
		t.Error(err)
	}

	_ = result
}
