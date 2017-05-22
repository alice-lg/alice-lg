package main

import (
	"testing"
)

// Test configuration loading and parsing
// using the default config

func TestLoadConfigs(t *testing.T) {

	config, err := loadConfig("../etc/alicelg/alice.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	if config.Server.Listen == "" {
		t.Error("Listen string not present.")
	}

	if len(config.Ui.RoutesColumns) == 0 {
		t.Error("Route columns settings missing")
	}

	if len(config.Ui.RoutesRejections.Reasons) == 0 {
		t.Error("Rejection reasons missing")
	}
}
