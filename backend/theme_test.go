package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func touchFile(path, filename string) error {
	target := filepath.Join(path, filename)
	return ioutil.WriteFile(target, []byte{}, 0644)
}

func TestThemeLoading(t *testing.T) {
	themePath, err := ioutil.TempDir("", "alice-lg-tmp-theme")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(themePath)

	// This should work aswell, as themes are optional
	_, err = NewTheme(ThemeConfig{
		BasePath: "/theme",
		Path:     themePath,
	})

	if err != nil {
		t.Error(err)
	}

	// This should not:
	_, err = NewTheme(ThemeConfig{
		Path: "/1ade5e183fd7b84a1590ad7144dbd6e0caed1b6a",
	})

	if err == nil {
		t.Error("Expected the theme loading to fail with unknown path.")
	}
}

func TestThemeFiles(t *testing.T) {
	themePath, err := ioutil.TempDir("", "alice-lg-tmp-theme")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(themePath)

	// Create some "stylesheets" and a "script"
	touchFile(themePath, "style.css")
	touchFile(themePath, "extra.css")
	touchFile(themePath, "script.js")

	// Load theme
	theme := NewTheme(ThemeConfig{
		BasePath: "/theme",
		Path:     themePath,
	})

	if err != nil {
		t.Error(err)
	}

	// Check file presence
	scripts := theme.Scripts()
	if len(scripts) != 1 {
		t.Error("Expected one script file: script.js")
	}

	stylesheets := theme.Stylesheets()
	if len(stylesheets) != 2 {
		t.Error("Expected two stylesheets: {style, extra}.css")
	}

	// Check uri / path mapping
	script := scripts[0]
	if script != "script.js" {
		t.Error("Expected script.js to be included in scripts")
	}
}

func TestThemeIncludeHash(t *testing.T) {
	themePath, err := ioutil.TempDir("", "alice-lg-tmp-theme")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(themePath)

	// Create some "stylesheets" and a "script"
	touchFile(themePath, "style.css")

	theme := NewTheme(ThemeConfig{
		BasePath: "/theme",
		Path:     themePath,
	})

	hash := theme.HashInclude("style.css")
	if hash == "" {
		t.Error("Something went wrong with hashing")
	}

	t.Log("Filehash:", hash)

}

func TestThemeIncludes(t *testing.T) {
	themePath, err := ioutil.TempDir("", "alice-lg-tmp-theme")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(themePath)

	// Create some "stylesheets" and a "script"
	touchFile(themePath, "style.css")
	touchFile(themePath, "extra.css")
	touchFile(themePath, "script.js")

	// Load theme
	theme := NewTheme(ThemeConfig{
		BasePath: "/theme",
		Path:     themePath,
	})

	stylesHtml := theme.StylesheetIncludes()
	scriptsHtml := theme.ScriptIncludes()

}
