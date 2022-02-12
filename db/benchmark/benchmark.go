package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store/backends/postgres"
)

func makeRoute() *api.LookupRoute {
	jjj
}

func makeRoutes(count uint) api.LookupRoutes {
	routes := make(api.LookupRoutes, 0, count)

	return routes
}

func main() {
	fmt.Println("benchmarking routes insert")

	configFilenameFlag := flag.String(
		"config", "/etc/alice-lg/alice.conf",
		"Alice looking glass configuration file",
	)

	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig(*configFilenameFlag)
	if err != nil {
		log.Fatal(err)
	}
	pool, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	m := postgres.NewManager(pool)

	if err := m.Initialize(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("database initialized")

	go m.Start(ctx)
	backend := postgres.NewRoutesBackend(pool)

	// Now insert tons of routes...
	routes := makeRoutes(1000)
	for i := 0; i < 10; i++ {
		t := time.Now()
		if err := backend.SetRoutes(
			ctx, "rs1-example-fra1", routes); err != nil {
			log.Fatal(err)
		}

		elapsed := time.Since(t)
		log.Println(
			"set routes", i, "finished after:", elapsed)
	}
}
