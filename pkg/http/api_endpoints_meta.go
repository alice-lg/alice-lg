package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
)

// Handle Status Endpoint, this is intended for
// monitoring and service health checks
func apiStatusShow(_req *http.Request, _params httprouter.Params) (api.Response, error) {
	status, err := CollectAppStatus()
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
func apiConfigShow(cfg *config.Config) apiEndpoint {
	return func(_req *http.Request, _params httprouter.Params) (api.Response, error) {
		result := api.ConfigResponse{
			Asn:            cfg.Server.Asn,
			BgpCommunities: cfg.UI.BgpCommunities,
			RejectReasons:  cdf.UI.RoutesRejections.Reasons,
			Noexport: api.Noexport{
				LoadOnDemand: cfg.UI.RoutesNoexports.LoadOnDemand,
			},
			NoexportReasons: cfg.UI.RoutesNoexports.Reasons,
			RejectCandidates: api.RejectCandidates{
				Communities: cfg.UI.RoutesRejectCandidates.Communities,
			},
			Rpki:                  api.Rpki(AliceConfig.UI.Rpki),
			RoutesColumns:         cfg.UI.RoutesColumns,
			RoutesColumnsOrder:    cfg.UI.RoutesColumnsOrder,
			NeighborsColumns:      cfg.UI.NeighborsColumns,
			NeighborsColumnsOrder: cfg.UI.NeighborsColumnsOrder,
			LookupColumns:         cfg.UI.LookupColumns,
			LookupColumnsOrder:    cfg.UI.LookupColumnsOrder,
			PrefixLookupEnabled:   cfg.Server.EnablePrefixLookup,
		}
		return result, nil
	}
}
