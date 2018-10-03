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

func TestRejectAndNoexportReasons(t *testing.T) {
	config, err := loadConfig("../etc/alicelg/alice.example.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	// Rejection reasons
	description, err := config.Ui.RoutesRejections.Reasons.Lookup("23:42:1")
	if err != nil {
		t.Error(err)
	}

	if description != "Some made up reason" {
		t.Error("Unexpected reason for 23:42:1 -", description)
	}

	// Noexport reasons
	description, err = config.Ui.RoutesNoexports.Reasons.Lookup("23:46:1")
	if err != nil {
		t.Error(err)
	}

	if description != "Some other made up reason" {
		t.Error("Unexpected reason for 23:46:1 -", description)
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

func TestOwnASN(t *testing.T) {
	config, err := loadConfig("../etc/alicelg/alice.example.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
	}

	if config.Server.Asn != 9033 {
		t.Error("Expected a set server asn")
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

func TestRejectCandidatesConfig(t *testing.T) {
	config, err := loadConfig("../etc/alicelg/alice.example.conf")
	if err != nil {
		t.Error("Could not load test config:", err)
		return
	}

	t.Log(config.Ui.RoutesRejectCandidates.Communities)

	description, err := config.Ui.RoutesRejectCandidates.Communities.Lookup("23:42:46")
	if err != nil {
		t.Error(err)
	}

	if description != "reject-candidate-3" {
		t.Error("expected 23:42:46 to be a 'reject-candidate'")
	}
}
