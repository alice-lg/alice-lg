package http

import (
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Handle get neighbors on routeserver
func (s *Server) apiNeighborsList(
	_req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	rsID, err := validateSourceID(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	var neighborsResponse *api.NeighborsResponse

	// Try to fetch neighbors from store, only fall back
	// to RS query if store is not ready yet
	sourceStatus := s.neighborsStore.SourceStatus(rsID)
	if sourceStatus.State == STATE_READY {
		neighbors := s.neighborsStore.GetNeighborsAt(rsID)
		// Make response
		neighborsResponse = &api.NeighborsResponse{
			Api: api.ApiStatus{
				Version: Version,
				CacheStatus: api.CacheStatus{
					OrigTtl:  0,
					CachedAt: sourceStatus.LastRefresh,
				},
				ResultFromCache: true, // you bet!
				Ttl: sourceStatus.LastRefresh.Add(
					s.neighborsStore.refreshInterval),
			},
			Neighbors: neighbors,
		}
	} else {
		source := s.neighborsStore.SourceInstanceByID(rsID)
		if source == nil {
			return nil, SOURCE_NOT_FOUND_ERROR
		}
		neighborsResponse, err = source.Neighbors()
		if err != nil {
			apiLogSourceError("neighbors", rsID, err)
			return nil, err
		}
	}

	// Sort result
	sort.Sort(&neighborsResponse.Neighbors)
	return neighborsResponse, nil
}
