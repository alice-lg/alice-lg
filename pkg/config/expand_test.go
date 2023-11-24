package config

import (
	"testing"
)

// Text variable pattern matching
func TestExpandMatch(t *testing.T) {
	exp := ExpandMap{
		"AS2342": "",
		"AS1111": "",
		"FOOBAR": "foo",
	}

	matches := exp.matchWildcard("AS*")
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}

	for _, m := range matches {
		t.Log("Match wildcard:", m)
	}
}

// Test variable expansion / substitution
func TestFindPlaceholders(t *testing.T) {
	s := "{FOO} BAR {AS{AS*}}"
	placeholders := expandFindPlaceholders(s)
	if len(placeholders) != 3 {
		t.Errorf("Expected 3 placeholders, got %d", len(placeholders))
	}
	t.Log(placeholders)
}

// Test variable expansion / substitution
func TestExpand(t *testing.T) {
	s := "{FOO} BAR {AS{AS*}} AS {AS*}"
	exp := ExpandMap{
		"AS2342": "AS2342",
		"AS1111": "AS1111",
		"FOO":    "foo",
	}

	results, err := exp.Expand(s)
	if err != nil {
		t.Error(err)
	}
	t.Log(results)
}

func TestExpandErr(t *testing.T) {
	s := "{FOO} BAR {AS{AS*}} AS {AS*} {UNKNOWN}"
	exp := ExpandMap{
		"AS2342": "AS2342",
		"AS1111": "AS1111",
		"FOO":    "foo",
		"FN":     "fn",
		"FA":     "fa",
	}

	_, err := exp.Expand(s)
	t.Log(err)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestExpandPreprocess(t *testing.T) {
	s := "FOO {FOO} {{AS*}} {F*} {{F*}} {X{X*}}"
	expect := "FOO {FOO} {AS{AS*}} {F*} {F{F*}} {X{X*}}"
	s = expandPreprocess(s)
	if s != expect {
		t.Errorf("Expected '%s', got '%s'", expect, s)
	}
	t.Log(s)

	s = "TEST {{FN}}"
	s = expandPreprocess(s)
	t.Log(s)

}

func TestExpandAddExpr(t *testing.T) {
	e := ExpandMap{
		"FOO":   "foo23",
		"BAR":   "bar42",
		"bar42": "BAM",
	}

	if err := e.AddExpr("FOOBAR = {FOO}{BAR}{{BAR}}"); err != nil {
		t.Fatal(err)
	}
	t.Log(e)

	if e["FOOBAR"] != "foo23bar42BAM" {
		t.Error("Expected 'foo23bar42BAM', got", e["FOOBAR"])
	}
}

func TestExpandBgpCommunities(t *testing.T) {
	e := ExpandMap{
		"ASRS01":   "6695",
		"ASRS02":   "4617",
		"SW1001":   "edge01.fra2",
		"SW1002":   "edge01.fra6",
		"SW2038":   "edge01.nyc1",
		"RDCTL911": "Redistribute",
		"RDCTL922": "Do not redistribute",
	}

	// Some large communities:
	expr := "{{AS*}}:{RDCTL*}:{SW*} = {{RDCTL*}} to {{SW*}}"
	exp, err := e.Expand(expr)
	if err != nil {
		t.Fatal(err)
	}
	expected := 2 * 3 * 2
	if len(exp) != expected {
		t.Errorf("Expected %d results, got %d", expected, len(exp))
	}
	t.Log(exp)
}
