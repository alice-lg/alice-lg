/*
 * Prefix Lookup Reducer
 */

import {LOAD_RESULTS_REQUEST,
        LOAD_RESULTS_SUCCESS,
        LOAD_RESULTS_ERROR,
        
        SET_LOOKUP_QUERY_VALUE}
 from './actions'

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE'

const initialState = {
  query: "",
  queryValue: "",

  results: [],
  error: null,

  queryDurationMs: 0.0,

  limit: 100,
  offset: 0,
  totalRoutes: 0,

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
        isLoading: true
      });
    case LOAD_RESULTS_SUCCESS:
      if (state.query != action.payload.query) {
        return state;
      }

      return Object.assign({}, state, {
        isLoading: false,
        query: action.payload.query,
        queryDurationMs: action.payload.results.query_duration_ms,
        results: action.payload.results.routes,
        limit: action.payload.results.limit,
        offset: action.payload.results.offset,
        totalRoutes: action.payload.results.total_routes,
        error: null
      });

    case LOAD_RESULTS_ERROR:
      if (state.query != action.payload.query) {
        return state;
      }

      return Object.assign({}, state, initialState, {
        query: action.payload.query,
        error: action.payload.error
      });
  }
  return state;
}


