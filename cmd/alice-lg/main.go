package main

import (
	"flag"
	"log"

	"github.com/alice-lg/alice-lg/pkg/backend"
)

func main() {
	quit := make(chan bool)

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

	log.Println("Using configuration:", backend.AliceConfig.File)

	// Setup local routes store
	backend.InitStores()

	// Start stores
	if backend.AliceConfig.Server.EnablePrefixLookup == true {
		go backend.AliceRoutesStore.Start()
	}

	// Setup local neighbours store
	if backend.AliceConfig.Server.EnablePrefixLookup == true {
		go backend.AliceNeighboursStore.Start()
	}

	// Start the Housekeeping
	go backend.Housekeeping(backend.AliceConfig)

	go backend.StartHTTPServer()

	<-quit
}
