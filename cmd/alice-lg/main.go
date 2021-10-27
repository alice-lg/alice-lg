package main

import (
	"context"
	"flag"
	"log"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/http"
	"github.com/alice-lg/alice-lg/pkg/store"
)

func main() {
	done := make(chan bool)
	ctx := context.Background()

	// Handle commandline parameters
	configFilenameFlag := flag.String(
		"config", "/etc/alice-lg/alice.conf",
		"Alice looking glass configuration file",
	)

	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFilenameFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Setup local routes store
	neighborsStore := store.NewNeighborsStore(cfg)
	routesStore := store.NewRoutesStore(neighborsStore, cfg)

	// Say hi
	printBanner(cfg, neighborsStore, routesStore)
	log.Println("Using configuration:", cfg.File)

	// Start stores
	if cfg.Server.EnablePrefixLookup == true {
		go neighborsStore.Start()
		go routesStore.Start()
	}

	// Start the Housekeeping
	go store.StartHousekeeping(ctx, cfg)

	// Start HTTP API
	server := http.NewServer(cfg, routesStore, neighborsStore)
	go server.Start()

	<-done
}
