package main

import (
	"testing"
)

// Test configuration loading and parsing
// using the default config

func TestLoadConfigs(t *testing.T) {

	config, err := loadConfigs("../etc/alicelg/alice.conf", "", "")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	t.Log(config)
}
