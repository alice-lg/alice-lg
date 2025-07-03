package http

import (
	"context"
	"log"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) registerMetrics(
	ctx context.Context,
	router *httprouter.Router,
) error {
	if !s.cfg.Server.EnableMetrics {
		return nil
	}

	log.Println("Metrics enabled and available on: /metrics")
	router.Handler("GET", "/metrics", promhttp.Handler())

	return nil
}
