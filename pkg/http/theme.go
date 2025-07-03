package http

/*
 The theme provides a method for adding customized CSS
 or Javascript to Alice:

 A theme directory can be specified in the config.
 Stylesheets and Javascript residing in the theme root
 directory will be included in the frontends HTML.

 Additional files can be added in subdirectories.
 These are served as well and can be used for additional
 assets. (E.g. a logo)
*/

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/config"
)

// Theme is a client customization through additional
// HTML, CSS and JS content.
type Theme struct {
	Config config.ThemeConfig
}

// NewTheme creates a theme from a config
func NewTheme(config config.ThemeConfig) *Theme {
	return &Theme{
		Config: config,
	}
}

// Get includable files from theme directory
func (t *Theme) listIncludes(suffix string) []string {
	includes := []string{}

	files, err := os.ReadDir(t.Config.Path)
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if strings.HasPrefix(filename, ".") {
			continue
		}

		if strings.HasSuffix(filename, suffix) {
			includes = append(includes, filename)
		}
	}
	return includes
}

// HashInclude calculates a hashvalue for an include file,
// to help with cache invalidation, when the file changes.
//
// We are using the timestamp of the last access as Unix()
// encoded as hex
func (t *Theme) HashInclude(include string) string {
	path := filepath.Join(t.Config.Path, include)
	stat, err := os.Stat(path)
	if err != nil {
		return ""
	}

	modTime := stat.ModTime().UTC()
	timestamp := modTime.Unix()

	return strconv.FormatInt(timestamp, 16)
}

// Stylesheets retrieve a list of includeable stylesheets, with
// their md5sum as hash
func (t *Theme) Stylesheets() []string {
	return t.listIncludes(".css")
}

// StylesheetIncludes make include statements for stylesheet
func (t *Theme) StylesheetIncludes() string {

	includes := []string{}
	for _, stylesheet := range t.Stylesheets() {
		hash := t.HashInclude(stylesheet)
		include := fmt.Sprintf(
			"<link rel=\"stylesheet\" href=\"%s/%s?%s\" />",
			t.Config.BasePath, stylesheet, hash,
		)
		includes = append(includes, include)
	}

	return strings.Join(includes, "\n")
}

// Scripts retrieve a list of includeable javascripts
func (t *Theme) Scripts() []string {
	return t.listIncludes(".js")
}

// ScriptIncludes makes include statement for scripts
func (t *Theme) ScriptIncludes() string {
	includes := []string{}
	for _, script := range t.Scripts() {
		hash := t.HashInclude(script)
		include := fmt.Sprintf(
			"<script type=\"text/javascript\" src=\"%s/%s?%s\" defer></script>",
			t.Config.BasePath, script, hash,
		)
		includes = append(includes, include)
	}

	return strings.Join(includes, "\n")
}

// Handler is the theme HTTP handler
func (t *Theme) Handler() http.Handler {
	// Serve the content using the file server
	path := t.Config.Path
	themeFilesHandler := http.StripPrefix(
		t.Config.BasePath, http.FileServer(http.Dir(path)))
	return themeFilesHandler
}

// RegisterThemeAssets registers the theme at path
func (t *Theme) RegisterThemeAssets(router *httprouter.Router) error {
	fsPath := t.Config.Path
	if fsPath == "" {
		return nil // nothing to do here
	}

	if _, err := os.Stat(fsPath); err != nil {
		return fmt.Errorf("theme path '%s' could not be found", fsPath)
	}

	log.Println("Using theme at:", fsPath)

	// We have a theme, install handler
	path := fmt.Sprintf("%s/*path", t.Config.BasePath)
	router.Handler("GET", path, t.Handler())

	return nil
}

// PrepareClientHTML prepares the document and fills placeholders
// with scripts and stylesheet.
func (t *Theme) PrepareClientHTML(html string) string {
	stylesheets := t.StylesheetIncludes()
	scripts := t.ScriptIncludes()

	html = strings.Replace(html,
		"<!-- ###THEME_STYLESHEETS### -->",
		stylesheets, 1)
	html = strings.Replace(html,
		"<!-- ###THEME_SCRIPTS### -->",
		scripts, 1)

	return html
}
