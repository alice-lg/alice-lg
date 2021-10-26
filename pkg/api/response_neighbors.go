// Package api contains the datastructures for the
// Alice API.
package api

import (
	"encoding/json"
	"strings"
	"time"
)

// Neighbor is a BGP peer on the RS
type Neighbor struct {
	ID string `json:"id"`

	// Mandatory fields
	Address         string        `json:"address"`
	ASN             int           `json:"asn"`
	State           string        `json:"state"`
	Description     string        `json:"description"`
	RoutesReceived  int           `json:"routes_received"`
	RoutesFiltered  int           `json:"routes_filtered"`
	RoutesExported  int           `json:"routes_exported"`
	RoutesPreferred int           `json:"routes_preferred"`
	RoutesAccepted  int           `json:"routes_accepted"`
	Uptime          time.Duration `json:"uptime"`
	LastError       string        `json:"last_error"`
	RouteServerID   string        `json:"routeserver_id"`

	// Original response
	Details map[string]interface{} `json:"details"`
}

// String encodes a neighbor as json. This is
// more readable than the golang default represenation.
func (n *Neighbor) String() string {
	repr, _ := json.Marshal(n)
	return string(repr)
}

// Neighbors is a collection of neighbors
type Neighbors []*Neighbor

func (neighbors Neighbors) Len() int {
	return len(neighbors)
}

func (neighbors Neighbors) Less(i, j int) bool {
	return neighbors[i].ASN < neighbors[j].ASN
}

func (neighbors Neighbors) Swap(i, j int) {
	neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
}

// MatchSourceID implements Filterable interface
func (n *Neighbor) MatchSourceID(id string) bool {
	return n.RouteServerID == id
}

// MatchASN compares the neighbor's ASN.
func (n *Neighbor) MatchASN(asn int) bool {
	return n.ASN == asn
}

// MatchCommunity is undefined for neighbors.
func (n *Neighbor) MatchCommunity(Community) bool {
	return true // Ignore
}

// MatchExtCommunity is undefined for neighbors.
func (n *Neighbor) MatchExtCommunity(Community) bool {
	return true // Ignore
}

// MatchLargeCommunity is undefined for neighbors.
func (n *Neighbor) MatchLargeCommunity(Community) bool {
	return true // Ignore
}

// MatchName is a case insensitive match of
// the neighbor's description
func (n *Neighbor) MatchName(name string) bool {
	name = strings.ToLower(name)
	neighName := strings.ToLower(n.Description)

	return strings.Contains(neighName, name)
}

// A NeighborsResponse is a list of neighbors with
// caching information.
type NeighborsResponse struct {
	BackendResponse
	Neighbors Neighbors `json:"neighbors"`
}

// CacheTTL returns the duration of validity
// of the neighbor response.
func (res *NeighborsResponse) CacheTTL() time.Duration {
	now := time.Now().UTC()
	return res.Meta.TTL.Sub(now)
}

// NeighborsLookupResults is a mapping of lookup neighbors.
// The sourceID is used as a key.
type NeighborsLookupResults map[string]Neighbors

// NeighborStatus contains only the neighbor state and
// uptime.
type NeighborStatus struct {
	ID    string        `json:"id"`
	State string        `json:"state"`
	Since time.Duration `json:"uptime"`
}

// NeighborsStatus is a list of statuses.
type NeighborsStatus []*NeighborStatus

func (neighbors NeighborsStatus) Len() int {
	return len(neighbors)
}

func (neighbors NeighborsStatus) Less(i, j int) bool {
	return neighbors[i].ID < neighbors[j].ID
}

func (neighbors NeighborsStatus) Swap(i, j int) {
	neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
}

// NeighborsStatusResponse contains the status of all neighbors
// on a RS.
type NeighborsStatusResponse struct {
	BackendResponse
	Neighbors NeighborsStatus `json:"neighbors"`
}
