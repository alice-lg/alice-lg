package main

import (
	"testing"
)

func TestContainsCi(t *testing.T) {

	if ContainsCi("foo bar", "BaR") != true {
		t.Error("An unexpected error occured.")
	}

}
