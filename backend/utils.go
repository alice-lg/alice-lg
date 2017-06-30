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
 Check array membership
*/
func MemberOf(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

/*
 Check if something could be a prefix
*/
func MaybePrefix(s string) bool {
	s = strings.ToLower(s)

	// Rule out anything which can not be
	if strings.ContainsAny(s, "ghijklmnopqrstuvwxyz][;'_") {
		return false
	}

	// Test using regex
	matches := regexp.MustCompile(`([a-f0-9/]+[\.:]?)+`).FindAllStringIndex(s, -1)
	if len(matches) == 1 {
		return true
	}

	return false
}
