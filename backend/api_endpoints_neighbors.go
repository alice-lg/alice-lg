package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

// Handle get neighbors on routeserver
func apiNeighborsList(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Neighbours()
	if err != nil {
		apiLogSourceError("neighbors", rsId, err)
	}

	return result, err
}
