package decoders

// Decode interfaces into expected types
// with a fallback.

import (
	"fmt"
	"strconv"
	"strings"
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
	ldata, ok := data.([]interface{})
	if !ok {
		return []string{}
	}
	list := make([]string, 0, len(ldata))
	for _, e := range ldata {
		s, ok := e.(string)
		if ok {
			list = append(list, s)
		}
	}
	return list
}

// TrimmedCSVStringList makes a trimmed list of CSV strings
// omitting empty values.
func TrimmedCSVStringList(s string) []string {
	tokens := strings.Split(s, ",")
	list := []string{}
	for _, t := range tokens {
		if t == "" {
			continue
		}
		list = append(list, strings.TrimSpace(t))
	}
	return list
}

// IntList decodes a list of integers
func IntList(data interface{}) []int {
	sdata := StringList(data)
	list := make([]int, 0, len(sdata))
	for _, e := range sdata {
		val, err := strconv.Atoi(e)
		if err == nil {
			list = append(list, val)
		}
	}
	return list
}

// IntListFromStrings decodes a list of strings
// into a list of integers.
func IntListFromStrings(strs []string) []int {
	list := make([]int, 0, len(strs))
	for _, s := range strs {
		v, err := strconv.Atoi(s)
		if err != nil {
			continue // skip this
		}
		list = append(list, v)
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

// IntFromString decodes an integer from a string
func IntFromString(s string, fallback int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return val
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

// DurationTimeframe decodes a duration: Bgpctl encodes
// this using fmt_timeframe, whiuch outputs a format similar
// to that being understood by time.ParseDuration - however
// the time unit "w" (weeks) is not supported.
// According to https://github.com/openbgpd-portable/openbgpd-openbsd/blob/master/src/usr.sbin/bgpctl/bgpctl.c#L586-L591
// we have to parse %02lluw%01ud%02uh, %01ud%02uh%02um and %02u:%02u:%02u.
// This yields three formats:
//
//	01w3d01h
//	1d02h03m
//	01:02:03
func DurationTimeframe(value interface{}, fallback time.Duration) time.Duration {
	var sec, min, hour, day uint
	var week uint64
	sval := String(value, "")
	if sval == "" {
		return fallback
	}

	n, _ := fmt.Sscanf(sval, "%02dw%01dd%02dh", &week, &day, &hour)
	if n == 3 {
		return time.Duration(week)*7*24*time.Hour +
			time.Duration(day)*24*time.Hour +
			time.Duration(hour)*time.Hour
	}

	n, _ = fmt.Sscanf(sval, "%01dd%02dh%02dm", &day, &hour, &min)
	if n == 3 {
		return time.Duration(day)*24*time.Hour +
			time.Duration(hour)*time.Hour +
			time.Duration(min)*time.Minute
	}

	n, _ = fmt.Sscanf(sval, "%02d:%02d:%02d", &hour, &min, &sec)
	if n == 3 {
		return time.Duration(hour)*time.Hour +
			time.Duration(min)*time.Minute +
			time.Duration(sec)*time.Second
	}

	return fallback
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
