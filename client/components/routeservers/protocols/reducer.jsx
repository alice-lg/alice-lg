
/*
 * This will migrate to become the neighbors
 * reducer. Currently neihgbors are stored in
 * the routeserver reducer.
 */

import {SET_FILTER_VALUE} from './actions'

import {LOAD_ROUTESERVER_PROTOCOL_REQUEST,
        LOAD_ROUTESERVER_PROTOCOL_SUCCESS,
        LOAD_ROUTESERVER_PROTOCOL_ERROR}
  from '../actions'

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE';

const DEFAULT_SORT_COLUMN = "asn";
const DEFAULT_SORT_ORDER = "asc";

const initialState = {
  sortColumn: DEFAULT_SORT_COLUMN,
  sortOrder: DEFAULT_SORT_ORDER,

  isLoading: false,

  cachedAt: null,
  cacheTtl: null,

  filterQuery: "",
  filterValue: ""
};


// Reducer functions

function _handleLocationChange(state, payload) {
  const query = payload.query;
  const sort = query["s"] || DEFAULT_SORT_COLUMN;
  const order = query["o"]  || DEFAULT_SORT_ORDER; 
  const filterQuery = query["q"] || "";

  return Object.assign({}, state, {
    sortColumn: sort,
    sortOrder: order,
    
    filterQuery: filterQuery,
    filterValue: filterQuery
  });
}


export default function(state=initialState, action) {
  switch (action.type) {
    case LOCATION_CHANGE:
      return _handleLocationChange(state, action.payload);

    case SET_FILTER_VALUE:
      return Object.assign({}, state, {
        filterValue: action.payload.value
      });

    case LOAD_ROUTESERVER_PROTOCOL_REQUEST:
      return Object.assign({}, state, {
        isLoading: true,
      });

    case LOAD_ROUTESERVER_PROTOCOL_ERROR:
      return Object.assign({}, state, {
        isLoading: false,
      });

    // TODO: move neighbors list here
    case LOAD_ROUTESERVER_PROTOCOL_SUCCESS:
      return Object.assign({}, state, {
        isLoading: false,
        cachedAt: action.payload.api.cache_status.cached_at,
        cacheTtl: action.payload.api.ttl,
      });

    default:
  }

  return state;
}

