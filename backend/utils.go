package main

// Some helper functions
import (
	"regexp"
	"strconv"
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

/*
 Since havin ints as keys in json is
 acutally undefined behaviour, we keep these interally
 but provide a string as a key for serialization
*/
func SerializeReasons(reasons map[int]string) map[string]string {
	res := make(map[string]string)
	for id, reason := range reasons {
		res[strconv.Itoa(id)] = reason
	}
	return res
}

/*
 Make trimmed list of CSV strings.
 Ommits empty values.
*/
func TrimmedStringList(s string) []string {
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
