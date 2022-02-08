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
	json.Unmarshal(data, &payload)
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
	if r.Bgp.Communities[0][0] != 20119 {
		t.Error("unexpected community:", r.Bgp.Communities[0])
	}
	if r.Bgp.Communities[0][1] != 3 {
		t.Error("unexpected community:", r.Bgp.Communities[0])
	}
	if r.Bgp.ExtCommunities[1][0] != "rt" {
		t.Error("unexpected community:", r.Bgp.ExtCommunities[0])
	}
	if r.Bgp.ExtCommunities[1][1] != 65000 {
		t.Error("unexpected community:", r.Bgp.ExtCommunities[0])
	}
	if r.Bgp.ExtCommunities[1][2] != 11000 {
		t.Error("unexpected community:", r.Bgp.ExtCommunities[0])
	}

	if r.Bgp.AsPath[0] != 1111 {
		t.Error("unexpected as_path:", r.Bgp.AsPath)
	}
	t.Log(r.Age)
}

func TestDecodeExtendedCommunities(t *testing.T) {
	data := []interface{}{"rt 123:456", "error invalid community"}
	comms := decodeExtendedCommunities(data)
	if len(comms) != 1 {
		t.Fatal("expected 1 valid community")
	}
	if comms[0][0] != "rt" && comms[0][1] != 123 && comms[0][2] != 456 {
		t.Fatal("unexpected result:", comms[0])
	}
}
