/*
 * Prefix Lookup Reducer
 */

import {LOAD_RESULTS_REQUEST,
        LOAD_RESULTS_SUCCESS,
        LOAD_RESULTS_ERROR}
 from './actions'

const initialState = {
  query: '',

  results: [],
  error: null,
  queryDurationMs: 0.0,

  isLoading: false
}

export default function reducer(state=initialState, action) {
  switch(action.type) {
    case LOAD_RESULTS_REQUEST:
      return Object.assign({}, state, initialState, {
        query: action.payload.query,
        isLoading: true
      });
    case LOAD_RESULTS_SUCCESS:
      return Object.assign({}, state, {
        isLoading: false,
        query: action.payload.query,
        queryDurationMs: action.payload.results.query_duration_ms,
        results: action.payload.results.routes,
        error: null,
      });
    case LOAD_RESULTS_ERROR:
      return Object.assign({}, state, initialState, {
        query: action.payload.query,
        error: action.payload.error
      });
  }
  return state;
}


