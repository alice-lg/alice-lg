package config

import (
	"fmt"
	"regexp"
	"strings"
)

// Compile input pattern regex
var (
	expandMatchPlaceholder       = regexp.MustCompile(`(?U:{.*}+?)`)
	expandMatchWildcardShorthard = regexp.MustCompile(`(?U:{{.*\*}}+?)`)
)

// Extract all matches from the input string.
// The pattern to find is {INPUT}. The input string
// itself can contain new matches.
func expandFindPlaceholders(s string) []string {

	// Find all matches
	results := expandMatchPlaceholder.FindAllString(s, -1)
	if len(results) == 0 {
		return []string{}
	}

	matches := []string{}
	for _, result := range results {
		key := expandGetKey(result)
		subP := expandFindPlaceholders(key)
		matches = append(matches, result)
		matches = append(matches, subP...)
	}

	return matches
}

// Extract the key from the placeholder
func expandGetKey(s string) string {
	// Strip the enclosing curly braces
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	return s
}

// ExpandMap holds the current state of variables
type ExpandMap map[string]string

// Retrieve a set of matching variables, by iterating variables.
// Whenever a key matches the wildcard, the prefix is removed.
// Example:
//
//	pattern = "AS*", key = "AS2342", value = "2342"
func (e ExpandMap) matchWildcard(pattern string) []string {
	matches := []string{}

	// Strip the wildcard from the pattern.
	pattern = strings.TrimSuffix(pattern, "*")

	// Iterate variables and add match to result set
	for k := range e {
		if strings.HasPrefix(k, pattern) {
			key := strings.TrimPrefix(k, pattern)
			matches = append(matches, key)
		}
	}
	return matches
}

// Get all substitutions for a given key.
// This method will return an error, if a placeholder
// does not match.
func (e ExpandMap) getSubstitutions(key string) []string {
	// Check if the placeholder is a wildcard
	if strings.HasSuffix(key, "*") {
		return e.matchWildcard(key)
	}

	// Check if the placeholder is direct match
	if val, ok := e[key]; ok {
		return []string{val}
	}

	return []string{}
}

// Get placeholder level. This is the number of opening
// curly braces in the placeholder.
func expandGetLevel(s string) int {
	level := 0
	for _, c := range s {
		if c == '{' {
			level++
		}
	}
	return level
}

// Preprocess input string and resolve syntactic sugar.
// Replace {{VAR}} with {VAR{VAR}} to make it easier
// to access the wildcard value.
func expandPreprocess(s string) string {
	// Find all access shorthands and replace them
	// with the full syntax
	results := expandMatchWildcardShorthard.FindAllString(s, -1)
	for _, match := range results {
		// Wildcard {{KEY*}} -> KEY
		key := match[2 : len(match)-3]
		expr := fmt.Sprintf("{%s{%s*}}", key, key)
		s = strings.Replace(s, match, expr, -1)
	}
	return s
}

// Expand variables by recursive substitution and expansion
func (e ExpandMap) Expand(s string) ([]string, error) {
	// Preprocess syntactic sugar: replace {{VAR}}
	// with {VAR{VAR}}
	s = expandPreprocess(s)

	// Find all placeholders and substitute them
	placeholders := expandFindPlaceholders(s)
	if len(placeholders) == 0 {
		return []string{s}, nil
	}

	// Find substitutions for each placeholder
	substitutions := map[string][]string{}
	for _, p := range placeholders {
		key := expandGetKey(p)
		subs := e.getSubstitutions(key)
		if len(subs) == 0 {
			level := expandGetLevel(p)
			if level == 1 {
				err := fmt.Errorf("no substitution for %s in '%s'", p, s)
				return []string{}, err
			}
			continue
		}
		substitutions[p] = subs
	}

	// Apply substitutions
	subsRes := []string{s}
	for p, subs := range substitutions {
		subsExp := []string{}
		for _, s := range subsRes {
			for _, sub := range subs {
				res := strings.Replace(s, p, sub, -1)
				subsExp = append(subsExp, res)
			}
		}
		subsRes = subsExp
	}

	// Expand recursively
	results := []string{}
	for _, s := range subsRes {
		res, err := e.Expand(s)
		if err != nil {
			return []string{}, err
		}
		results = append(results, res...)
	}

	return results, nil
}

// AddExpr inserts a new variable to the map. Key and value are
// expanded.
func (e ExpandMap) AddExpr(expr string) error {
	// Expand expression
	res, err := e.Expand(expr)
	if err != nil {
		return err
	}
	for _, exp := range res {
		// Split key and value
		parts := strings.SplitN(exp, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid expression '%s'", expr)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		e[key] = val
	}

	return nil
}
