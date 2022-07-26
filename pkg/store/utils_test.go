package store

import (
	"testing"
)

func TestContainsCi(t *testing.T) {
	if ContainsCi("foo bar", "BaR") != true {
		t.Error("An unexpected error occured.")
	}
	if ContainsCi("Luxembourg Online SA", "Goo") == true {
		t.Error("Should ne no match")
	}
}
