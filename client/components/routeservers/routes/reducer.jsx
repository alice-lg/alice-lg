

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

import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';

const initialState = {

  received: [],
  receivedLoading: false,
  receivedError: null,
  receivedPage: 0,
  receivedPageSize: 0,
  receivedTotalPages: 0,
  receivedTotalResults: 0,
  receivedApiStatus: {},

  filtered: [],
  filteredLoading: false,
  filteredError: null,
  filteredPage: 0,
  filteredPageSize: 0,
  filteredTotalPages: 0,
  filteredTotalResults: 0,
  filteredApiStatus: {},

  notExported: [],
  notExportedLoading: false,
  notExportedError: null,
  notExportedPage: 0,
  notExportedPageSize: 0,
  notExportedTotalPages: 0,
  notExportedTotalResults: 0,
  notExportedApiStatus: {},

  filterQuery: "",
}


// Helpers
function _stateType(type) {
  let stype = type;
  if (stype == ROUTES_NOT_EXPORTED) {
    stype = "notExported"; // TODO: This lacks elegance.
  }
  return stype;
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

function _handleFetchRoutesRequest(type, state, payload) {
  const stype = _stateType(type);
  let nextState = Object.assign({}, state, {
    [stype+'Loading']: true,
  });

  return nextState;
}


function _handleFetchRoutesSuccess(type, state, payload) {
  const stype = _stateType(type);
  const pagination = payload.pagination;
  const apiStatus = payload.api;

  let nextState = Object.assign({}, state, {
    [stype]: payload.routes,

    [stype+'Page']:         pagination.page,
    [stype+'PageSize']:     pagination.page_size,
    [stype+'TotalPages']:   pagination.total_pages,
    [stype+'TotalResults']: pagination.total_results,

    [stype+'ApiStatus']: apiStatus,

    [stype+'Loading']: false,
  });

  return nextState;
}

function _handleFetchRoutesError(type, state, payload) {
  const stype = _stateType(type);
  let nextState = Object.assign({}, state, {
    [stype+'Loading']: false,
    [stype+'Error']: payload.error
  });

  return nextState;
}

export default function reducer(state=initialState, action) {

  switch(action.type) {
    case LOCATION_CHANGE:
      return _handleLocationChange(state, action.payload);

    // Routes Received
    case FETCH_ROUTES_RECEIVED_REQUEST:
      return _handleFetchRoutesRequest(ROUTES_RECEIVED,
                                       state,
                                       action.payload);
    case FETCH_ROUTES_RECEIVED_SUCCESS:
      return _handleFetchRoutesSuccess(ROUTES_RECEIVED,
                                       state,
                                       action.payload);
    case FETCH_ROUTES_RECEIVED_ERROR:
      return _handleFetchRoutesError(ROUTES_RECEIVED,
                                     state,
                                     action.payload);

    // Routes Filtered
    case FETCH_ROUTES_FILTERED_REQUEST:
      return _handleFetchRoutesRequest(ROUTES_FILTERED,
                                       state,
                                       action.payload);
    case FETCH_ROUTES_FILTERED_SUCCESS:
      return _handleFetchRoutesSuccess(ROUTES_FILTERED,
                                       state,
                                       action.payload);
    case FETCH_ROUTES_FILTERED_ERROR:
      return _handleFetchRoutesError(ROUTES_FILTERED,
                                     state,
                                     action.payload);

    // Routes Not Exported
    case FETCH_ROUTES_NOT_EXPORTED_REQUEST:
      return _handleFetchRoutesRequest(ROUTES_NOT_EXPORTED,
                                       state,
                                       action.payload);
    case FETCH_ROUTES_NOT_EXPORTED_SUCCESS:
      return _handleFetchRoutesSuccess(ROUTES_NOT_EXPORTED,
                                       state,
                                       action.payload);
    case FETCH_ROUTES_NOT_EXPORTED_ERROR:
      return _handleFetchRoutesError(ROUTES_NOT_EXPORTED,
                                     state,
                                     action.payload);
  }

  return state;
}


