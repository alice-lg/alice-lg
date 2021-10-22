package http

import (
	"fmt"
	"testing"
)

func TestApiLogSourceError(t *testing.T) {
	err := fmt.Errorf("an unexpected error occured")

	conf := &Config{
		Sources: []*SourceConfig{
			&SourceConfig{
				ID:   "rs1v4",
				Name: "rs1.example.net (IPv4)",
			},
		},
	}

	AliceConfig = conf

	apiLogSourceError("foo.bar", "rs1v4", 23, "Test")
	apiLogSourceError("foo.bam", "rs1v4", err)
	apiLogSourceError("foo.baz", "rs1v4", 23, 42, "foo", err)
}
