package api

import (
	"encoding/json"
	"strings"
	"time"
)

// Neighbours
type Neighbours []*Neighbour

type Neighbour struct {
	Id string `json:"id"`

	// Mandatory fields
	Address         string        `json:"address"`
	Asn             int           `json:"asn"`
	State           string        `json:"state"`
	Description     string        `json:"description"`
	RoutesReceived  int           `json:"routes_received"`
	RoutesFiltered  int           `json:"routes_filtered"`
	RoutesExported  int           `json:"routes_exported"`
	RoutesPreferred int           `json:"routes_preferred"`
	RoutesAccepted  int           `json:"routes_accepted"`
	Uptime          time.Duration `json:"uptime"`
	LastError       string        `json:"last_error"`
	RouteServerId   string        `json:"routeserver_id"`

	// Original response
	Details map[string]interface{} `json:"details"`
}

// String encodes a neighbor as json. This is
// more readable than the golang default represenation.
func (n *Neighbour) String() string {
	repr, _ := json.Marshal(n)
	return string(repr)
}

// Implement sorting interface for routes
func (neighbours Neighbours) Len() int {
	return len(neighbours)
}

func (neighbours Neighbours) Less(i, j int) bool {
	return neighbours[i].Asn < neighbours[j].Asn
}

func (neighbours Neighbours) Swap(i, j int) {
	neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
}

type NeighboursResponse struct {
	Api        ApiStatus  `json:"api"`
	Neighbours Neighbours `json:"neighbours"`
}

// Implement Filterable interface
func (self *Neighbour) MatchSourceId(id string) bool {
	return self.RouteServerId == id
}

func (self *Neighbour) MatchAsn(asn int) bool {
	return self.Asn == asn
}

func (self *Neighbour) MatchCommunity(_community Community) bool {
	return true // Ignore
}

func (self *Neighbour) MatchExtCommunity(_community Community) bool {
	return true // Ignore
}

func (self *Neighbour) MatchLargeCommunity(_community Community) bool {
	return true // Ignore
}

func (self *Neighbour) MatchName(name string) bool {
	name = strings.ToLower(name)
	neighName := strings.ToLower(self.Description)

	return strings.Contains(neighName, name)
}

// Neighbours response is cacheable
func (self *NeighboursResponse) CacheTTL() time.Duration {
	now := time.Now().UTC()
	return self.Api.Ttl.Sub(now)
}

type NeighboursLookupResults map[string]Neighbours

type NeighboursStatus []*NeighbourStatus

type NeighbourStatus struct {
	Id    string        `json:"id"`
	State string        `json:"state"`
	Since time.Duration `json:"uptime"`
}

// Implement sorting interface for status
func (neighbours NeighboursStatus) Len() int {
	return len(neighbours)
}

func (neighbours NeighboursStatus) Less(i, j int) bool {
	return neighbours[i].Id < neighbours[j].Id
}

func (neighbours NeighboursStatus) Swap(i, j int) {
	neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
}

type NeighboursStatusResponse struct {
	Api        ApiStatus        `json:"api"`
	Neighbours NeighboursStatus `json:"neighbours"`
}
