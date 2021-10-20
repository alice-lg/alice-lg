package main

import (
	"flag"
	"log"

	"github.com/alice-lg/alice-lg/pkg/store"
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
	cfg, err = config.LoadConfig(filename)
  if err != nil {
		log.Fatal(err)
	}

	// Say hi
	printBanner()
	log.Println("Using configuration:", cfg.File)

	// Setup local routes store
	neighborsStore := store.NewNeighborsStore(cfg)
	routesStore := store.NewRoutesStore(cfg)

	// Start stores
	if backend.AliceConfig.Server.EnablePrefixLookup == true {
		go neighborsStore.Start()
    go routesStore.Start()
	}

	// Start the Housekeeping
	go store.Housekeeping(cfg)

  // Start HTTP API
	go backend.StartHTTPServer()

	<-quit
}
