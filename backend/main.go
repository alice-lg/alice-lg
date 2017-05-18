package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var AliceConfig *Config

func main() {
	var err error

	printBanner()

	// Load configuration
	AliceConfig, err = loadConfigs("../etc/alicelg/alice.conf", "", "")
	log.Println("Using configuration: ...")

	// Setup request routing
	router := httprouter.New()

	// Serve static content
	err = httpRegisterAssets(router)
	if err != nil {
		log.Fatal(err)
	}

	err = apiRegisterEndpoints(router)
	if err != nil {
		log.Fatal(err)
	}

	// Start http server
	log.Fatal(http.ListenAndServe(":7340", router))
}
