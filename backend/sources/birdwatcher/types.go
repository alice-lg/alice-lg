package birdwatcher

import (
	"strconv"
)

/*
 * Types helper for parser
 */

// Assert string, provide default
func mustString(value interface{}, fallback string) string {
	sval, ok := value.(string)
	if !ok {
		return fallback
	}
	return sval
}

// Assert list of strings
func mustStringList(data interface{}) []string {
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

// Convert list of strings to int
func mustIntList(data interface{}) []int {
	list := []int{}
	sdata := mustStringList(data)
	for _, e := range sdata {
		val, _ := strconv.Atoi(e)
		list = append(list, val)
	}
	return list
}

func mustInt(value interface{}, fallback int) int {
	fval, ok := value.(float64)
	if !ok {
		return fallback
	}
	return int(fval)
}

func mustBool(value interface{}, fallback bool) bool {
	val, ok := value.(bool)
	if !ok {
		return fallback
	}
	return val
}
