/*
 * Prefix Lookup Reducer
 */

import {LOAD_RESULTS_REQUEST,
        LOAD_RESULTS_SUCCESS,
        LOAD_RESULTS_ERROR,
        
        SET_LOOKUP_QUERY_VALUE,

        RESET}
 from './actions'

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE'

const initialState = {
  query: "",
  queryValue: "",

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
 * Restore lookup query state from location paramenters
 */
const _restoreQueryState = function(state, payload) {
  const params = payload.query;
  const query = params["q"] || "";

  return Object.assign({}, state, {
    query: query,
    queryValue: query
  });
}

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
      return _restoreQueryState(state, action.payload);
      
    case SET_LOOKUP_QUERY_VALUE:
      return Object.assign({}, state, {
        queryValue: action.payload.value,
      });

    case LOAD_RESULTS_REQUEST:
      return Object.assign({}, state, initialState, {
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
        error: action.payload.error
      });

    case RESET:
      return Object.assign({}, state, initialState);
  }
  return state;
}


