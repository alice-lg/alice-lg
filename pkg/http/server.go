package http

import (
	"log"
	"net/http"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store"
	"github.com/julienschmidt/httprouter"
)

// Server provides the HTTP server for the API
// and the assets.
type Server struct {
	*http.Server
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
	router := httprouter.New()

	// Register routes
	if err := s.webRegisterAssets(router); err != nil {
		log.Fatal(err)
	}
	if err := s.apiRegisterEndpoints(router); err != nil {
		log.Fatal(err)
	}
	httpTimeout := time.Duration(s.cfg.Server.HTTPTimeout) * time.Second
	log.Println("Web server HTTP timeout set to:", httpTimeout)

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
