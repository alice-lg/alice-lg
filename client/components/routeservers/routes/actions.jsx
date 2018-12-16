
import axios from 'axios'

import {apiError} from 'components/errors/actions'
import {filtersUrlEncode} from 'components/filters/encoding'

export const ROUTES_RECEIVED = "received";
export const ROUTES_FILTERED = "filtered";
export const ROUTES_NOT_EXPORTED = "notExported";

export const FETCH_ROUTES_RECEIVED_REQUEST = "@routes/FETCH_ROUTES_RECEIVED_REQUEST";
export const FETCH_ROUTES_RECEIVED_SUCCESS = "@routes/FETCH_ROUTES_RECEIVED_SUCCESS";
export const FETCH_ROUTES_RECEIVED_ERROR   = "@routes/FETCH_ROUTES_RECEIVED_ERROR";

export const FETCH_ROUTES_FILTERED_REQUEST = "@routes/FETCH_ROUTES_FILTERED_REQUEST";
export const FETCH_ROUTES_FILTERED_SUCCESS = "@routes/FETCH_ROUTES_FILTERED_SUCCESS";
export const FETCH_ROUTES_FILTERED_ERROR   = "@routes/FETCH_ROUTES_FILTERED_ERROR";

export const FETCH_ROUTES_NOT_EXPORTED_REQUEST = "@routes/FETCH_ROUTES_NOT_EXPORTED_REQUEST";
export const FETCH_ROUTES_NOT_EXPORTED_SUCCESS = "@routes/FETCH_ROUTES_NOT_EXPORTED_SUCCESS";
export const FETCH_ROUTES_NOT_EXPORTED_ERROR   = "@routes/FETCH_ROUTES_NOT_EXPORTED_ERROR";

export const SET_FILTER_QUERY_VALUE = "@routes/SET_FILTER_QUERY_VALUE";

// Url helper
function routesUrl(type, rsId, pId, page, query, filters) {
    let rtype = type;
    if (type == ROUTES_NOT_EXPORTED) {
      rtype = "not-exported"; // This is a bit ugly
    }

    const filtersEncoded = filtersUrlEncode(filters);
    const base = `/api/v1/routeservers/${rsId}/neighbors/${pId}/routes/${rtype}`
    const params = `?page=${page}&q=${query}${filtersEncoded}`

    return base + params;
};


// Meta Creators
function routesRequest(type) {
  return () => ({
    type: type,
    payload: {},
  });
}

function routesSuccess(type) {
  const rtype = {
    [FETCH_ROUTES_RECEIVED_SUCCESS]: 'imported',
    [FETCH_ROUTES_FILTERED_SUCCESS]: 'filtered',
    [FETCH_ROUTES_NOT_EXPORTED_SUCCESS]: 'not_exported',
  }[type];

  return (data) => ({
    type: type,
    payload: {
      routes: data[rtype],
      pagination: data.pagination,
      filtersApplied: data.filters_applied,
      filtersAvailable: data.filters_available,
      apiStatus: data.api,
    }
  });
}

function routesError(type) {
  return (error) => ({
    type: type,
    payload: {
      error: error
    }
  });
};

// Action Creators: Routes Received
export const fetchRoutesReceivedRequest = routesRequest(FETCH_ROUTES_RECEIVED_REQUEST);
export const fetchRoutesReceivedSuccess = routesSuccess(FETCH_ROUTES_RECEIVED_SUCCESS);
export const fetchRoutesReceivedError   = routesError(FETCH_ROUTES_RECEIVED_ERROR);
export const fetchRoutesReceived        = fetchRoutes(ROUTES_RECEIVED);

// Action Creators: Routes Filtered
export const fetchRoutesFilteredRequest = routesRequest(FETCH_ROUTES_FILTERED_REQUEST);
export const fetchRoutesFilteredSuccess = routesSuccess(FETCH_ROUTES_FILTERED_SUCCESS);
export const fetchRoutesFilteredError   = routesError(FETCH_ROUTES_FILTERED_ERROR);
export const fetchRoutesFiltered        = fetchRoutes(ROUTES_FILTERED);

// Action Creators: Routes Filtered
export const fetchRoutesNotExportedRequest = routesRequest(FETCH_ROUTES_NOT_EXPORTED_REQUEST);
export const fetchRoutesNotExportedSuccess = routesSuccess(FETCH_ROUTES_NOT_EXPORTED_SUCCESS);
export const fetchRoutesNotExportedError   = routesError(FETCH_ROUTES_NOT_EXPORTED_ERROR);
export const fetchRoutesNotExported        = fetchRoutes(ROUTES_NOT_EXPORTED);

function fetchRoutes(type) {
  const requestAction = {
    [ROUTES_RECEIVED]:     fetchRoutesReceivedRequest,
    [ROUTES_FILTERED]:     fetchRoutesFilteredRequest,
    [ROUTES_NOT_EXPORTED]: fetchRoutesNotExportedRequest
  }[type];

  const successAction = {
    [ROUTES_RECEIVED]:     fetchRoutesReceivedSuccess,
    [ROUTES_FILTERED]:     fetchRoutesFilteredSuccess,
    [ROUTES_NOT_EXPORTED]: fetchRoutesNotExportedSuccess
  }[type];

  const errorAction = {
    [ROUTES_RECEIVED]:     fetchRoutesReceivedError,
    [ROUTES_FILTERED]:     fetchRoutesFilteredError,
    [ROUTES_NOT_EXPORTED]: fetchRoutesNotExportedError
  }[type];

  return (rsId, pId, page, query, filters) => {
    return (dispatch) => {
      dispatch(requestAction());
      axios.get(routesUrl(type, rsId, pId, page, query, filters))
        .then(
          ({data}) => {
            dispatch(successAction(data));
          },
          (error) => {
            dispatch(errorAction(error));
            dispatch(apiError(error));
          });
    }
  }
};

// Action Creators: Set Filter Query
export function setFilterQueryValue(value) {
  return {
    type: SET_FILTER_QUERY_VALUE,
    payload: {
      value: value
    }
  }
}

