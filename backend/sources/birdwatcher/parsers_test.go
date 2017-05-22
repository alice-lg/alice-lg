package birdwatcher

import (
	"encoding/json"
	"testing"
)

const API_RESPONSE_NEIGHBOURS = `
{"api":{"Version":"1.7.11","result_from_cache":true,"cache_status":{"orig_ttl":0,"cached_at":{"date":"","timezone_type":"","timezone":""}}},"protocols":{"ID103_AS25074_194.9.117.1":{"action":"restart","bgp_state":"Established","bird_protocol":"BGP","connection":"Established","description":"AS25074 194.9.117.1 MESH GmbH","export_withdraws":"45756        ---        ---        ---      45654","hold_timer":"146/180","import_limit":16000,"input_filter":"(unnamed)","keepalive_timer":"8/60","neighbor_address":"194.9.117.1","neighbor_as":25074,"neighbor_caps":"refresh AS4","neighbor_id":"212.162.48.85","output_filter":"(unnamed)","preference":100,"protocol":"ID103_AS25074_194.9.117.1","route_change_stats":"received   rejected   filtered    ignored   accepted","route_changes":{"export_updates":{"accepted":200340,"ignored":10,"received":202884,"rejected":117},"import_updates":{"accepted":150,"filtered":388,"ignored":12965,"received":13503,"rejected":0},"import_withdraws":{"accepted":15,"filtered":4,"received":15,"rejected":0}},"route_limit":"139/16000","routes":{"exported":35707,"filtered":4,"imported":135,"preferred":114},"session":"external route-server AS4","source_address":"194.9.117.253","state":"up","state_changed":"2017-05-17 03:20:28","table":"master"},"ID109_AS31078_194.9.117.4":{"action":"restart","bgp_state":"Established","bird_protocol":"BGP","connection":"Established","description":"AS31078 194.9.117.4 Netsign GmbH","export_withdraws":"115690        ---        ---        ---     115445","hold_timer":"146/180","import_limit":16000,"input_filter":"(unnamed)","keepalive_timer":"16/60","neighbor_address":"194.9.117.4","neighbor_as":31078,"neighbor_caps":"refresh","neighbor_id":"217.115.0.29","output_filter":"(unnamed)","preference":100,"protocol":"ID109_AS31078_194.9.117.4","route_change_stats":"received   rejected   filtered    ignored   accepted","route_changes":{"export_updates":{"accepted":442671,"ignored":10,"received":448284,"rejected":14},"import_updates":{"accepted":14,"filtered":0,"ignored":1800,"received":1814,"rejected":0},"import_withdraws":{"accepted":6,"filtered":0,"received":6,"rejected":0}},"route_limit":"8/16000","routes":{"exported":35806,"filtered":0,"imported":8,"preferred":8},"session":"external route-server","source_address":"194.9.117.253","state":"up","state_changed":"2017-05-10 14:47:27","table":"master"}}, "ttl":"2017-05-22T08:34:04.008634978Z"}`

const API_RESPONSE_ROUTES = `
{"api":{"Version":"1.7.11","result_from_cache":false,"cache_status":{"orig_ttl":0,"cached_at":{"date":"","timezone_type":"","timezone":""}}},"routes":[{"age":"2017-05-19 08:12:44","bgp":{"aggregator":"62.69.151.1 AS201785","as_path":["31078","201785"],"communities":[[65000,65000],[31078,200],[31078,211],[65011,1],[9033,3051]],"local_pref":"100","next_hop":"194.9.117.4","origin":"IGP"},"from_protocol":"ID109_AS31078_194.9.117.4","gateway":"194.9.117.4","interface":"eno7","learnt_from":"","metric":100,"network":"193.200.230.0/24","primary":true,"type":["BGP","unicast","univ"]}], "ttl":"2017-05-22T10:22:39.732071843Z"}`

const API_RESPONSE_ROUTES_FILTERED = `
{"api":{"Version":"1.7.11","result_from_cache":true,"cache_status":{"orig_ttl":0,"cached_at":{"date":"","timezone_type":"","timezone":""}}},"routes":[{"age":"2017-05-17 03:20:31","bgp":{"as_path":["25074","15368"],"communities":[[25074,123],[25074,333],[25074,2070],[25074,20702],[65000,29208]],"large_communities":[[9033,65666,9]],"local_pref":"100","med":"1","next_hop":"194.9.117.1","origin":"IGP"},"from_protocol":"ID103_AS25074_194.9.117.1","gateway":"194.9.117.1","interface":"eno7","learnt_from":"","metric":100,"network":"192.111.47.0/24","primary":true,"type":["BGP","unicast","univ"]}], "ttl":"2017-05-22T10:22:39.732071843Z"}`

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
	config := Config{Timezone: "UTC"} // Or ""
	bird := parseTestResponse(API_RESPONSE_NEIGHBOURS)

	neighbours, err := parseNeighbours(bird, config)
	if err != nil {
		t.Error(err)
	}

	// We have 4 neighbours in our test response
	if len(neighbours) != 2 {
		t.Error("Number of neighbours should be 2, is:", len(neighbours))
	}

	// Test neighbour parsing
	neighbour := neighbours[0]
	if neighbour.Asn == 0 {
		t.Error("Expected ASN to be <> 0")
	}

	if neighbour.Address != "194.9.117.1" {
		t.Error("Expected neighbour address to be: 194.9.117.1, not:", neighbour.Address)
	}

	if neighbour.Description == "" {
		t.Error("Expected description to be set")
	}

}

func Test_RoutesParsing(t *testing.T) {
	config := Config{Timezone: "UTC"} // Or ""
	bird := parseTestResponse(API_RESPONSE_ROUTES)

	routes, err := parseRoutes(bird, config)
	if err != nil {
		t.Error(err)
	}

	if len(routes) != 1 {
		t.Error("Expected parsed routes to be 1, not:", len(routes))
	}

	// TODO: addo more tests
}
