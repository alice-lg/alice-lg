package http

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// All metrics are grouped here.
type metrics struct {
	server *Server

	neighborInfo   *prometheus.GaugeVec
	neighborUptime *prometheus.GaugeVec

	routesReceived  *prometheus.GaugeVec
	routesFiltered  *prometheus.GaugeVec
	routesPreferred *prometheus.GaugeVec
	routesAccepted  *prometheus.GaugeVec
}

// Create and initialize metrics
func initMetrics(server *Server) *metrics {
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
		server: server,

		neighborInfo:   neighborInfo,
		neighborUptime: neighborUptime,

		routesReceived:  routesReceived,
		routesFiltered:  routesFiltered,
		routesPreferred: routesPreferred,
	}
}

// Update metrics
func (m *metrics) update(ctx context.Context) error {
	res, err := m.server.apiRouteServersList(ctx, nil, nil)
	if err != nil {
		return err
	}
	routeServers := res.(api.RouteServersResponse).RouteServers

	// For all route servers, fetch neighbors list and
	// update statistics.
	for _, rs := range routeServers {
		res, err = m.server.apiNeighborsList(
			ctx, nil, httprouter.Params{
				{Key: "id", Value: rs.ID},
			})
		if err != nil {
			return err
		}
		neighbors := res.(*api.NeighborsResponse).Neighbors

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

func (s *Server) registerMetrics(
	ctx context.Context,
	router *httprouter.Router,
) error {
	if s.cfg.Server.EnableMetrics == false {
		return nil
	}

	m := initMetrics(s)

	log.Println("Metrics enabled and available on: /metrics")
	router.Handler("GET", "/metrics", promhttp.Handler())

	// Every 5 second, update the metrics
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

	return nil
}
