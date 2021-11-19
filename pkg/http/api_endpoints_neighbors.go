package http

import (
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/store"
)

// Handle get neighbors on routeserver
func (s *Server) apiNeighborsList(
	_req *http.Request,
	params httprouter.Params,
) (response, error) {
	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	var neighborsResponse *api.NeighborsResponse

	// Try to fetch neighbors from store, only fall back
	// to RS query if store is not ready yet
	sourceStatus, err := s.neighborsStore.SourceStatus(rsID)
	if err != nil {
		return nil, err
	}
	if sourceStatus.State == store.StateReady {
		neighbors := s.neighborsStore.GetNeighborsAt(rsID)
		// Make response
		neighborsResponse = &api.NeighborsResponse{
			Response: api.Response{
				Meta: &api.Meta{
					Version: config.Version,
					CacheStatus: api.CacheStatus{
						OrigTTL:  0,
						CachedAt: sourceStatus.LastRefresh,
					},
					ResultFromCache: true, // you bet!
					TTL: sourceStatus.LastRefresh.Add(
						s.neighborsStore.RefreshInterval),
				},
			},
			Neighbors: neighbors,
		}
	} else {
		source := s.cfg.SourceInstanceByID(rsID)
		if source == nil {
			return nil, ErrSourceNotFound
		}
		neighborsResponse, err = source.Neighbors()
		if err != nil {
			s.logSourceError("neighbors", rsID, err)
			return nil, err
		}
	}

	// Sort result
	sort.Sort(&neighborsResponse.Neighbors)
	return neighborsResponse, nil
}
