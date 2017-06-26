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
        isLoading: true,
      });
    case LOAD_RESULTS_SUCCESS:
      return Object.assign({}, state, {
        isLoading: false,
        queryDurationMs: action.payload.results.query_duration_ms,
        results: action.payload.results.routes,
        error: null,
      });
    case LOAD_RESULTS_ERROR:
      return Object.assign({}, state, initialState, {
        error: action.payload.error,
      });
  }
  return state;
}


