
/*
 * Prefix lookup actions
 */

import axios from 'axios'

export const SET_LOOKUP_QUERY_VALUE = '@lookup/SET_LOOKUP_QUERY_VALUE';

export const LOAD_RESULTS_REQUEST = '@lookup/LOAD_RESULTS_REQUEST';
export const LOAD_RESULTS_SUCCESS = '@lookup/LOAD_RESULTS_SUCCESS';
export const LOAD_RESULTS_ERROR   = '@lookup/LOAD_RESULTS_ERROR';

export const RESET = "@lookup/RESET";

// Action creators
export function setLookupQueryValue(value) {
  return {
    type: SET_LOOKUP_QUERY_VALUE,
    payload: {
      value: value,
    }
  }
}


export function loadResultsRequest(query) {
  return {
    type: LOAD_RESULTS_REQUEST,
    payload: {
      query: query
    }
  }
}

export function loadResultsSuccess(query, results) {
  return {
    type: LOAD_RESULTS_SUCCESS,
    payload: {
      query: query,
      results: results
    }
  }
}

export function loadResultsError(query, error) {
  return {
    type: LOAD_RESULTS_ERROR,
    payload: {
      query: query,
      error: error
    }
  }
}

export function loadResults(query, pageImported=0, pageFiltered=0) {
  return (dispatch) => {
    dispatch(loadResultsRequest(query));

    // Build querystring
    let q = `q=${query}&page_filtered=${pageFiltered}&page_imported=${pageImported}`;
    axios.get(`/api/v1/lookup/prefix?${q}`)
      .then((res) => {
        dispatch(loadResultsSuccess(query, res.data));
      })
      .catch((error) => {
        dispatch(loadResultsError(query, error));
      });
  }
}

export function reset() {
  return {
    type: RESET,
    payload: {}
  }
}


