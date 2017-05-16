package main

import (
	"github.com/GeertJohan/go.rice"

	"log"
	"net/http"
)

func main() {
	printBanner()

	// Load configuration

	// Serve static assets
	assets := rice.MustFindBox("../client/build")
	assetsHandler := http.StripPrefix(
		"/static/",
		http.FileServer(assets.HTTPBox()))

	index, err := assets.String("index.html")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(index)

	http.Handle("/static/", assetsHandler)

	// Start http server
	http.ListenAndServe(":7340", nil)
}
