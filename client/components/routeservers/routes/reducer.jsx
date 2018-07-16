

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE'

import {FETCH_ROUTES_RECEIVED_REQUEST,
        FETCH_ROUTES_RECEIVED_SUCCESS,
        FETCH_ROUTES_RECEIVED_ERROR,

        FETCH_ROUTES_FILTERED_REQUEST,
        FETCH_ROUTES_FILTERED_SUCCESS,
        FETCH_ROUTES_FILTERED_ERROR,

        FETCH_ROUTES_NOT_EXPORTED_REQUEST,
        FETCH_ROUTES_NOT_EXPORTED_SUCCESS,
        FETCH_ROUTES_NOT_EXPORTED_ERROR} from './actions'


const initialState = {

  received: [],
  receivedLoading: false,
  receivedPage: 0,
  receivedTotalPages: 0,
  receivedTotalResults: 0,
  receivedApiStatus: {},

  filtered: [],
  filteredLoading: false,
  filteredPage: 0,
  filteredTotalPages: 0,
  filteredTotalResults: 0,
  filteredApiStatus: {},

  notExported: [],
  notExportedLoading: false,
  notExportedPage: 0,
  notExportedTotalPages: 0,
  notExportedTotalResults: 0,
  notExportedApiStatus: {},

  filterQuery: "",
}


// Handlers:
function _handleLocationChange(state, payload) {
  // Check query payload
  let query = payload.query;

  let filterQuery = query["q"] || "";

  let receivedPage    = query["pr"] || 0;
  let filteredPage    = query["pf"] || 0;
  let notExportedPage = query["pn"] || 0;

  // Assert numeric
  receivedPage    = parseInt(receivedPage);
  filteredPage    = parseInt(filteredPage);
  notExportedPage = parseInt(notExportedPage);

  let nextState = Object.assign({}, state, {
    filterQuery: filterQuery,

    receivedPage:    receivedPage,
    filteredPage:    filteredPage,
    notExportedPage: notExportedPage,
  });

  return nextState;
}



export default function reducer(state=initialState, action) {

  switch(action.type) {
    case LOCATION_CHANGE:
      return _handleLocationChange(state, action.payload);
  }

  return state;
}


