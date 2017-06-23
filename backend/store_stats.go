package main

import (
	"log"
	"time"
)

type RoutesStats struct {
	Filtered int `json:"filtered"`
	Imported int `json:"imported"`
}

type RouteServerStats struct {
	Name   string      `json:"name"`
	Routes RoutesStats `json:"routes"`

	State     string
	UpdatedAt time.Time `json:"updated_at"`
}

type RoutesStoreStats struct {
	TotalRoutes  RoutesStats        `json:"total_routes"`
	RouteServers []RouteServerStats `json:"route_servers"`
}

// Write stats to the log
func (stats RoutesStoreStats) Log() {
	log.Println("Routes store:")

	log.Println("    Routes Imported:",
		stats.TotalRoutes.Imported,
		"Filtered:",
		stats.TotalRoutes.Filtered)
	log.Println("    Routeservers:")

	for _, rs := range stats.RouteServers {
		log.Println("      -", rs.Name)
		log.Println("        State:", rs.State)
		log.Println("        UpdatedAt:", rs.UpdatedAt)
		log.Println("        Routes Imported:",
			rs.Routes.Imported,
			"Filtered:",
			rs.Routes.Filtered)
	}
}
