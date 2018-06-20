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
	"strings"

	"io/ioutil"
	"net/http"
)

type Theme struct {
	Config ThemeConfig
}

func NewTheme(config ThemeConfig) (*Theme, error) {
	theme := &Theme{
		Config: config,
	}

	// Check if the the path is present and readable
	path := theme.Config.Path
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("Theme path could not be found: %s", path)
	}

	return theme, nil
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

		if strings.HasSuffix(file.Name(), suffix) {
			include := self.Config.BasePath + "/" + file.Name()
			includes = append(includes, include)
		}
	}

	return includes
}

/*
 Retrieve a list of includeable stylesheets
*/
func (self *Theme) Stylesheets() []string {
	return self.listIncludes(".css")
}

/*
 Retrieve a list of includeable javascipts
*/
func (self *Theme) Scripts() []string {
	return self.listIncludes(".js")
}

func (self *Theme) Handler() http.Handler {
	log.Println("Using theme at:", self.Config.Path)

	// Serve the content using the file server
	path := self.Config.Path
	themeFilesHandler := http.StripPrefix(
		self.Config.BasePath, http.FileServer(http.Dir(path)))

	return themeFilesHandler
}
