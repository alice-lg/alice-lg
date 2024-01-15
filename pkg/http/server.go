// Package http provides the server and API implementation
// for the webclient. The webclient's static files are also served.
package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

// Server provides the HTTP server for the API
// and the assets.
type Server struct {
	*http.Server
	cfg            *config.Config
	routesStore    *store.RoutesStore
	neighborsStore *store.NeighborsStore
	pool           *pgxpool.Pool
}

// NewServer creates a new server
func NewServer(
	cfg *config.Config,
	pool *pgxpool.Pool,
	routesStore *store.RoutesStore,
	neighborsStore *store.NeighborsStore,
) *Server {
	return &Server{
		cfg:            cfg,
		routesStore:    routesStore,
		neighborsStore: neighborsStore,
		pool:           pool,
	}
}

// Start starts a HTTP server and begins to listen
// on the configured port.
func (s *Server) Start(ctx context.Context) {
	router := httprouter.New()

	// Register routes
	if err := s.webRegisterAssets(ctx, router); err != nil {
		log.Fatal(err)
	}
	if err := s.apiRegisterEndpoints(router); err != nil {
		log.Fatal(err)
	}

	httpTimeout := time.Duration(s.cfg.Server.HTTPTimeout) * time.Second
	log.Println("Web server HTTP timeout set to:", httpTimeout)
	log.Println("Listening on:", s.cfg.Server.Listen)

	if s.cfg.Server.EnablePrefixLookup {
		log.Println("Prefix Lookup (Search): enabled")
		log.Println("Prefix Lookup Community Filter Cutoff:",
			s.cfg.Server.PrefixLookupCommunityFilterCutoff)
	} else {
		log.Println("Prefix Lookup (Search): disabled")
	}

	s.Server = &http.Server{
		Addr:         s.cfg.Server.Listen,
		Handler:      router,
		ReadTimeout:  httpTimeout,
		WriteTimeout: httpTimeout,
		IdleTimeout:  httpTimeout,
	}

	// Start http server
	log.Fatal(s.ListenAndServe())
}
