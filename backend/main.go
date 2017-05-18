package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	printBanner()

	// Load configuration
	log.Println("Using configuration: ...")

	// Setup request routing
	router := httprouter.New()

	// Serve static content
	err := httpRegisterAssets(router)
	if err != nil {
		log.Fatal(err)
	}

	err = apiRegisterEndpints(router)
	if err != nil {
		log.Fatal(err)
	}

	// Start http server
	log.Fatal(http.ListenAndServe(":7340", router))
}
