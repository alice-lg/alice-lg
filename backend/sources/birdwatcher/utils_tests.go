package birdwatcher

import (
	"testing"
)

func TestIsProtocolUp(t *testing.T) {
	tests := map[string]bool{
		"up":   true,
		"uP":   true,
		"Up":   true,
		"UP":   true,
		"down": false,
	}

	for up, expected := range tests {
		if isProtocolUp(up) != expected {
			t.Error("f(", up, ") != ", expected)
		}
	}
}
