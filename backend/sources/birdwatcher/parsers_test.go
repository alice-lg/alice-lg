package birdwatcher

import (
	"encoding/json"
	"testing"
)

const API_RESPONSE_NEIGHBOURS = `
{"api":{"Version":"1.7.11","result_from_cache":true,"cache_status":{"orig_ttl":0,"cached_at":{"date":"","timezone_type":"","timezone":""}}},"protocols":{"ID103_AS25074_194.9.117.1":{"action":"restart","bgp_state":"Established","bird_protocol":"BGP","connection":"Established","description":"AS25074 194.9.117.1 MESH GmbH","export_withdraws":"45756        ---        ---        ---      45654","hold_timer":"146/180","import_limit":16000,"input_filter":"(unnamed)","keepalive_timer":"8/60","neighbor_address":"194.9.117.1","neighbor_as":25074,"neighbor_caps":"refresh AS4","neighbor_id":"212.162.48.85","output_filter":"(unnamed)","preference":100,"protocol":"ID103_AS25074_194.9.117.1","route_change_stats":"received   rejected   filtered    ignored   accepted","route_changes":{"export_updates":{"accepted":200340,"ignored":10,"received":202884,"rejected":117},"import_updates":{"accepted":150,"filtered":388,"ignored":12965,"received":13503,"rejected":0},"import_withdraws":{"accepted":15,"filtered":4,"received":15,"rejected":0}},"route_limit":"139/16000","routes":{"exported":35707,"filtered":4,"imported":135,"preferred":114},"session":"external route-server AS4","source_address":"194.9.117.253","state":"up","state_changed":"2017-05-17 03:20:28","table":"master"},"ID109_AS31078_194.9.117.4":{"action":"restart","bgp_state":"Established","bird_protocol":"BGP","connection":"Established","description":"AS31078 194.9.117.4 Netsign GmbH","export_withdraws":"115690        ---        ---        ---     115445","hold_timer":"146/180","import_limit":16000,"input_filter":"(unnamed)","keepalive_timer":"16/60","neighbor_address":"194.9.117.4","neighbor_as":31078,"neighbor_caps":"refresh","neighbor_id":"217.115.0.29","output_filter":"(unnamed)","preference":100,"protocol":"ID109_AS31078_194.9.117.4","route_change_stats":"received   rejected   filtered    ignored   accepted","route_changes":{"export_updates":{"accepted":442671,"ignored":10,"received":448284,"rejected":14},"import_updates":{"accepted":14,"filtered":0,"ignored":1800,"received":1814,"rejected":0},"import_withdraws":{"accepted":6,"filtered":0,"received":6,"rejected":0}},"route_limit":"8/16000","routes":{"exported":35806,"filtered":0,"imported":8,"preferred":8},"session":"external route-server","source_address":"194.9.117.253","state":"up","state_changed":"2017-05-10 14:47:27","table":"master"}}, "ttl":"2017-05-22T08:34:04.008634978Z"}`

// Load test response json
func parseTestResponse(payload string) ClientResponse {
	result := make(ClientResponse)
	_ = json.Unmarshal([]byte(payload), &result)
	return result
}

func Test_ParseApiStatus(t *testing.T) {
	bird := parseTestResponse(API_RESPONSE_NEIGHBOURS)

	// mock config
	config := Config{Timezone: "UTC"} // Or ""

	apiStatus, err := parseApiStatus(bird, config)
	if err != nil {
		t.Error(err)
	}

	// Assertations
	if apiStatus.Version != "1.7.11" {
		t.Error("Expected version: 1.7.11, got:", apiStatus.Version)
	}

	if apiStatus.ResultFromCache == false {
		t.Error("Expected result_from_cache to be true")
	}

	// TODO: Test cache status and TTL parsing

}

func Test_NeighboursParsing(t *testing.T) {

}
