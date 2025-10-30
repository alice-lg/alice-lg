package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store/backends/postgres"
)

func randomNet() string {
	return fmt.Sprintf("fd42:%d:%d:%d::%d",
		1+rand.Intn(9999),
		1+rand.Intn(9999),
		1+rand.Intn(9999),
		1+rand.Intn(99))

}

func makeRoute(n int) *api.LookupRoute {
	//id := fmt.Sprintf("route_%d", n)
	nid := fmt.Sprintf("neighbor_%d", n%50)
	gw := fmt.Sprintf("fd23:2342:%04d::1", n%50)
	intf := "enp0s23"
	net := randomNet()
	return &api.LookupRoute{
		Route: &api.Route{
			//ID:         id,
			NeighborID: &nid,
			Network:    net,
			Interface:  &intf,
			Gateway:    &gw,
			Metric:     100,
			Age:        30 * time.Second,
			Type:       []string{"BGP", "unicast", "univ"},
			Primary:    true,
			LearntFrom: &gw,
		},
		State: "imported",
		Neighbor: &api.Neighbor{
			ID:              nid,
			Address:         gw,
			State:           "up",
			Description:     fmt.Sprintf("Neighbor AS 2342%d", n%50),
			RoutesReceived:  23000,
			RoutesFiltered:  0,
			RoutesExported:  42000,
			RoutesPreferred: 23,
			RoutesAccepted:  23910,
			Uptime:          20 * time.Minute,
			RouteServerID:   "rsid",
		},
	}
}

func makeRoutes(count int) api.LookupRoutes {
	routes := make(api.LookupRoutes, 0, count)
	for i := range count {
		routes = append(routes, makeRoute(i))
	}
	return routes
}

func main() {
	fmt.Println("benchmarking routes insert")

	configFilenameFlag := flag.String(
		"config", "/etc/alice-lg/alice.conf",
		"Alice looking glass configuration file",
	)

	flag.Parse()

	ctx := context.Background()
	fmt.Println("using config:", *configFilenameFlag)

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
	backend := postgres.NewRoutesBackend(pool, cfg.Sources)

	// Now insert tons of routes...
	for i := range 10 {
		routes := makeRoutes(100000)
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
