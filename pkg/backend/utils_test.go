package backend

import (
	"testing"
)

func TestContainsCi(t *testing.T) {
	if ContainsCi("foo bar", "BaR") != true {
		t.Error("An unexpected error occured.")
	}
}

func TestMaybePrefix(t *testing.T) {
	expected := []struct {
		string
		bool
	}{
		{"10.0.0", true},
		{"23.42.11.42/23", true},
		{"fa42:2342::/32", true},
		{"200", true},
		{"2001:", true},
		{"A", true},
		{"A b", false},
		{"23 Foo", false},
		{"Nordfoo", false},
		{"122.beef:", true}, // sloppy
		{"122.beef:", true}, // very
		{"122:beef", true},  // sloppy.
	}

	for _, e := range expected {
		if MaybePrefix(e.string) != e.bool {
			t.Error("Expected", e.string, "to be prefix:", e.bool)
		}
	}
}

func TestTrimmedStringList(t *testing.T) {
	l := TrimmedStringList("foo, bar   , dreiundzwanzig,")

	if len(l) != 3 {
		t.Error("Expected length to be 3, got:", len(l))
	}

	if l[0] != "foo" || l[1] != "bar" || l[2] != "dreiundzwanzig" {
		t.Error("Expected list of [foo, bar, dreiundzwanzig], got:", l)
	}
}
