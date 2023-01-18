package http

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Handle Status Endpoint, this is intended for
// monitoring and service health checks
func (s *Server) apiStatusShow(
	ctx context.Context,
	req *http.Request,
	_params httprouter.Params,
) (response, error) {
	status, err := CollectAppStatus(ctx, s.pool, s.routesStore, s.neighborsStore)
	return status, err
}

// Handle Config Endpoint
func (s *Server) apiConfigShow(
	_ctx context.Context,
	_req *http.Request,
	_params httprouter.Params,
) (response, error) {
	result := api.ConfigResponse{
		Asn:                     s.cfg.Server.Asn,
		BGPCommunities:          s.cfg.UI.BGPCommunities,
		BGPBlackholeCommunities: s.cfg.UI.BGPBlackholeCommunities,
		RejectReasons:           s.cfg.UI.RoutesRejections.Reasons,
		Noexport: api.Noexport{
			LoadOnDemand: s.cfg.UI.RoutesNoexports.LoadOnDemand,
		},
		NoexportReasons: s.cfg.UI.RoutesNoexports.Reasons,
		RejectCandidates: api.RejectCandidates{
			Communities: s.cfg.UI.RoutesRejectCandidates.Communities,
		},
		Rpki:                  api.Rpki(s.cfg.UI.Rpki),
		RoutesColumns:         s.cfg.UI.RoutesColumns,
		RoutesColumnsOrder:    s.cfg.UI.RoutesColumnsOrder,
		NeighborsColumns:      s.cfg.UI.NeighborsColumns,
		NeighborsColumnsOrder: s.cfg.UI.NeighborsColumnsOrder,
		LookupColumns:         s.cfg.UI.LookupColumns,
		LookupColumnsOrder:    s.cfg.UI.LookupColumnsOrder,
		PrefixLookupEnabled:   s.cfg.Server.EnablePrefixLookup,
	}
	return result, nil
}
