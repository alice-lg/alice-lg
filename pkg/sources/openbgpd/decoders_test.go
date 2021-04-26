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
