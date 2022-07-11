package store

import (
	"context"
	"sort"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store/backends/memory"
)

// Make a store and populate it with data
func makeTestNeighborsStore() *NeighborsStore {
	be := memory.NewNeighborsBackend()

	cfg := &config.Config{
		Server: config.ServerConfig{
			NeighborsStoreRefreshInterval: 1,
		},
		Sources: []*config.SourceConfig{
			{
				ID:   "rs1",
				Name: "rs1",
			},
			{
				ID:   "rs2",
				Name: "rs2",
			},
		},
	}

	rs1 := api.Neighbors{
		&api.Neighbor{
			ID:          "ID2233_AS2342",
			ASN:         2342,
			Description: "PEER AS2342 192.9.23.42 Customer Peer 1",
		},
		&api.Neighbor{
			ID:          "ID2233_AS2343",
			ASN:         2343,
			Description: "PEER AS2343 192.9.23.43 Different Peer 1",
		},
		&api.Neighbor{
			ID:          "ID2233_AS2344",
			ASN:         2344,
			Description: "PEER AS2344 192.9.23.44 3rd Peer from the sun",
		},
		&api.Neighbor{
			ID:          "ID163_AS31078",
			ASN:         31078,
			Description: "PEER AS31078 1.2.3.4 Peer Peer",
		},
		&api.Neighbor{
			ID:          "ID7254_AS31334",
			ASN:         31078,
			Description: "PEER AS31334 4.3.2.1 Peer",
		},
	}
	rs2 := api.Neighbors{
		&api.Neighbor{
			ID:          "ID2233_AS2342",
			ASN:         2342,
			Description: "PEER AS2342 192.9.23.42 Customer Peer 1",
		},
		&api.Neighbor{
			ID:          "ID2233_AS4223",
			ASN:         4223,
			Description: "PEER AS4223 192.9.42.23 Cloudfoo Inc.",
		},
	}

	be.SetNeighbors(context.Background(), "rs1", rs1)
	be.SetNeighbors(context.Background(), "rs2", rs2)

	// Create store
	store := NewNeighborsStore(cfg, be)
	return store
}

func TestGetNeighborsMapAt(t *testing.T) {
	store := makeTestNeighborsStore()

	neighbors, err := store.GetNeighborsMapAt(context.Background(), "rs1")
	if err != nil {
		t.Fatal(err)
	}
	neighbor := neighbors["ID2233_AS2343"]
	if neighbor.ID != "ID2233_AS2343" {
		t.Error("unexpected neighbor:", neighbor)
	}
}

func TestGetNeighbors(t *testing.T) {
	store := makeTestNeighborsStore()
	neighbors, err := store.GetNeighborsAt(context.Background(), "rs2")
	if err != nil {
		t.Fatal(err)
	}

	if len(neighbors) != 2 {
		t.Error("Expected 2 neighbors, got:", len(neighbors))
	}

	sort.Sort(neighbors)

	if neighbors[0].ID != "ID2233_AS2342" {
		t.Error("Expected neighbor: ID2233_AS2342, got:",
			neighbors[0])
	}

	neighbors, err = store.GetNeighborsAt(context.Background(), "rs3")
	if err == nil {
		t.Error("Unknown source should have yielded zero results")
	}
	t.Log(neighbors)

}

func TestNeighborLookup(t *testing.T) {
	store := makeTestNeighborsStore()

	results, err := store.LookupNeighbors(context.Background(), "Cloudfoo")
	if err != nil {
		t.Fatal(err)
	}

	// Peer should be present at RS2
	neighbors, ok := results["rs2"]
	if !ok {
		t.Error("Lookup on rs2 unsuccessful.")
	}

	if len(neighbors) > 1 {
		t.Error("Lookup should match exact 1 peer.")
	}

	n := neighbors[0]
	if n.ID != "ID2233_AS4223" {
		t.Error("Wrong peer in lookup response")
	}
}

func TestNeighborFilter(t *testing.T) {
	ctx := context.Background()
	store := makeTestNeighborsStore()
	filter := api.NeighborFilterFromQueryString("asn=2342")
	neighbors, err := store.FilterNeighbors(ctx, filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(neighbors) != 2 {
		t.Error("Expected two results")
	}

	filter = api.NeighborFilterFromQueryString("")
	neighbors, err = store.FilterNeighbors(ctx, filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(neighbors) != 0 {
		t.Error("Expected empty result set")
	}

}

func TestReMatchASLookup(t *testing.T) {
	if !ReMatchASLookup.MatchString("AS2342") {
		t.Error("should be ASN")
	}
	if ReMatchASLookup.MatchString("Goo") {
		t.Error("should not be ASN")
	}
}
