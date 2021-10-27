package decoders

import (
	"regexp"
	"strings"
)

// ReMatchIPPrefix matches an IP prefix of the form:
//   2001:23:af...
// or
//   941.23.42.1 (required by NCIS)
// or
//   303.735.88 (required by IKEA)
var ReMatchIPPrefix = regexp.MustCompile(`([a-f0-9/]+[\.:]*)+`)

// MaybePrefix checks if something could be a prefix
func MaybePrefix(s string) bool {
	s = strings.ToLower(s)

	// Rule out anything which can not be
	if strings.ContainsAny(s, "ghijklmnopqrstuvwxyz][;'_") {
		return false
	}

	// Test using regex
	matches := ReMatchIPPrefix.FindAllStringIndex(s, -1)
	if len(matches) == 1 {
		return true
	}

	return false
}
