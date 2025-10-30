package http

import (
	"fmt"
	"log"
	"strings"
)

// Log an api error
func (s *Server) logSourceError(
	module string,
	sourceID string,
	params ...any,
) {
	var err error
	args := []string{}

	// Get source configuration
	source := s.cfg.SourceByID(sourceID)
	sourceName := "unknown"
	if source != nil {
		sourceName = source.Name
	}

	// Build args string and get error from params
	for _, p := range params {
		// We have our error
		if e, ok := p.(error); ok {
			err = e
			continue
		}

		args = append(args, fmt.Sprintf("%v", p))
	}

	if err != nil {
		log.Printf(
			"API ERROR :: %s.%s(%s) :: %v",
			sourceName, module, strings.Join(args, ", "), err,
		)
	} else {
		log.Printf(
			"API ERROR :: %s.%s(%s)",
			sourceName, module, strings.Join(args, ", "),
		)
	}
}
