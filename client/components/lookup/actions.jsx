
/*
 * Prefix lookup actions
 */

import axios from 'axios'

export const LOAD_RESULTS_REQUEST = '@lookup/LOAD_RESULTS_REQUEST';
export const LOAD_RESULTS_SUCCESS = '@lookup/LOAD_RESULTS_SUCCESS';
export const LOAD_RESULTS_ERROR   = '@lookup/LOAD_RESULTS_ERROR';

// Action creators
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

export function loadResults(query, limit=50, offset=0) {
  return (dispatch) => {
    dispatch(loadResultsRequest(query));

    axios.get(`/api/lookup/prefix?q=${query}&limit=${limit}&offset=${offset}`)
      .then((res) => {
        dispatch(loadResultsSuccess(query, res.data));
      })
      .catch((error) => {
        dispatch(loadResultsError(query, error));
      });
  }
}


