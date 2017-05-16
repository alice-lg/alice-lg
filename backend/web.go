package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
)

// Web Client
// Handle assets and client app preprarations

// Register assets handler and index handler
// at /static and /
func httpRegisterAssets(router *httprouter.Router) error {
	log.Println("Preparing and installing assets")

	// Serve static assets
	assets := rice.MustFindBox("../client/build")
	assetsHandler := http.StripPrefix(
		"/static/",
		http.FileServer(assets.HTTPBox()))

	// Register static assets
	router.Handler("GET", "/static/*path", assetsHandler)

	// Prepare client html: Rewrite paths
	indexHtml, err := assets.String("index.html")
	if err != nil {
		return err
	}

	pathRewriter := strings.NewReplacer(
		"js/", "/static/js/",
		"css/", "/static/css/")
	indexHtml = pathRewriter.Replace(indexHtml)

	// Rewrite paths
	// Serve index html as root
	router.GET("/", func(res http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		io.WriteString(res, indexHtml)
	})

	return nil
}
