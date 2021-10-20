package store

import (
	"log"
	"time"
)

// Routes Store

type RoutesStats struct {
	Filtered int `json:"filtered"`
	Imported int `json:"imported"`
}

type RouteServerRoutesStats struct {
	Name   string      `json:"name"`
	Routes RoutesStats `json:"routes"`

	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoutesStoreStats struct {
	TotalRoutes  RoutesStats              `json:"total_routes"`
	RouteServers []RouteServerRoutesStats `json:"route_servers"`
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

// Neighbors Store

type RouteServerNeighborsStats struct {
	Name       string    `json:"name"`
	State      string    `json:"state"`
	Neighbors int       `json:"neighbors"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NeighborsStoreStats struct {
	TotalNeighbors int `json:"total_neighbors"`

	RouteServers []RouteServerNeighborsStats `json:"route_servers"`
}

// Print stats
func (stats NeighborsStoreStats) Log() {
	log.Println("Neighbors store:")

	log.Println("    Neighbors:",
		stats.TotalNeighbors)

	for _, rs := range stats.RouteServers {
		log.Println("      -", rs.Name)
		log.Println("        State:", rs.State)
		log.Println("        UpdatedAt:", rs.UpdatedAt)
		log.Println("        Neighbors:",
			rs.Neighbors)
	}
}
