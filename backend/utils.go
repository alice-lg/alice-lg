package main

// Some helper functions
import (
	"regexp"
	"strings"
)

/*
 Case Insensitive Contains
*/
func ContainsCi(s, substr string) bool {
	return strings.Contains(
		strings.ToLower(s),
		strings.ToLower(substr),
	)
}

/*
 Check if something could be a prefix
*/
func MaybePrefix(s string) bool {
	s = strings.ToLower(s)

	// Test using regex
	matches := regexp.MustCompile(`([a-f0-9/]+[\.:]?)+`).FindAllStringIndex(s, -1)
	if len(matches) == 1 {
		return true
	}

	return false
}
