package testdata

import (
	_ "embed" // testdata
	"encoding/json"
	"log"

	"github.com/alice-lg/alice-lg/pkg/api"
)

//go:embed routes_response.json
var testRoutesResponse []byte

// RoutesResponse returns the routes response from testdata
func RoutesResponse() *api.RoutesResponse {
	response := &api.RoutesResponse{}
	err := json.Unmarshal(testRoutesResponse, &response)
	if err != nil {
		log.Panic("could not unmarshal response test data:", err)
	}
	return response
}

// LoadTestLookupRoutes loads the testdata routes and converts
// them to lookup routes.
func LoadTestLookupRoutes(srcID, srcName string) api.LookupRoutes {
	res := RoutesResponse()

	// Prepare imported routes for lookup
	neighbors := map[string]*api.Neighbor{
		"ID163_AS31078": &api.Neighbor{
			ID: "ID163_AS31078",
		},
		"ID7254_AS31334": &api.Neighbor{
			ID: "ID7254_AS31334",
		},
	}
	rs := &api.RouteServer{
		ID:   srcID,
		Name: srcName,
	}
	imported := res.Imported.ToLookupRoutes("imported", rs, neighbors)
	filtered := res.Filtered.ToLookupRoutes("filtered", rs, neighbors)
	lookupRoutes := append(imported, filtered...)
	return lookupRoutes
}
