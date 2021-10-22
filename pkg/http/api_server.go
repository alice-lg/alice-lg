package http

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// StartHTTPServer starts a HTTP server with the config
// in the global AliceConfig. TODO: refactor.
func StartHTTPServer() {

	// Setup request routing
	router := httprouter.New()

	// Serve static content
	if err := webRegisterAssets(AliceConfig.UI, router); err != nil {
		log.Fatal(err)
	}

	if err := apiRegisterEndpoints(router); err != nil {
		log.Fatal(err)
	}

	// Start http server
	log.Fatal(http.ListenAndServe(AliceConfig.Server.Listen, router))
}
