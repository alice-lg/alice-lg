package main

import (
	"fmt"
	"testing"
)

func TestApiLogError(t *testing.T) {
	err := fmt.Errorf("an unexpected error occured")

	apiLogError("foo.bar", 23, "Test")
	apiLogError("foo.bam", err)
	apiLogError("foo.baz", 23, 42, "foo", err)
}
