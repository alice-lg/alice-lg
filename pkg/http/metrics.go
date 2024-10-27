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

func (s *Server) registerMetrics(
	ctx context.Context,
	router *httprouter.Router,
) error {
	if s.cfg.Server.EnablePrometheus == false {
		return nil
	}

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

	s.routesStore.Stats(ctx)

	neighborState := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "neighbor_state",
			Help: "The state of a neighbor in a route server (0 = down, 1 = up)",
		},
		labels,
	)

	neighborUptime := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "neighbor_uptime",
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

	prometheus.MustRegister(neighborState)
	prometheus.MustRegister(neighborUptime)
	prometheus.MustRegister(routesReceived)
	prometheus.MustRegister(routesFiltered)
	prometheus.MustRegister(routesPreferred)
	prometheus.MustRegister(routesAccepted)

	log.Println("Prometheus metrics enabled, listening on /metrics")
	router.Handler("GET", "/metrics", promhttp.Handler())

	go func() {
		// Every 5 second, update the metrics with the latest data from the cache
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 5):
				routeServers, err := s.apiRouteServersList(ctx, nil, nil)
				if err != nil {
					log.Println("[metrics] error getting route servers: ", err)
					continue
				}

				for _, routeServer := range routeServers.(api.RouteServersResponse).RouteServers {
					neighbors, err := s.apiNeighborsList(ctx, nil, httprouter.Params{{Key: "id", Value: routeServer.ID}})
					if err != nil {
						log.Println("[metrics] error getting neighbors:", err)
						continue
					}

					for _, neighbor := range neighbors.(*api.NeighborsResponse).Neighbors {
						labels := prometheus.Labels{
							"route_server_id":      routeServer.ID,
							"route_server_name":    routeServer.Name,
							"route_server_group":   routeServer.Group,
							"neighbor_id":          neighbor.ID,
							"neighbor_description": neighbor.Description,
							"neighbor_asn":         strconv.Itoa(neighbor.ASN),
							"neighbor_address":     neighbor.Address,
						}

						state := float64(0)
						if neighbor.State == "up" {
							state = 1
						}

						neighborState.With(labels).Set(state)
						neighborUptime.With(labels).Set(neighbor.Uptime.Seconds())
						routesReceived.With(labels).Set(float64(neighbor.RoutesReceived))
						routesFiltered.With(labels).Set(float64(neighbor.RoutesFiltered))
						routesPreferred.With(labels).Set(float64(neighbor.RoutesPreferred))
						routesAccepted.With(labels).Set(float64(neighbor.RoutesAccepted))
					}
				}
			}

		}

	}()
	return nil
}
