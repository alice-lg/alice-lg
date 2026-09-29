package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/sources/birdwatcher"
)

// writeConfig writes a test configuration file and returns its path.
func writeConfig(t *testing.T, data string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "alice.conf")
	if err := os.WriteFile(file, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	return file
}

func birdwatcherSource(id string) *SourceConfig {
	return &SourceConfig{
		ID:      id,
		Name:    id,
		Backend: SourceBackendBirdwatcher,
		Birdwatcher: birdwatcher.Config{
			ID:   id,
			Name: id,
			API:  "http://192.0.2.1:29184/",
		},
	}
}

// TestLoadConfigReturnsErrorForUnknownBirdwatcherType proves a bad
// source type is returned, so a reload does not exit the process.
func TestLoadConfigReturnsErrorForUnknownBirdwatcherType(t *testing.T) {
	file := writeConfig(t, `
[source.rs1-ipv4]
name = rs1 (IPv4)
[source.rs1-ipv4.birdwatcher]
api = http://192.0.2.1:29184/
type = bogus_table
`)
	_, err := LoadConfig(file)
	if err == nil || !strings.Contains(err.Error(), "unknown birdwatcher type") {
		t.Fatalf("expected an unknown birdwatcher type error, got %v", err)
	}
}

// TestLoadConfigReturnsErrorForUnexpandableExpression proves an
// expansion failure is reported to the caller instead of exiting.
func TestLoadConfigReturnsErrorForUnexpandableExpression(t *testing.T) {
	file := writeConfig(t, `
[source.rs1-ipv4]
name = {UNDEFINED} (IPv4)
`)
	_, err := LoadConfig(file)
	if err == nil || !strings.Contains(err.Error(), "error expanding expression") {
		t.Fatalf("expected an expansion error, got %v", err)
	}
}

// TestSourceConfigEqualIgnoresDerivedFields proves a moved source,
// or a new stream parser throttle, is not marked as changed.
func TestSourceConfigEqualIgnoresDerivedFields(t *testing.T) {
	held := birdwatcherSource("rs1")
	moved := birdwatcherSource("rs1")
	moved.Order = held.Order + 4
	moved.Birdwatcher.StreamParserThrottle = 10
	if !held.Equal(moved) {
		t.Error("a moved source should compare equal")
	}

	edited := birdwatcherSource("rs1")
	edited.Birdwatcher.API = "http://192.0.2.2:29184/"
	if held.Equal(edited) {
		t.Error("an edited API URL should compare different")
	}
}

// TestSourceConfigWithOrderSharesInstance proves a moved source
// keeps its client.
func TestSourceConfigWithOrderSharesInstance(t *testing.T) {
	held := birdwatcherSource("rs1")
	instance := held.GetInstance()

	moved := held.WithOrder(7)
	if moved.Order != 7 || held.Order != 0 {
		t.Errorf("expected order 7 on the copy and 0 on the original, got %d and %d",
			moved.Order, held.Order)
	}
	if moved.GetInstance() != instance {
		t.Error("the moved copy should share the memoized instance")
	}
}

// TestSourceConfigEqualDoesNotRaceGetInstance runs Equal while the
// instance is memoized; under -race it fails on an unguarded read.
func TestSourceConfigEqualDoesNotRaceGetInstance(t *testing.T) {
	for i := 0; i < 50; i++ {
		held := birdwatcherSource("rs1")
		desired := birdwatcherSource("rs1")
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			held.GetInstance()
		}()
		go func() {
			defer wg.Done()
			held.Equal(desired)
		}()
		wg.Wait()
	}
}
