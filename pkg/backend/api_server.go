package backend

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// StartHTTPServer starts a HTTP server with the config
// in the global AliceConfig. TODO: refactor.
func StartHTTPServer() {

	// Setup request routing
	router := httprouter.New()

	// Serve static content
	if err := webRegisterAssets(AliceConfig.Ui, router); err != nil {
		log.Fatal(err)
	}

	if err := apiRegisterEndpoints(router); err != nil {
		log.Fatal(err)
	}

	httpTimeout := time.Duration(AliceConfig.Server.HttpTimeout) * time.Second
	log.Println("Web server HTTP timeout set to:", httpTimeout)

	server := &http.Server{
		Addr:         AliceConfig.Server.Listen,
		Handler:      router,
		ReadTimeout:  httpTimeout,
		WriteTimeout: httpTimeout,
		IdleTimeout:  httpTimeout,
	}

	// Start http server
	log.Fatal(server.ListenAndServe())
}
