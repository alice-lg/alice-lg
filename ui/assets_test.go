package ui

import (
	"testing"
)

func TestPresenceOfIndexHTML(t *testing.T) {
	content, err := Assets.ReadFile("build/index.html")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
}
