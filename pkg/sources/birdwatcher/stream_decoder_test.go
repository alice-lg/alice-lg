package birdwatcher

import (
	"os"
	"testing"
)

func TestParseRoutesResponseStream(t *testing.T) {
	f, err := os.Open("../../../../routes_dump_v4.json")
	if err != nil {
		return
	}
	defer f.Close()

	cfg := Config{
		Timezone:   "Europe/Berlin",
		ServerTime: "2006-01-02T15:04:05.999999999Z07:00",
	}

	meta, routes, err := parseRoutesResponseStream(f, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(meta)
	t.Log(len(routes))

}
