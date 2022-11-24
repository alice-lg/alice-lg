package config

import (
	"testing"

	"github.com/alice-lg/alice-lg/pkg/sources/birdwatcher"
	"github.com/alice-lg/alice-lg/pkg/sources/gobgp"
)

// Test configuration loading and parsing
// using the default config

func TestLoadConfigs(t *testing.T) {

	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	if config.Server.Listen == "" {
		t.Error("Listen string not present.")
	}

	if len(config.UI.RoutesColumns) == 0 {
		t.Error("Route columns settings missing")
	}

	if len(config.UI.RoutesRejections.Reasons) == 0 {
		t.Error("Rejection reasons missing")
	}

	// Check communities
	label, err := config.UI.BGPCommunities.Lookup("1:23")
	if err != nil {
		t.Error(err)
	}
	if label != "some tag" {
		t.Error("expcted to find example community 1:23 with 'some tag'",
			"but got:", label)
	}
}

// TestSourceConfig checks that the proper backend type was identified for each
// example routeserver
func TestSourceConfig(t *testing.T) {
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	// Get sources
	rs1 := config.Sources[0] // Birdwatcher v4
	rs2 := config.Sources[1] // Birdwatcher v6
	rs3 := config.Sources[2] // GoBGP

	nilBirdwatcherConfig := birdwatcher.Config{}
	if rs1.Birdwatcher == nilBirdwatcherConfig {
		t.Errorf(
			"Example routeserver %s should have been identified as a birdwatcher source but was not",
			rs1.Name,
		)
	}
	if rs2.Birdwatcher == nilBirdwatcherConfig {
		t.Errorf(
			"Example routeserver %s should have been identified as a birdwatcher source but was not",
			rs2.Name,
		)
	} else {
		if rs2.Birdwatcher.AltPipeProtocolSuffix != "_lg" {
			t.Error("unexpected alt_pipe_suffix:", rs2.Birdwatcher.AltPipeProtocolSuffix)
		}
	}
	nilGoBGPConfig := gobgp.Config{}
	if rs3.GoBGP == nilGoBGPConfig {
		t.Errorf(
			"Example routeserver %s should have been identified as a gobgp source but was not",
			rs3.Name,
		)
	}
}

func TestSourceConfigDefaultsOverride(t *testing.T) {
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	// Get sources
	rs1 := config.Sources[0] // Birdwatcher v4
	rs2 := config.Sources[1] // Birdwatcher v6
	rs3 := config.Sources[2] // GoBGP

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

	if rs3.GoBGP.ProcessingTimeout != 300 {
		t.Error(
			"Expected GoBGP example to set 300s 'processing_timeout', got",
			rs3.GoBGP.ProcessingTimeout,
		)
	}
}

func TestRejectAndNoexportReasons(t *testing.T) {
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	// Rejection reasons
	description, err := config.UI.RoutesRejections.Reasons.Lookup("23:42:1")
	if err != nil {
		t.Error(err)
	}

	if description != "Some made up reason" {
		t.Error("Unexpected reason for 23:42:1 -", description)
	}

	// Noexport reasons
	description, err = config.UI.RoutesNoexports.Reasons.Lookup("23:46:1")
	if err != nil {
		t.Error(err)
	}

	if description != "Some other made up reason" {
		t.Error("Unexpected reason for 23:46:1 -", description)
	}
}

func TestBlackholeParsing(t *testing.T) {
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
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
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	if config.Server.Asn != 9033 {
		t.Error("Expected a set server asn")
	}
}

func TestRpkiConfig(t *testing.T) {
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	if len(config.UI.Rpki.Valid) != 3 {
		t.Error("Unexpected RPKI:VALID,", config.UI.Rpki.Valid)
	}
	if len(config.UI.Rpki.Invalid) != 4 {
		t.Fatal("Unexpected RPKI:INVALID,", config.UI.Rpki.Invalid)
	}

	// Check fallback
	if config.UI.Rpki.NotChecked[0] != "9033" {
		t.Error(
			"Expected NotChecked to fall back to defaults, got:",
			config.UI.Rpki.NotChecked,
		)
	}

	// Check range postprocessing
	if config.UI.Rpki.Invalid[3] != "*" {
		t.Error("Missing wildcard from config")
	}

	t.Log(config.UI.Rpki)
}

func TestRejectCandidatesConfig(t *testing.T) {
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	t.Log(config.UI.RoutesRejectCandidates.Communities)

	description, err := config.UI.RoutesRejectCandidates.Communities.Lookup("23:42:46")
	if err != nil {
		t.Error(err)
	}

	if description != "reject-candidate-3" {
		t.Error("expected 23:42:46 to be a 'reject-candidate'")
	}
}

// TestDefaultHTTPTimeout checks that the default HTTP timeout
// be set when not configured from a config file
func TestDefaultHTTPTimeout(t *testing.T) {
	config, err := LoadConfig("testdata/alice.conf")
	if err != nil {
		t.Fatal("Could not load test config:", err)
	}

	if config.Server.HTTPTimeout != DefaultHTTPTimeout {
		t.Error("Expected HTTP timeout be set to", DefaultHTTPTimeout,
			"but got", config.Server.HTTPTimeout)
	}
}

func TestPostgresStoreConfig(t *testing.T) {
	config, _ := LoadConfig("testdata/alice.conf")
	if config.Server.StoreBackend != "postgres" {
		t.Error("unexpected StoreBackend:", config.Server.StoreBackend)
	}
	if config.Postgres.MinConns != 10 {
		t.Error("unexpected MinConns:", config.Postgres.MinConns)
	}
	if config.Postgres.MaxConns != 128 {
		t.Error("unexpected MaxConns:", config.Postgres.MaxConns)
	}
	t.Log(config.Postgres)
}

func TestGetBlackholeCommunities(t *testing.T) {
	config, _ := LoadConfig("testdata/alice.conf")
	comms := config.UI.BGPBlackholeCommunities

	if comms.Communities[0][0].([]int)[0] != 1337 {
		t.Error("unexpected community:", comms.Communities[0])
	}
	if len(comms.ExtCommunities) != 1 {
		t.Error("unexpected communities:", comms.ExtCommunities)
	}
	if len(comms.LargeCommunities) != 1 {
		t.Error("unexpected communities:", comms.LargeCommunities)
	}
	t.Log(comms)
}
