package openbgpd

import (
	"testing"
)

func TestConfigAPIURL(t *testing.T) {
	cfg := &Config{
		API: "http://a",
	}

	url := cfg.APIURL("/%d/bgpd", 42)
	if url != "http://a/42/bgpd" {
		t.Error("unexpected url:", url)
	}
}
