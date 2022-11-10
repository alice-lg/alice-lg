package http

import (
	"fmt"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/config"
)

func TestApiLogSourceError(t *testing.T) {
	err := fmt.Errorf("an unexpected error occurred")

	cfg := &config.Config{
		Sources: []*config.SourceConfig{
			{
				ID:   "rs1v4",
				Name: "rs1.example.net (IPv4)",
			},
		},
	}

	s := &Server{cfg: cfg}

	s.logSourceError("foo.bar", "rs1v4", 23, "Test")
	s.logSourceError("foo.bam", "rs1v4", err)
	s.logSourceError("foo.baz", "rs1v4", 23, 42, "foo", err)
}
