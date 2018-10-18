
/*
 * Prefix lookup actions
 */

import axios from 'axios'

import {filtersUrlEncode} from './filter-encoding'

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

export function loadResults(query, filters, pageImported=0, pageFiltered=0) {
  return (dispatch) => {
    dispatch(loadResultsRequest(query));

    // Build querystring
    const q = `q=${query}&page_filtered=${pageFiltered}&page_imported=${pageImported}`;
    const f = filtersUrlEncode(filters); 
    axios.get(`/api/v1/lookup/prefix?${q}${f}`)
      .then(
        (res) => {
          dispatch(loadResultsSuccess(query, res.data));
        },
        (error) => {
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


