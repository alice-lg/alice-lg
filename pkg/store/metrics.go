package store

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	neighborsStore *NeighborsStore

	neighborInfo   *prometheus.GaugeVec
	neighborUptime *prometheus.GaugeVec

	routesReceived  *prometheus.GaugeVec
	routesFiltered  *prometheus.GaugeVec
	routesPreferred *prometheus.GaugeVec
	routesAccepted  *prometheus.GaugeVec
}

// Initialize
func initMetrics(s *NeighborsStore) *metrics {
	log.Println(
		"[metrics] Initializing export")

	labels := []string{
		// The route server ID
		"route_server_id",
		"route_server_name",
		"route_server_group",
		"neighbor_id",
		"neighbor_description",
		"neighbor_asn",
		"neighbor_address",
	}

	neighborInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "neighbor_info",
			Help: "Information about the neighbor including the state",
		},
		append(labels, "neighbor_state"),
	)

	neighborUptime := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "neighbor_uptime_seconds_total",
			Help: "The uptime of a neighbor on a route server in seconds",
		},
		labels,
	)

	routesReceived := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "routes_received",
			Help: "Total number of routes received by a route server for a given neighbor",
		},
		labels,
	)

	routesFiltered := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "routes_filtered",
			Help: "Total number of routes filtered by a route server for a given neighbor",
		},
		labels,
	)

	routesPreferred := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "routes_preferred",
			Help: "Total number of routes preferred by a route server for a given neighbor",
		},
		labels,
	)

	routesAccepted := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "routes_accepted",
			Help: "Total number of routes accepted by a route server for a given neighbor",
		},
		labels,
	)

	prometheus.MustRegister(neighborInfo)
	prometheus.MustRegister(neighborUptime)
	prometheus.MustRegister(routesReceived)
	prometheus.MustRegister(routesFiltered)
	prometheus.MustRegister(routesPreferred)
	prometheus.MustRegister(routesAccepted)

	return &metrics{
		neighborsStore: s,

		neighborInfo:   neighborInfo,
		neighborUptime: neighborUptime,

		routesReceived:  routesReceived,
		routesFiltered:  routesFiltered,
		routesPreferred: routesPreferred,
	}
}

// Update the metrics with information from the store.
func (m *metrics) update(ctx context.Context) error {
	rsIDs := m.neighborsStore.sources.GetSourceIDs()

	// For all route servers, fetch neighbors list and
	// update statistics.
	for _, rsID := range rsIDs {
		rs := m.neighborsStore.sources.Get(rsID)
		neighbors, err := m.neighborsStore.GetNeighborsAt(ctx, rsID)
		if err != nil {
			return err
		}

		// Get neighbors
		for _, neighbor := range neighbors {

			m.neighborInfo.With(prometheus.Labels{
				"route_server_id":      rs.ID,
				"route_server_name":    rs.Name,
				"route_server_group":   rs.Group,
				"neighbor_id":          neighbor.ID,
				"neighbor_description": neighbor.Description,
				"neighbor_asn":         strconv.Itoa(neighbor.ASN),
				"neighbor_address":     neighbor.Address,
				"neighbor_state":       neighbor.State,
			}).Set(1.0)

			labels := prometheus.Labels{
				"route_server_id":      rs.ID,
				"route_server_name":    rs.Name,
				"route_server_group":   rs.Group,
				"neighbor_id":          neighbor.ID,
				"neighbor_description": neighbor.Description,
				"neighbor_asn":         strconv.Itoa(neighbor.ASN),
				"neighbor_address":     neighbor.Address,
			}
			m.neighborUptime.
				With(labels).
				Set(neighbor.Uptime.Seconds())
			m.routesReceived.
				With(labels).
				Set(float64(neighbor.RoutesReceived))
			m.routesFiltered.
				With(labels).
				Set(float64(neighbor.RoutesFiltered))
			m.routesPreferred.
				With(labels).
				Set(float64(neighbor.RoutesPreferred))
			m.routesAccepted.
				With(labels).
				Set(float64(neighbor.RoutesAccepted))
		}
	}

	return nil
}

// Startmetrics registers the metrics and starts a
// periodical refresh.
func StartMetrics(
	ctx context.Context,
	neighborsStore *NeighborsStore,
) {
	m := initMetrics(neighborsStore)

	// Every 5 second, update the metrics
	log.Println(
		"[metrics] Starting refresh.")
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 5):
				err := m.update(ctx)
				if err != nil {
					log.Println(
						"[metrics] Error while updating:", err)
				}
			}
		}
	}()
}
