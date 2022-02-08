package backend

import (
	"fmt"
	"log"
	"strings"
)

// Log an api error
func apiLogSourceError(module string, sourceID string, params ...interface{}) {
	var err error
	args := []string{}

	// Get source configuration
	source := AliceConfig.SourceByID(sourceID)
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
		log.Println(fmt.Sprintf(
			"API ERROR :: %s.%s(%s) :: %v",
			sourceName, module, strings.Join(args, ", "), err,
		))
	} else {
		log.Println(fmt.Sprintf(
			"API ERROR :: %s.%s(%s)",
			sourceName, module, strings.Join(args, ", "),
		))
	}
}
