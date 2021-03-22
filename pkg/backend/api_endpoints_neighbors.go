package backend

import (
	"github.com/alice-lg/alice-lg/backend/api"
	"github.com/julienschmidt/httprouter"

	"net/http"
	"sort"
)

// Handle get neighbors on routeserver
func apiNeighborsList(
	_req *http.Request,
	params httprouter.Params,
) (api.Response, error) {
	rsId, err := validateSourceId(params.ByName("id"))
	if err != nil {
		return nil, err
	}

	var neighborsResponse *api.NeighboursResponse

	// Try to fetch neighbors from store, only fall back
	// to RS query if store is not ready yet
	sourceStatus := AliceNeighboursStore.SourceStatus(rsId)
	if sourceStatus.State == STATE_READY {
		neighbors := AliceNeighboursStore.GetNeighborsAt(rsId)
		// Make response
		neighborsResponse = &api.NeighboursResponse{
			Api: api.ApiStatus{
				Version: version,
				CacheStatus: api.CacheStatus{
					OrigTtl:  0,
					CachedAt: sourceStatus.LastRefresh,
				},
				ResultFromCache: true, // you bet!
				Ttl: sourceStatus.LastRefresh.Add(
					AliceNeighboursStore.refreshInterval),
			},
			Neighbours: neighbors,
		}
	} else {
		source := AliceConfig.SourceInstanceById(rsId)
		if source == nil {
			return nil, SOURCE_NOT_FOUND_ERROR
		}

		neighborsResponse, err = source.Neighbours()
		if err != nil {
			apiLogSourceError("neighbors", rsId, err)
			return nil, err
		}
	}

	// Sort result
	sort.Sort(&neighborsResponse.Neighbours)

	return neighborsResponse, nil
}
