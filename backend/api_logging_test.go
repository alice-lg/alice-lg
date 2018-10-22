package main

import (
	"fmt"
	"testing"
)

func TestApiLogSourceError(t *testing.T) {
	err := fmt.Errorf("an unexpected error occured")

	conf := &Config{
		Sources: []*SourceConfig{
			&SourceConfig{
				Id:   0,
				Name: "rs1.example.net (IPv4)",
			},
		},
	}

	AliceConfig = conf

	apiLogSourceError("foo.bar", 0, 23, "Test")
	apiLogSourceError("foo.bam", 0, err)
	apiLogSourceError("foo.baz", 0, 23, 42, "foo", err)
}
