package main

/*
 The theme provides a method for adding customized CSS
 or Javascript to Alice:

 A theme directory can be specified in the config.
 Stylesheets and Javascript residing in the theme root
 directory will be included in the frontends HTML.

 Additional files can be added in subdirectories.
 These are served aswell and can be used for additional
 assets. (E.g. a logo)
*/

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

type Theme struct {
	Config ThemeConfig
}

func NewTheme(config ThemeConfig) *Theme {
	theme := &Theme{
		Config: config,
	}

	return theme
}

/*
 Get includable files from theme directory
*/
func (self *Theme) listIncludes(suffix string) []string {
	includes := []string{}

	files, err := ioutil.ReadDir(self.Config.Path)
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

/*
Calculate a hashvalue for an include file,
to help with cache invalidation, when the file changes.

We are using the timestamp of the last access as Unix()
encoded as hex
*/
func (self *Theme) HashInclude(include string) string {
	path := filepath.Join(self.Config.Path, include)
	stat, err := os.Stat(path)
	if err != nil {
		return ""
	}

	modTime := stat.ModTime().UTC()
	timestamp := modTime.Unix()

	return strconv.FormatInt(timestamp, 16)
}

/*
 Retrieve a list of includeable stylesheets, with
 their md5sum as hash
*/
func (self *Theme) Stylesheets() []string {
	return self.listIncludes(".css")
}

/*
 Make include statement: stylesheet
*/
func (self *Theme) StylesheetIncludes() string {

	includes := []string{}
	for _, stylesheet := range self.Stylesheets() {
		hash := self.HashInclude(stylesheet)
		include := fmt.Sprintf(
			"<link rel=\"stylesheet\" href=\"%s/%s?%s\" />",
			self.Config.BasePath, stylesheet, hash,
		)
		includes = append(includes, include)
	}

	return strings.Join(includes, "\n")
}

/*
 Retrieve a list of includeable javascipts
*/
func (self *Theme) Scripts() []string {
	return self.listIncludes(".js")
}

/*
 Make include statement: script
*/
func (self *Theme) ScriptIncludes() string {
	includes := []string{}
	for _, script := range self.Scripts() {
		hash := self.HashInclude(script)
		include := fmt.Sprintf(
			"<script type=\"text/javascript\" src=\"%s/%s?%s\"></script>",
			self.Config.BasePath, script, hash,
		)
		includes = append(includes, include)
	}

	return strings.Join(includes, "\n")
}

/*
 Theme HTTP Handler
*/
func (self *Theme) Handler() http.Handler {

	// Serve the content using the file server
	path := self.Config.Path
	themeFilesHandler := http.StripPrefix(
		self.Config.BasePath, http.FileServer(http.Dir(path)))

	return themeFilesHandler
}

/*
 Register theme at path
*/
func (self *Theme) RegisterThemeAssets(router *httprouter.Router) error {
	fsPath := self.Config.Path
	if fsPath == "" {
		return nil // nothing to do here
	}

	if _, err := os.Stat(fsPath); err != nil {
		return fmt.Errorf("Theme path '%s' could not be found!", fsPath)
	}

	log.Println("Using theme at:", fsPath)

	// We have a theme, install handler
	path := fmt.Sprintf("%s/*path", self.Config.BasePath)
	router.Handler("GET", path, self.Handler())

	return nil
}

/*
 Prepare document, fill placeholder with scripts and stylesheet
*/
func (self *Theme) PrepareClientHtml(html string) string {
	stylesheets := self.StylesheetIncludes()
	scripts := self.ScriptIncludes()

	html = strings.Replace(html,
		"<!-- ###THEME_STYLESHEETS### -->",
		stylesheets, 1)
	html = strings.Replace(html,
		"<!-- ###THEME_SCRIPTS### -->",
		scripts, 1)

	return html
}
