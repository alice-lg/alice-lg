package store

// Some helper functions
import (
	"strconv"
	"strings"
)

// ContainsCi is like `strings.Contains` but case insensitive
func ContainsCi(s, substr string) bool {
	return strings.Contains(
		strings.ToLower(s),
		strings.ToLower(substr),
	)
}

// SerializeReasons asserts the bgp community parts are
// actually strings, because there are no such things as
// integers as keys in json.
// Serialization of this is undefined behaviour, so we
// keep these internally but provide a string as a key for
// serialization
func SerializeReasons(reasons map[int]string) map[string]string {
	res := make(map[string]string)
	for id, reason := range reasons {
		res[strconv.Itoa(id)] = reason
	}
	return res
}
