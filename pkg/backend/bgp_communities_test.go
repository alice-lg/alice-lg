package backend

import (
	"testing"
)

func TestCommunityLookup(t *testing.T) {

	c := MakeWellKnownBgpCommunities()

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

func TestSetCommunity(t *testing.T) {
	c := MakeWellKnownBgpCommunities()

	c.Set("2342:10", "foo")
	c.Set("2342:42:23", "bar")

	// Simple lookup
	label, err := c.Lookup("2342:10")
	if err != nil {
		t.Error(err)
	}
	if label != "foo" {
		t.Error("Expected foo for 2342:10, got:", label)
	}

	label, err = c.Lookup("2342:42:23")
	if err != nil {
		t.Error(err)
	}
	if label != "bar" {
		t.Error("Expected bar for 2342:42:23, got:", label)
	}
}

func TestWildcardLookup(t *testing.T) {
	c := MakeWellKnownBgpCommunities()

	c.Set("2342:*", "foobar $0")
	c.Set("42:*:1", "baz")

	// This should work
	label, err := c.Lookup("2342:23")
	if err != nil {
		t.Error(err)
	}
	if label != "foobar $0" {
		t.Error("Did not get expected label.")
	}

	// This however not
	label, err = c.Lookup("2342:23:666")
	if err == nil {
		t.Error("Lookup should have failed, got label:", label)
	}

	// This should again work
	label, err = c.Lookup("42:123:1")
	if err != nil {
		t.Error(err)
	}
	if label != "baz" {
		t.Error("Unexpected label for key")
	}
}

func TestAPICommunities(t *testing.T) {
	c := MakeWellKnownBgpCommunities()
	comm := c.APICommunities()
	if len(comm) != 14 {
		t.Error("unexpected len(communities) = ", len(comm))
	}
}
