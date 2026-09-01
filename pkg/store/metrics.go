package store

import (
	"context"
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	neighborsStore *NeighborsStore
	ctx            context.Context

	neighborInfo   *prometheus.Desc
	neighborUptime *prometheus.Desc

	routesReceived  *prometheus.Desc
	routesFiltered  *prometheus.Desc
	routesPreferred *prometheus.Desc
	routesAccepted  *prometheus.Desc
}

// StartMetrics initializes a metrics object and registers it as a prometheus.Collector
func StartMetrics(ctx context.Context, s *NeighborsStore) {
	log.Println("[metrics] Initializing export.")

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

	neighborInfo := prometheus.NewDesc(
		"neighbor_info",
		"Information about the neighbor including the state",
		append(labels, "neighbor_state"), nil,
	)

	neighborUptime := prometheus.NewDesc(
		"neighbor_uptime_seconds_total",
		"The uptime of a neighbor on a route server in seconds",
		labels, nil,
	)

	routesReceived := prometheus.NewDesc(
		"routes_received",
		"Total number of routes received by a route server for a given neighbor",
		labels, nil,
	)

	routesFiltered := prometheus.NewDesc(
		"routes_filtered",
		"Total number of routes filtered by a route server for a given neighbor",
		labels, nil,
	)

	routesPreferred := prometheus.NewDesc(
		"routes_preferred",
		"Total number of routes preferred by a route server for a given neighbor",
		labels, nil,
	)

	routesAccepted := prometheus.NewDesc(
		"routes_accepted",
		"Total number of routes accepted by a route server for a given neighbor",
		labels, nil,
	)

	m := &metrics{
		neighborsStore: s,
		ctx:            ctx,

		neighborInfo:   neighborInfo,
		neighborUptime: neighborUptime,

		routesReceived:  routesReceived,
		routesFiltered:  routesFiltered,
		routesPreferred: routesPreferred,
		routesAccepted:  routesAccepted,
	}

	prometheus.MustRegister(m)
}

// Describe implements prometheus.Collector and passes the descriptions of all the metrics we'll report
func (m *metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.neighborInfo
	ch <- m.neighborUptime
	ch <- m.routesReceived
	ch <- m.routesFiltered
	ch <- m.routesPreferred
	ch <- m.routesAccepted
}

// Collect implements prometheus.Collector and returns the current state of all metrics so they can be scraped
func (m *metrics) Collect(ch chan<- prometheus.Metric) {
	rsIDs := m.neighborsStore.sources.GetSourceIDs()

	// For all route servers, fetch neighbors list and
	// update statistics.
	for _, rsID := range rsIDs {
		if !m.neighborsStore.IsInitialized(rsID) {
			continue // No data from RS yet
		}
		rs := m.neighborsStore.sources.Get(rsID)

		neighbors, err := m.neighborsStore.GetNeighborsAt(m.ctx, rsID)
		if err != nil {
			continue
		}

		// Get neighbors
		for _, neighbor := range neighbors {
			// label values, in the same order as defined in StartMetrics
			labels := []string{
				rs.ID,                      // route_server_id
				rs.Name,                    // route_server_name
				rs.Group,                   // route_server_group
				neighbor.ID,                // neighbor_id
				neighbor.Description,       // neighbor_description
				strconv.Itoa(neighbor.ASN), // neighbor_asn
				neighbor.Address,           // neighbor_address
			}

			ch <- prometheus.MustNewConstMetric(
				m.neighborInfo,
				prometheus.GaugeValue,
				1.0,
				append(labels, neighbor.State)..., // neighbor_info has the additional label neighbor_state
			)
			ch <- prometheus.MustNewConstMetric(
				m.neighborUptime,
				prometheus.GaugeValue,
				neighbor.Uptime.Seconds(),
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				m.routesReceived,
				prometheus.GaugeValue,
				float64(neighbor.RoutesReceived),
				labels...,
			)
			ch <- prometheus.MustNewConstMetric(
				m.routesFiltered,
				prometheus.GaugeValue,
				float64(neighbor.RoutesFiltered),
				labels...,
			)
			ch <- prometheus.MustNewConstMetric(
				m.routesPreferred,
				prometheus.GaugeValue,
				float64(neighbor.RoutesPreferred),
				labels...,
			)
			ch <- prometheus.MustNewConstMetric(
				m.routesAccepted,
				prometheus.GaugeValue,
				float64(neighbor.RoutesAccepted),
				labels...,
			)
		}
	}
}
