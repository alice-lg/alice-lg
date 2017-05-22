package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var AliceConfig *Config

func main() {
	var err error

	// Handle commandline parameters
	configFilenameFlag := flag.String(
		"config", "/etc/alicelg/alice.conf",
		"Alice looking glass configuration file",
	)

	flag.Parse()

	// Load configuration
	AliceConfig, err = loadConfig(*configFilenameFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Say hi
	printBanner()

	log.Println("Using configuration:", AliceConfig.File)

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
	log.Fatal(http.ListenAndServe(AliceConfig.Server.Listen, router))
}
