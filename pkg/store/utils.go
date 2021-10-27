package store

// Some helper functions
import (
	"strconv"
	"strings"
	"time"
)

// ContainsCi is like `strings.Contains` but case insensitive
func ContainsCi(s, substr string) bool {
	return strings.Contains(
		strings.ToLower(s),
		strings.ToLower(substr),
	)
}

// MemberOf checks if a key is present in
// a list of strings.
func MemberOf(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

// SerializeReasons asserts the bgp communitiy parts are
// actually strings, because there are no such things as
// integers as keys in json.
// Serialization of this is undefined behaviour, so we
// keep these interallybut provide a string as a key for
// serialization
func SerializeReasons(reasons map[int]string) map[string]string {
	res := make(map[string]string)
	for id, reason := range reasons {
		res[strconv.Itoa(id)] = reason
	}
	return res
}

// DurationMs converts time.Duration to milliseconds
func DurationMs(d time.Duration) float64 {
	return float64(d) / 1000.0 / 1000.0 // nano -> micro -> milli
}
