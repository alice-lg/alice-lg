package main

import (
	"fmt"
	"log"
	"strings"
)

// Log an api error
func apiLogError(module string, params ...interface{}) {
	var err error
	args := []string{}

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
			"API :: %s(%s) :: ERROR: %v",
			module, strings.Join(args, ", "), err,
		))
	} else {
		log.Println(fmt.Sprintf(
			"API :: %s(%s)",
			module, strings.Join(args, ", "),
		))
	}
}
