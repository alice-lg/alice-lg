
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
        ROUTES_NOT_EXPORTED} from './actions'

import {SET_FILTER_QUERY_VALUE} from './actions'

import {cloneFilters, decodeFiltersApplied, initializeFilterState}
  from 'components/filters/state'

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE'

const initialState = {

  received: [],
  receivedLoading: false,
  receivedRequested: false,
  receivedError: null,
  receivedPage: 0,
  receivedPageSize: 0,
  receivedTotalPages: 0,
  receivedTotalResults: 0,
  receivedApiStatus: {},
  receivedFiltersApplied: initializeFilterState(),
  receivedFiltersAvailable: initializeFilterState(),

  filtered: [],
  filteredLoading: false,
  filteredRequested: false,
  filteredError: null,
  filteredPage: 0,
  filteredPageSize: 0,
  filteredTotalPages: 0,
  filteredTotalResults: 0,
  filteredApiStatus: {},
  filteredFiltersApplied: initializeFilterState(),
  filteredFiltersAvailable: initializeFilterState(),

  notExported: [],
  notExportedLoading: false,
  notExportedRequested: false,
  notExportedError: null,
  notExportedPage: 0,
  notExportedPageSize: 0,
  notExportedTotalPages: 0,
  notExportedTotalResults: 0,
  notExportedApiStatus: {},
  notExportedFiltersApplied: initializeFilterState(),
  notExportedFiltersAvailable: initializeFilterState(),

  // Derived state from location
  loadNotExported: false,

  filterValue: "",
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
  const query = payload.query;

  const filterQuery = query["q"] || "";

  const receivedPage    = parseInt(query["pr"] || 0, 10);
  const filteredPage    = parseInt(query["pf"] || 0, 10);
  const notExportedPage = parseInt(query["pn"] || 0, 10);

  // Determine on demand loading state
  const loadNotExported = parseInt(query["ne"] || 0, 10) === 1 ? true : false;

  // Restore filters applied from location
  const filtersApplied = decodeFiltersApplied(query);

  const nextState = Object.assign({}, state, {
    filterQuery: filterQuery,
    filterValue: filterQuery, // location overrides form

    receivedPage:    receivedPage,
    filteredPage:    filteredPage,
    notExportedPage: notExportedPage,

    loadNotExported: loadNotExported,

    receivedFiltersApplied:    filtersApplied,
    filteredFiltersApplied:    filtersApplied,
    notExportedFiltersApplied: filtersApplied,
  });

  return nextState;
}


function _handleFetchRoutesRequest(type, state, payload) {
  const stype = _stateType(type);
  const nextState = Object.assign({}, state, {
    [stype+'Loading']: true,
    [stype+'Requested']: true,
    [stype+'FiltersAvailable']: initializeFilterState(),
  });

  return nextState;
}

function _handleFetchRoutesSuccess(type, state, payload) {
  const stype = _stateType(type);
  const pagination = payload.pagination;
  const apiStatus = payload.apiStatus;

  let nextState = Object.assign({}, state, {
    [stype]: payload.routes,

    [stype+'Page']:         pagination.page,
    [stype+'PageSize']:     pagination.page_size,
    [stype+'TotalPages']:   pagination.total_pages,
    [stype+'TotalResults']: pagination.total_results,

    [stype+'ApiStatus']: apiStatus,

    [stype+'Loading']: false,

    [stype+'FiltersAvailable']: cloneFilters(payload.filtersAvailable),
    [stype+'FiltersApplied']:   cloneFilters(payload.filtersApplied),
  });

  return nextState;
}

function _handleFetchRoutesError(type, state, payload) {
  const stype = _stateType(type);
  let nextState = Object.assign({}, state, {
    [stype+'Loading']: false,
    [stype+'Requested']: false,
    [stype+'Error']: payload.error
  });

  return nextState;
}

function _handleFilterQueryValueChange(state, payload) {
  return Object.assign({}, state, {
    filterValue: payload.value
  });
}


export default function reducer(state=initialState, action) {

  switch(action.type) {
    case LOCATION_CHANGE:
      return _handleLocationChange(state, action.payload);

    case SET_FILTER_QUERY_VALUE:
      return _handleFilterQueryValueChange(state, action.payload);

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


