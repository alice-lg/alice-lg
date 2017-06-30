package main

// Some helper functions
import (
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
