package main

import (
	"github.com/ecix/alice-lg/backend/api"
	"github.com/ecix/alice-lg/backend/sources"

	"log"
	"time"
)

type RoutesStore struct {
	routesMap map[sources.Source]api.RoutesResponse
	ttlMap    map[string]time.Time
}

func NewRoutesStore(config *Config) *RoutesStore {

	// Build mapping based on source instances
	routesMap := make(map[sources.Source]api.RoutesResponse)
	for _, source := range config.Sources {
		routesMap[source.getInstance()] = api.RoutesResponse{}
	}

	store := &RoutesStore{
		routesMap: routesMap,
		ttlMap:    make(map[string]time.Time),
	}
	return store
}

func (self *RoutesStore) Start() {
	log.Println("Starting local routes store")
}

func (self *RoutesStore) init() {

}

// Update all routes
func (self *RoutesStore) update() {
	for source, _ := range self.routesMap {
		routes, err := source.AllRoutes()
		if err != nil {
			log.Println("Error while updating routes cache:", err)
			continue
		}
		self.routesMap[source] = routes
	}
}

func (self *RoutesStore) Lookup(prefix string) []api.LookupRoute {
	result := []api.LookupRoute{}

	return result
}
