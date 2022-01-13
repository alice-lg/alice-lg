package main

import (
	"context"
	"flag"
	"log"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/http"
	"github.com/alice-lg/alice-lg/pkg/store"
	"github.com/alice-lg/alice-lg/pkg/store/backends/memory"
	"github.com/alice-lg/alice-lg/pkg/store/backends/postgres"
)

func main() {
	ctx := context.Background()

	// Handle commandline parameters
	configFilenameFlag := flag.String(
		"config", "/etc/alice-lg/alice.conf",
		"Alice looking glass configuration file",
	)
	dbInitFlag := flag.Bool(
		"db-init", false,
		"Initialize the database. Clears all data.",
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFilenameFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Setup local routes store and use backend from configuration
	var (
		neighborsBackend store.NeighborsStoreBackend = memory.NewNeighborsBackend()
		routesBackend    store.RoutesStoreBackend    = memory.NewRoutesBackend()
	)
	if cfg.Server.StoreBackend == "postgres" {
		pool, err := postgres.Connect(ctx, cfg.Postgres)
		if err != nil {
			log.Fatal(err)
		}
		m := postgres.NewManager(pool)

		// Initialize db if required
		if *dbInitFlag {
			if err := m.Initialize(ctx); err != nil {
				log.Fatal(err)
			}
			log.Println("database initialized")
			return
		}

		go m.Start(ctx)

		neighborsBackend = postgres.NewNeighborsBackend(pool)
		routesBackend = postgres.NewRoutesBackend(pool)
	}

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
