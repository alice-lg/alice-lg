package api

import (
	"log"
	"time"
)

// RoutesStats provides number of filtered and
// imported prefixes in the store
type RoutesStats struct {
	Filtered int `json:"filtered"`
	Imported int `json:"imported"`
}

// RouteServerRoutesStats provides the number of
// filtered and exported routes for a route server.
type RouteServerRoutesStats struct {
	Name   string      `json:"name"`
	Routes RoutesStats `json:"routes"`

	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoutesStoreStats are statistics about the
// stored prefixes per route server
type RoutesStoreStats struct {
	TotalRoutes  RoutesStats              `json:"total_routes"`
	RouteServers []RouteServerRoutesStats `json:"route_servers"`
}

// Log writes stats to the log
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

// RouteServerNeighborsStats are statistics about the
// neighbors store for a single route server.
type RouteServerNeighborsStats struct {
	Name      string    `json:"name"`
	State     string    `json:"state"`
	Neighbors int       `json:"neighbors"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NeighborsStoreStats are the stats for all sources
// of neighbors.
type NeighborsStoreStats struct {
	TotalNeighbors int                         `json:"total_neighbors"`
	RouteServers   []RouteServerNeighborsStats `json:"route_servers"`
}

// Log prints the stats
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
