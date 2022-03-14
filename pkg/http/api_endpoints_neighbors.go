package http

import (
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
)

// Handle get neighbors on routeserver
func (s *Server) apiNeighborsList(
	req *http.Request,
	params httprouter.Params,
) (response, error) {
	ctx := req.Context()
	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	var neighborsResponse *api.NeighborsResponse

	// Try to fetch neighbors from store, only fall back
	// to RS query if store is not ready yet.
	// The stored neighbors response includes details like
	// the number of filtered routes which might be lacking
	// from the summary.
	if s.neighborsStore.IsInitialized(rsID) {
		status, err := s.neighborsStore.GetStatus(rsID)
		if err != nil {
			return nil, err
		}
		neighbors, err := s.neighborsStore.GetNeighborsAt(ctx, rsID)
		if err != nil {
			return nil, err
		}
		// Make response
		neighborsResponse = &api.NeighborsResponse{
			Response: api.Response{
				Meta: &api.Meta{
					Version: config.Version,
					CacheStatus: api.CacheStatus{
						OrigTTL:  0,
						CachedAt: status.LastRefresh,
					},
					ResultFromCache: true, // you bet!
					TTL:             s.neighborsStore.SourceCacheTTL(ctx, rsID),
				},
			},
			Neighbors: neighbors,
		}
	} else {
		source := s.cfg.SourceInstanceByID(rsID)
		if source == nil {
			return nil, ErrSourceNotFound
		}
		neighborsResponse, err = source.NeighborsSummary()
		if err != nil {
			s.logSourceError("neighbors", rsID, err)
			return nil, err
		}
	}

	// Sort result
	sort.Sort(&neighborsResponse.Neighbors)
	return neighborsResponse, nil
}
