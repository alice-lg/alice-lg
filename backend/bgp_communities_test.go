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

func TestCommunityLookup(t *testing.T) {

	c := NgMakeWellKnownBgpCommunities()

	label, err := c.Lookup("65535:666")
	if err != nil {
		t.Error(err)
	}
	if label != "blackhole" {
		t.Error("Label should have been: blackhole, got:", label)
	}

	// Okay now try some fails
	label, err = c.Lookup("65535")
	if err == nil {
		t.Error("Expected error!")
	}

	label, err = c.Lookup("65535:23:42")
	if err == nil {
		t.Error("Expected not found error!")
	}

}
