package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/http"
	"github.com/alice-lg/alice-lg/pkg/store"
	"github.com/alice-lg/alice-lg/pkg/store/backends/memory"
	"github.com/alice-lg/alice-lg/pkg/store/backends/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

func createHeapProfile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func createAllocProfile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create alloc profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
		log.Fatal("could not write alloc profile: ", err)
	}
}

func startMemoryProfile(prefix string) {
	t := 0
	for {
		filename := fmt.Sprintf("%s-heap-%03d", prefix, t)
		runtime.GC() // get up-to-date statistics (according to docs)
		createHeapProfile(filename)
		log.Println("wrote memory heap profile:", filename)
		filename = fmt.Sprintf("%s-allocs-%03d", prefix, t)
		log.Println("wrote memory allocs profile:", filename)
		createAllocProfile(filename)
		time.Sleep(30 * time.Second)
		t++
	}
}

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
	memprofile := flag.String(
		"memprofile", "", "write memory profile to `file`",
	)
	flag.Parse()

	if *memprofile != "" {
		go startMemoryProfile(*memprofile)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configFilenameFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Tune garbage collection
	debug.SetGCPercent(10)

	// Setup local routes store and use backend from configuration
	var (
		neighborsBackend store.NeighborsStoreBackend = memory.NewNeighborsBackend()
		routesBackend    store.RoutesStoreBackend    = memory.NewRoutesBackend()

		pool *pgxpool.Pool
	)
	if cfg.Server.StoreBackend == "postgres" {
		pool, err = postgres.Connect(ctx, cfg.Postgres)
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
		routesBackend = postgres.NewRoutesBackend(
			pool, cfg.Sources)
		if err := routesBackend.(*postgres.RoutesBackend).Init(ctx); err != nil {
			log.Println("error while initializing routes backend:", err)
		}
	}

	neighborsStore := store.NewNeighborsStore(cfg, neighborsBackend)
	routesStore := store.NewRoutesStore(neighborsStore, cfg, routesBackend)

	// Say hi
	printBanner(cfg, neighborsStore, routesStore)
	log.Println("Using configuration:", cfg.File)

	// Start stores
	if cfg.Server.EnablePrefixLookup {
		go neighborsStore.Start(ctx)
		go routesStore.Start(ctx)
	}

	// Start exporting metrics
	if cfg.Server.EnableMetrics {
		go store.StartMetrics(ctx, neighborsStore)
	}

	// Start the Housekeeping
	go store.StartHousekeeping(ctx, cfg)

	// Start HTTP API
	server := http.NewServer(cfg, pool, routesStore, neighborsStore)
	go server.Start(ctx)

	<-ctx.Done()
}
