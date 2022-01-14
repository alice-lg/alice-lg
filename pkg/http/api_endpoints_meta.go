package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Handle Status Endpoint, this is intended for
// monitoring and service health checks
func (s *Server) apiStatusShow(
	req *http.Request,
	_params httprouter.Params,
) (response, error) {
	ctx := req.Context()
	status, err := CollectAppStatus(ctx, s.pool, s.routesStore, s.neighborsStore)
	return status, err
}

// Handle status
func (s *Server) apiStatus(
	_req *http.Request,
	params httprouter.Params,
) (response, error) {
	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	source := s.cfg.SourceInstanceByID(rsID)
	if source == nil {
		return nil, ErrSourceNotFound
	}

	result, err := source.Status()
	if err != nil {
		s.logSourceError("status", rsID, err)
	}

	return result, err
}

// Handle Config Endpoint
func (s *Server) apiConfigShow(
	_req *http.Request,
	_params httprouter.Params,
) (response, error) {
	result := api.ConfigResponse{
		Asn:            s.cfg.Server.Asn,
		BGPCommunities: s.cfg.UI.BGPCommunities,
		RejectReasons:  s.cfg.UI.RoutesRejections.Reasons,
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
