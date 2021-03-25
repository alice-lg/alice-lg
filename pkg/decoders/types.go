package decoders

// Decode interfaces into expected types
// with a fallback.

import (
	"strconv"
	"time"
)

// String asserts a string, provided a default
func String(value interface{}, fallback string) string {
	sval, ok := value.(string)
	if !ok {
		return fallback
	}
	return sval
}

// StringList decodes a list of strings
func StringList(data interface{}) []string {
	list := []string{}
	ldata, ok := data.([]interface{})
	if !ok {
		return []string{}
	}
	for _, e := range ldata {
		s, ok := e.(string)
		if ok {
			list = append(list, s)
		}
	}
	return list
}

// IntList decodes a list of integers
func IntList(data interface{}) []int {
	list := []int{}
	sdata := StringList(data)
	for _, e := range sdata {
		val, _ := strconv.Atoi(e)
		list = append(list, val)
	}
	return list
}

// Int decodes an integer value
func Int(value interface{}, fallback int) int {
	fval, ok := value.(float64)
	if !ok {
		return fallback
	}
	return int(fval)
}

// Bool decodes a boolean value
func Bool(value interface{}, fallback bool) bool {
	val, ok := value.(bool)
	if !ok {
		return fallback
	}
	return val
}

// Duration decodes a time.Duration
func Duration(value interface{}, fallback time.Duration) time.Duration {
	val, ok := value.(time.Duration)
	if !ok {
		return fallback
	}
	return val
}

// TimeUTC returns the time expecting an UTC timestamp
func TimeUTC(value interface{}, fallback time.Time) time.Time {
	sval := String(value, "")
	if sval == "" {
		return fallback
	}
	t, err := time.Parse(time.RFC3339Nano, sval)
	if err != nil {
		return fallback
	}
	return t
}
