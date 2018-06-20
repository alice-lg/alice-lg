package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
)

// Web Client
// Handle assets and client app preprarations

// Prepare client HTML:
// Set paths and add version to assets.
func webPrepareClientHtml(html string) string {
	status, _ := NewAppStatus()

	// Replace paths and tags
	rewriter := strings.NewReplacer(
		// Paths
		"js/", "/static/js/",
		"css/", "/static/css/",

		// Tags
		"APP_VERSION", status.Version,
	)
	html = rewriter.Replace(html)
	return html
}

func webNewThemeHandler(path string) http.Handler {
	if path == "" {
		return nil // Nothing to do here (Null Theme)
	}

	// Check if the the path is present and readable
	if _, err := os.Stat(path); err != nil {
		log.Println("WARNING - Could not read theme path:", path)
		return nil
	}

	log.Println("Using theme at:", path)

	// Looks like we are okay
	// Serve the content using the file server
	themeFilesHandler := http.StripPrefix(
		"/theme",
		http.FileServer(http.Dir(path)))

	return themeFilesHandler
}

// Register assets handler and index handler
// at /static and /
func webRegisterAssets(ui UiConfig, router *httprouter.Router) error {
	log.Println("Preparing and installing assets")

	// Serve static assets
	assets := rice.MustFindBox("../client/build")
	assetsHandler := http.StripPrefix(
		"/static/",
		http.FileServer(assets.HTTPBox()))

	// Prepare client html: Rewrite paths
	indexHtml, err := assets.String("index.html")
	if err != nil {
		return err
	}

	// Update paths
	indexHtml = webPrepareClientHtml(indexHtml)

	themeHandler := webNewThemeHandler(ui.Theme.Path)
	if themeHandler != nil {
		// We have a theme
		router.Handler("GET", "/theme/*path", themeHandler)
	}

	// Register static assets
	router.Handler("GET", "/static/*path", assetsHandler)

	// Rewrite paths
	// Serve index html as root...
	router.GET("/", func(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		io.WriteString(res, indexHtml)
	})

	// ...and as catch all
	router.GET("/alice/*path", func(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		io.WriteString(res, indexHtml)
	})

	return nil
}
