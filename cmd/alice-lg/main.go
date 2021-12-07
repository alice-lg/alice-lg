package main

import (
	"context"
	"flag"
	"log"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/http"
	"github.com/alice-lg/alice-lg/pkg/store"
	"github.com/alice-lg/alice-lg/pkg/store/backends/memory"
)

func main() {
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
	neighborsBackend := memory.NewNeighborsBackend()
	routesBackend := memory.NewRoutesBackend()

	neighborsStore := store.NewNeighborsStore(cfg, neighborsBackend)
	routesStore := store.NewRoutesStore(neighborsStore, cfg, routesBackend)

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

	<-ctx.Done()
}
