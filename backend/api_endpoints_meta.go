package main

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

// Handle Status Endpoint, this is intended for
// monitoring and service health checks
func apiStatusShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	status, err := NewAppStatus()
	return status, err
}

// Handle status
func apiStatus(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}
	source := AliceConfig.Sources[rsId].getInstance()
	result, err := source.Status()
	if err != nil {
		apiLogSourceError("status", rsId, err)
	}

	return result, err
}

// Handle Config Endpoint
func apiConfigShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	result := api.ConfigResponse{
		Rejection: api.Rejection{
			Asn:      AliceConfig.Ui.RoutesRejections.Asn,
			RejectId: AliceConfig.Ui.RoutesRejections.RejectId,
		},
		RejectReasons: SerializeReasons(
			AliceConfig.Ui.RoutesRejections.Reasons),
		Noexport: api.Noexport{
			Asn:          AliceConfig.Ui.RoutesNoexports.Asn,
			NoexportId:   AliceConfig.Ui.RoutesNoexports.NoexportId,
			LoadOnDemand: AliceConfig.Ui.RoutesNoexports.LoadOnDemand,
		},
		Rpki:           api.Rpki(AliceConfig.Ui.Rpki),
		BgpCommunities: AliceConfig.Ui.BgpCommunities,
		NoexportReasons: SerializeReasons(
			AliceConfig.Ui.RoutesNoexports.Reasons),
		RoutesColumns:          AliceConfig.Ui.RoutesColumns,
		RoutesColumnsOrder:     AliceConfig.Ui.RoutesColumnsOrder,
		NeighboursColumns:      AliceConfig.Ui.NeighboursColumns,
		NeighboursColumnsOrder: AliceConfig.Ui.NeighboursColumnsOrder,
		LookupColumns:          AliceConfig.Ui.LookupColumns,
		LookupColumnsOrder:     AliceConfig.Ui.LookupColumnsOrder,
		PrefixLookupEnabled:    AliceConfig.Server.EnablePrefixLookup,
	}
	return result, nil
}
