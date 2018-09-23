package main

import (
	"testing"
)

// Test configuration loading and parsing
// using the default config

func TestLoadConfigs(t *testing.T) {

	config, err := loadConfig("../etc/alicelg/alice.example.conf")
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

	// Check communities
	label, err := config.Ui.BgpCommunities.Lookup("1:23")
	if err != nil {
		t.Error(err)
	}
	if label != "some tag" {
		t.Error("expcted to find example community 1:23 with 'some tag'",
			"but got:", label)
	}
}

func TestSourceConfigDefaultsOverride(t *testing.T) {

	config, err := loadConfig("../etc/alicelg/alice.example.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	// Get sources
	rs1 := config.Sources[0]
	rs2 := config.Sources[1]

	// Source 1 should be on default time
	// Source 2 should have an override
	// For now it should be sufficient to test if
	// the serverTime(rs1) != serverTime(rs2)
	if rs1.Birdwatcher.ServerTime == rs2.Birdwatcher.ServerTime {
		t.Error("Server times should be different between",
			"source 1 and 2 in example configuration",
			"(alice.example.conf)")
	}

	// Check presence of timezone, default: UTC (rs1)
	// override: Europe/Bruessels (rs2)
	if rs1.Birdwatcher.Timezone != "UTC" {
		t.Error("Expected RS1 Timezone to be default: UTC")
	}

	if rs2.Birdwatcher.Timezone != "Europe/Brussels" {
		t.Error("Expected 'Europe/Brussels', got", rs2.Birdwatcher.Timezone)
	}
}

func TestBlackholeParsing(t *testing.T) {
	config, err := loadConfig("../etc/alicelg/alice.example.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	// Get first source
	rs1 := config.Sources[0]

	if len(rs1.Blackholes) != 2 {
		t.Error("Rs1 should have configured 2 blackholes. Got:", rs1.Blackholes)
		return
	}

	if rs1.Blackholes[0] != "10.23.6.666" {
		t.Error("Unexpected blackhole, got:", rs1.Blackholes[0])
	}
}

func TestSourceASN(t *testing.T) {
	config, err := loadConfig("../etc/alicelg/alice.example.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	// Get sources
	rs1 := config.Sources[0]
	if rs1.Asn != 99999 {
		t.Error("Expected RS1 ASN to be: 99999, not:", rs1.Asn)
	}

	rs2 := config.Sources[1]
	if rs2.Asn != 9033 {
		t.Error("Expected RS2 to fall back to AS9033, not:", rs2.Asn)
	}
}

func TestRpkiConfig(t *testing.T) {
	config, err := loadConfig("../etc/alicelg/alice.example.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	if len(config.Ui.Rpki.Valid) != 3 {
		t.Error("Unexpected RPKI:VALID,", config.Ui.Rpki.Valid)
	}
	if len(config.Ui.Rpki.Invalid) != 4 {
		t.Error("Unexpected RPKI:INVALID,", config.Ui.Rpki.Invalid)
		return // We would fail hard later
	}

	// Check fallback
	if config.Ui.Rpki.NotChecked[0] != "9033" {
		t.Error(
			"Expected NotChecked to fall back to defaults, got:",
			config.Ui.Rpki.NotChecked,
		)
	}

	// Check range postprocessing
	if config.Ui.Rpki.Invalid[3] != "*" {
		t.Error("Missing wildcard from config")
	}

	t.Log(config.Ui.Rpki)
}
