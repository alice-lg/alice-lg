package http

import "testing"

func TestValidateNeighborID(t *testing.T) {
	for _, id := range []string{"AS64402_1", "R194_42", "198.51.100.10", "2001:db8::10"} {
		if _, err := validateNeighborID(id); err != nil {
			t.Errorf("validateNeighborID(%q) = %v; want accepted", id, err)
		}
	}
	for _, id := range []string{"", "AS1?x=1", "AS1#f", "AS1&a=b", "../status", "AS1/x", "AS1%3F", "AS 1", "fe80::1%eth0"} {
		if _, err := validateNeighborID(id); err == nil {
			t.Errorf("validateNeighborID(%q) accepted; want rejected", id)
		}
	}
}
