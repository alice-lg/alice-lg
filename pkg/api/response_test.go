package api

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStatusResponseSerialization(t *testing.T) {

	// Make status
	response := StatusResponse{
		Response: Response{
			Meta: &Meta{
				Version:         "2.0.0",
				CacheStatus:     CacheStatus{},
				ResultFromCache: false,
				TTL:             time.Now(),
			},
		},
		Status: Status{
			Message:  "Server is up and running",
			RouterID: "testrouter",
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

func TestNeighborSerialization(t *testing.T) {

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

	// Make neighbor
	neighbor := Neighbor{
		ID:          "PROTOCOL_23_42_",
		State:       "Established",
		Description: "Some peer",
		Address:     "10.10.10.1",
		Details:     details,
	}

	result, err := json.Marshal(neighbor)
	if err != nil {
		t.Error(err)
	}

	_ = result
}

func TestCommunityStringify(t *testing.T) {
	com := Community{23, 42}
	if com.String() != "23:42" {
		t.Error("Expected 23:42, got:", com.String())
	}

	extCom := ExtCommunity{"ro", 42, 123}
	if extCom.String() != "ro:42:123" {
		t.Error("Expected ro:42:123, but got:", extCom.String())
	}

	emptyCom := Community{}
	if emptyCom.String() != "" {
		t.Error("Unexpected result:", emptyCom.String())
	}

	emptyExtCom := ExtCommunity{}
	if emptyExtCom.String() != "" {
		t.Error("Unexpected result:", emptyExtCom.String())
	}
}

func TestHasCommunity(t *testing.T) {
	com := Community{23, 42}

	bgp := &BGPInfo{
		Communities: []Community{
			{42, 123},
			{23, 42},
			{42, 23},
		},
		ExtCommunities: []ExtCommunity{
			{"rt", "23", "42"},
			{"ro", "123", "456"},
		},
		LargeCommunities: []Community{
			{1000, 23, 42},
			{2000, 123, 456},
		},
	}

	if bgp.HasCommunity(com) == false {
		t.Error("Expected community 23:42 to be present")
	}

	if bgp.HasCommunity(Community{111, 11}) != false {
		t.Error("Expected community 111:11 to be not present")
	}

	if bgp.HasExtCommunity(ExtCommunity{"ro", "123", "456"}) == false {
		t.Error("Expected ro:123:456 in ext community set")
	}

	if bgp.HasExtCommunity(ExtCommunity{"ro", "111", "11"}) != false {
		t.Error("Expected ro:111:111 not in ext community set")
	}

	if bgp.HasLargeCommunity(Community{2000, 123, 456}) == false {
		t.Error("Expected community 2000:123:456 present")
	}

	if bgp.HasLargeCommunity(Community{23, 42}) != false {
		t.Error("23:42 should not be present in large communities")
	}
}

/*
func TestUniqueCommunities(t *testing.T) {
	all := Communities{Community{23, 42}, Community{42, 123}, Community{23, 42}}
	unique := all.Unique()
	if len(unique) != 2 {
		t.Error("len(unique) should be < len(all)")
	}
	t.Log("All:", all, "Unique:", unique)
}

func TestUniqueExtCommunities(t *testing.T) {
	all := ExtCommunities{
		ExtCommunity{"rt", "23", "42"},
		ExtCommunity{"ro", "42", "123"},
		ExtCommunity{"rt", "23", "42"}}
	unique := all.Unique()
	if len(unique) != 2 {
		t.Error("len(unique) should be < len(all)")
	}
	t.Log("All:", all, "Unique:", unique)
}
*/
