package openbgpd

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/pools"
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
	ipPtr := pools.Networks4.Acquire("23.42.1.0/24")
	if r.Network != ipPtr {
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
	t.Log(r.BGP.ExtCommunities)

	if r.BGP.AsPath[0] != 1111 {
		t.Error("unexpected as_path:", r.BGP.AsPath)
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

func TestDecodeMalformedExtendedCommunities(t *testing.T) {
	data := []interface{}{
		"0x8000000000000000",
		"8000000000000000",
		"rt 1239", "generic :123", "generic ro-23:123",
		"generic 123123192399281398193489:asd",
		"[0] 0x8000000000000000",
		"[0] 0x800000000:0000000",
		"foo bar:23:42",
		"foo 2342:bar",
		"foo 23:bar:42",
		"foo",
		"b 9223372036854775808",
		922337203685477580,
		"ro  2::42",
		"generic rt a:b"}
	comms := decodeExtendedCommunities(data)
	t.Log(comms)
	if len(comms) > 0 {
		t.Error("expected empty communities")
	}
}
