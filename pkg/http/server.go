package http

import (
	"log"
	"net/http"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store"
	"github.com/julienschmidt/httprouter"
)

// Server provides the HTTP server for the API
// and the assets.
type Server struct {
	cfg            *config.Config
	routesStore    *store.RoutesStore
	neighborsStore *store.NeighborsStore
}

// NewServer creates a new server
func NewServer(
	cfg *config.Config,
	routesStore *store.RoutesStore,
	neighborsStore *store.NeighborsStore,
) *Server {
	return &Server{
		cfg:            cfg,
		routesStore:    routesStore,
		neighborsStore: neighborsStore,
	}
}

// Start starts a HTTP server and begins to listen
// on the configured port.
func (s *Server) Start() {

	// Setup request routing
	router := httprouter.New()

	// Serve static content
	if err := webRegisterAssets(cfg, router); err != nil {
		log.Fatal(err)
	}

	if err := apiRegisterEndpoints(cfg, router); err != nil {
		log.Fatal(err)
	}

	// Start http server
	log.Fatal(http.ListenAndServe(cfg.Server.Listen, router))
}
