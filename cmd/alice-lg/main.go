package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/backend"
)

func main() {
	// Handle commandline parameters
	configFilenameFlag := flag.String(
		"config", "/etc/alice-lg/alice.conf",
		"Alice looking glass configuration file",
	)

	flag.Parse()

	// Load configuration
	if err := backend.InitConfig(*configFilenameFlag); err != nil {
		log.Fatal(err)
	}

	// Say hi
	printBanner()

	log.Println("Using configuration:", AliceConfig.File)

	// Setup local routes store
	AliceRoutesStore = NewRoutesStore(AliceConfig)

	if AliceConfig.Server.EnablePrefixLookup == true {
		AliceRoutesStore.Start()
	}

	// Setup local neighbours store
	AliceNeighboursStore = NewNeighboursStore(AliceConfig)
	if AliceConfig.Server.EnablePrefixLookup == true {
		AliceNeighboursStore.Start()
	}

	// Start the Housekeeping
	go Housekeeping(AliceConfig)

	// Setup request routing
	router := httprouter.New()

	// Serve static content
	err = webRegisterAssets(AliceConfig.Ui, router)
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
