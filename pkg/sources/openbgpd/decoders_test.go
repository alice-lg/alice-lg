package openbgpd

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func readTestData(filename string) map[string]interface{} {
	data, _ := ioutil.ReadFile(filepath.Join("testdata", filename))
	payload := make(map[string]interface{})
	_ = json.Unmarshal(data, &payload)
	return payload
}

func TestDecodeAPIStatus(t *testing.T) {
	res := readTestData("status.json")
	s := decodeAPIStatus(res)
	t.Log(s.ServerTime)
	t.Log(s.LastReboot)
}

func TestDecodeNeighbors(t *testing.T) {
	res := readTestData("show.neighbor.json")
	n, err := decodeNeighbors(res)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n[0])
}

func TestDecodeNeighborsStatus(t *testing.T) {
	res := readTestData("show.summary.json")
	n, err := decodeNeighborsStatus(res)
	if err != nil {
		t.Fatal(err)
	}
	if len(n) != 3 {
		t.Error("unexpected length:", len(n))
	}
	t.Log(*n[0])
}

func TestDecodeRoutes(t *testing.T) {
	res := readTestData("rib.json")
	routes, err := decodeRoutes(res)
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 2 {
		t.Error("unexpected length:", len(routes))
	}

	// Check first route
	r := routes[0]
	if r.Network != "23.42.1.0/24" {
		t.Error("unexpected network:", r.Network)
	}
	// Community decoding
	if r.BGP.Communities[0][0] != 20119 {
		t.Error("unexpected community:", r.BGP.Communities[0])
	}
	if r.BGP.Communities[0][1] != 3 {
		t.Error("unexpected community:", r.BGP.Communities[0])
	}
	if r.BGP.ExtCommunities[1][0] != "rt" {
		t.Error("unexpected community:", r.BGP.ExtCommunities[0])
	}
	if r.BGP.ExtCommunities[1][1] != 65000 {
		t.Error("unexpected community:", r.BGP.ExtCommunities[0])
	}
	if r.BGP.ExtCommunities[1][2] != 11000 {
		t.Error("unexpected community:", r.BGP.ExtCommunities[0])
	}

	if r.BGP.AsPath[0] != 1111 {
		t.Error("unexpected as_path:", r.BGP.AsPath)
	}
	t.Log(r.Age)
}
