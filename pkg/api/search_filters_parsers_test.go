package api

import (
	"testing"
)

func TestParseQueryValueIntList(t *testing.T) {
	res, err := parseQueryValueList(parseIntValue, "23,42, 10, 11")
	if err != nil {
		t.Error(err)
		return
	}

	expected := []int{23, 42, 10, 11}
	for i, _ := range expected {
		if res[i].Value.(int) != expected[i] {
			t.Error("Expected:", expected[i], "got:", res[i])
		}
	}

	_, err = parseQueryValueList(parseIntValue, "23,a,b,42")
	if err == nil {
		t.Error("Expected err to be not nil with invalid integer list")
	}
}

func TestParseCommunityValueList(t *testing.T) {
	res, err := parseQueryValueList(parseCommunityValue, "23:42,10:42:123")
	if err != nil {
		t.Error(err)
	}
	if len(res) != 2 {
		t.Error("Expected length 2, got:", len(res))
	}

	filter := res[1]
	if filter.Name != "10:42:123" {
		t.Error("Expected name: '10:42:123', but got:", filter.Name)
	}
	com := filter.Value.(Community)
	if com[0] != 10 || com[1] != 42 || com[2] != 123 {
		t.Error("Expected [10, 42, 123] but got:", com)
	}
}

func TestParseExtCommunityValue(t *testing.T) {
	filter, err := parseExtCommunityValue("rt:23:42")
	if err != nil {
		t.Error(err)
		return
	}

	com := filter.Value.(ExtCommunity)

	if com[0].(string) != "rt" &&
		com[1].(string) != "23" &&
		com[2].(string) != "42" {
		t.Error("Expected community to be: ['rt', '23', '42'] but got:", com)
	}

}
