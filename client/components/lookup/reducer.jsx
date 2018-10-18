/*
 * Prefix Lookup Reducer
 */

import {
  LOAD_RESULTS_REQUEST,
  LOAD_RESULTS_SUCCESS,
  LOAD_RESULTS_ERROR,

  SET_LOOKUP_QUERY_VALUE,

  RESET,
} from './actions'

import {
  APPLY_FILTER,
} from './filter-actions'

import {
  FILTER_GROUP_SOURCES,
  FILTER_GROUP_ASNS,
  FILTER_GROUP_COMMUNITIES,
  FILTER_GROUP_EXT_COMMUNITIES,
  FILTER_GROUP_LARGE_COMMUNITIES,
} from './filter-groups'

import {
  decodeFiltersSources,
  decodeFiltersAsns,
  decodeFiltersCommunities,
  decodeFiltersExtCommunities,
  decodeFiltersLargeCommunities,
} from './filter-encoding'


const LOCATION_CHANGE = '@@router/LOCATION_CHANGE'

const initialFilterState = [
  {"key": "sources", "filters": []}, 
  {"key": "asns", "filters": []}, 
  {"key": "communities", "filters": []}, 
  {"key": "ext_communities", "filters": []}, 
  {"key": "large_communities", "filters": []}, 
];

const initialState = {
  query: "",
  queryValue: "",

  anchor: "",

  filtersAvailable: initialFilterState,
  filtersApplied: initialFilterState,

  routesImported: [],
  routesFiltered: [],

  error: null,

  queryDurationMs: 0.0,

  cachedAt: false,
  cacheTtl: false,

  pageImported: 0,
  pageFiltered: 0,

  pageSizeImported: 0,
  pageSizeFiltered: 0,

  totalPagesImported: 0,
  totalPagesFiltered: 0,

  totalRoutesImported: 0,
  totalRoutesFiltered: 0,

  isLoading: false
}

/*
 * Helper: Get scroll anchor from hash
 */
const getScrollAnchor = function(hash) {
  return hash.substr(hash.indexOf('-')+1);
}

/*
 * Decode filters applied from params
 */
const _decodeFiltersApplied = function(params) {
  let groups = Object.assign([], initialFilterState);

  groups[FILTER_GROUP_SOURCES].filters =           decodeFiltersSources(params);
  groups[FILTER_GROUP_ASNS].filters =              decodeFiltersAsns(params);
  groups[FILTER_GROUP_COMMUNITIES].filters =       decodeFiltersCommunities(params);
  groups[FILTER_GROUP_EXT_COMMUNITIES].filters =   decodeFiltersExtCommunities(params);
  groups[FILTER_GROUP_LARGE_COMMUNITIES].filters = decodeFiltersLargeCommunities(params);

  return groups;
}


/*
 * Restore lookup query state from location paramenters
 */
const _handleLocationChange = function(state, payload) {
  const params = payload.query;
  const query = params["q"] || "";
  const pageFiltered = parseInt(params["pf"] || 0, 10);
  const pageReceived = parseInt(params["pr"] || 0, 10);
  const anchor = getScrollAnchor(payload.hash);

  // Restore filters applied from location
  const filtersApplied = _decodeFiltersApplied(params);

  return Object.assign({}, state, {
    anchor: anchor,
    query: query,
    queryValue: query,
    pageImported: pageReceived,
    pageFiltered: pageFiltered,
    filtersApplied: filtersApplied,
  });
}

/*
 * Receive query results
 */
const _loadQueryResult = function(state, payload) {
  const results = payload.results;
  const imported = results.imported;
  const filtered = results.filtered;
  const api = results.api;

  return Object.assign({}, state, {
    isLoading: false,

    // Cache Status
    cachedAt: api.cache_status.cached_at, // I don't like this style. 
    cacheTtl: api.ttl, 
    
    // Routes
    routesImported: imported.routes,
    routesFiltered: filtered.routes,

    // Filters available
    filtersAvailable: results.filters_available,

    // Pagination
    pageImported:        imported.pagination.page,
    pageFiltered:        filtered.pagination.page,
    pageSizeImported:    imported.pagination.page_size,
    pageSizeFiltered:    filtered.pagination.page_size,
    totalPagesImported:  imported.pagination.total_pages, 
    totalPagesFiltered:  filtered.pagination.total_pages,
    totalRoutesImported: imported.pagination.total_results,
    totalRoutesFiltered: filtered.pagination.total_results,

    // Statistics
    queryDurationMs: results.request_duration_ms,
    totalRoutes:     imported.pagination.total_results + filtered.pagination.total_results
  });
}


export default function reducer(state=initialState, action) {
  switch(action.type) {
    case LOCATION_CHANGE:
      return _handleLocationChange(state, action.payload);
      
    case SET_LOOKUP_QUERY_VALUE:
      return Object.assign({}, state, {
        queryValue: action.payload.value,
      });

    case LOAD_RESULTS_REQUEST:
      return Object.assign({}, state, {
        query: action.payload.query,
        queryValue: action.payload.query,
        isLoading: true
      });
    case LOAD_RESULTS_SUCCESS:
      if (state.query != action.payload.query) {
        return state;
      }
      return _loadQueryResult(state, action.payload);

    case LOAD_RESULTS_ERROR:
      if (state.query != action.payload.query) {
        return state;
      }

      return Object.assign({}, state, initialState, {
        query: action.payload.query,
        queryValue: action.payload.query,
        error: action.payload.error
      });

    case RESET:
      return Object.assign({}, state, initialState);
  }
  return state;
}


