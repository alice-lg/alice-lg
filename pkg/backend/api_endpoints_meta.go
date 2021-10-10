package backend

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Handle Status Endpoint, this is intended for
// monitoring and service health checks
func apiStatusShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	status, err := NewAppStatus()
	return status, err
}

// Handle status
func apiStatus(_req *http.Request, params httprouter.Params) (api.Response, error) {
	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	source := AliceConfig.SourceInstanceByID(rsID)
	if source == nil {
		return nil, SOURCE_NOT_FOUND_ERROR
	}

	result, err := source.Status()
	if err != nil {
		apiLogSourceError("status", rsID, err)
	}

	return result, err
}

// Handle Config Endpoint
func apiConfigShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	result := api.ConfigResponse{
		Asn:            AliceConfig.Server.Asn,
		BgpCommunities: AliceConfig.UI.BgpCommunities,
		RejectReasons:  AliceConfig.UI.RoutesRejections.Reasons,
		Noexport: api.Noexport{
			LoadOnDemand: AliceConfig.UI.RoutesNoexports.LoadOnDemand,
		},
		NoexportReasons: AliceConfig.UI.RoutesNoexports.Reasons,
		RejectCandidates: api.RejectCandidates{
			Communities: AliceConfig.UI.RoutesRejectCandidates.Communities,
		},
		Rpki:                   api.Rpki(AliceConfig.UI.Rpki),
		RoutesColumns:          AliceConfig.UI.RoutesColumns,
		RoutesColumnsOrder:     AliceConfig.UI.RoutesColumnsOrder,
		NeighboursColumns:      AliceConfig.UI.NeighboursColumns,
		NeighboursColumnsOrder: AliceConfig.UI.NeighboursColumnsOrder,
		LookupColumns:          AliceConfig.UI.LookupColumns,
		LookupColumnsOrder:     AliceConfig.UI.LookupColumnsOrder,
		PrefixLookupEnabled:    AliceConfig.Server.EnablePrefixLookup,
	}
	return result, nil
}
