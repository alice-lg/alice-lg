package backend

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

// Neighbours Store

type RouteServerNeighboursStats struct {
	Name       string    `json:"name"`
	State      string    `json:"state"`
	Neighbours int       `json:"neighbours"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NeighboursStoreStats struct {
	TotalNeighbours int `json:"total_neighbours"`

	RouteServers []RouteServerNeighboursStats `json:"route_servers"`
}

// Print stats
func (stats NeighboursStoreStats) Log() {
	log.Println("Neighbours store:")

	log.Println("    Neighbours:",
		stats.TotalNeighbours)

	for _, rs := range stats.RouteServers {
		log.Println("      -", rs.Name)
		log.Println("        State:", rs.State)
		log.Println("        UpdatedAt:", rs.UpdatedAt)
		log.Println("        Neighbours:",
			rs.Neighbours)
	}
}
