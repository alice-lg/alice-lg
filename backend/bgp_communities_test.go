package main

import (
	"testing"
)

func TestMergeCommunities(t *testing.T) {

	c := MakeWellKnownBgpCommunities()

	merged := c.Merge(BgpCommunities{
		"2342:0":   "foo",
		"2342:123": "bar",
	})

	if merged["65535:666"] != "blackhole" {
		t.Error("old values should be present")
	}

	if merged["2342:123"] != "bar" {
		t.Error("new values should be present")
	}
}
