
/**
 * Routeservers Actions
 */

import axios from 'axios'

import {apiError} from 'components/errors/actions'

export const LOAD_ROUTESERVERS_REQUEST = '@routeservers/LOAD_ROUTESERVERS_REQUEST';
export const LOAD_ROUTESERVERS_SUCCESS = '@routeservers/LOAD_ROUTESERVERS_SUCCESS';
export const LOAD_ROUTESERVERS_ERROR   = '@routeservers/LOAD_ROUTESERVERS_ERROR';

export const LOAD_ROUTESERVER_STATUS_REQUEST = '@routeservers/LOAD_ROUTESERVER_STATUS_REQUEST';
export const LOAD_ROUTESERVER_STATUS_SUCCESS = '@routeservers/LOAD_ROUTESERVER_STATUS_SUCCESS';
export const LOAD_ROUTESERVER_STATUS_ERROR   = '@routeservers/LOAD_ROUTESERVER_STATUS_ERROR';

export const LOAD_ROUTESERVER_PROTOCOL_REQUEST = '@routeservers/LOAD_ROUTESERVER_PROTOCOL_REQUEST';
export const LOAD_ROUTESERVER_PROTOCOL_SUCCESS = '@routeservers/LOAD_ROUTESERVER_PROTOCOL_SUCCESS';
export const LOAD_ROUTESERVER_PROTOCOL_ERROR   = '@routeservers/LOAD_ROUTESERVER_PROTOCOL_ERROR';

export const SELECT_GROUP = "@routeservers/SELECT_GROUP";


// Action Creators
export function loadRouteserversRequest() {
  return {
    type: LOAD_ROUTESERVERS_REQUEST
  }
}

export function loadRouteserversSuccess(routeservers) {
  return {
    type: LOAD_ROUTESERVERS_SUCCESS,
    payload: {
      routeservers: routeservers
    }
  }
}

export function loadRouteserversError(error) {
  return {
    type: LOAD_ROUTESERVERS_ERROR,
    payload: {
      error: error
    }
  }
}

export function loadRouteservers() {
  return (dispatch) => {
    dispatch(loadRouteserversRequest())

    axios.get('/api/v1/routeservers')
      .then(
        ({data}) => {
          dispatch(loadRouteserversSuccess(data["routeservers"]));
        },
        (error) => {
          dispatch(apiError(error));
          dispatch(loadRouteserversError(error.data));
        });
  }
}



export function loadRouteserverStatusRequest(routeserverId) {
  return {
    type: LOAD_ROUTESERVER_STATUS_REQUEST,
    payload: {
      routeserverId: routeserverId
    }
  }
}

export function loadRouteserverStatusSuccess(routeserverId, status) {
  return {
    type: LOAD_ROUTESERVER_STATUS_SUCCESS,
    payload: {
      status: status,
      routeserverId: routeserverId
    }
  }
}

export function loadRouteserverStatusError(routeserverId, error) {
  return {
    type: LOAD_ROUTESERVER_STATUS_ERROR,
    payload: {
      error: error,
      routeserverId: routeserverId
    }
  }
}

export function loadRouteserverStatus(routeserverId) {
  return (dispatch) => {
    dispatch(loadRouteserverStatusRequest(routeserverId));
    axios.get(`/api/v1/routeservers/${routeserverId}/status`)
      .then(
        ({data}) => {
          dispatch(loadRouteserverStatusSuccess(routeserverId, data.status));
        },
        (error) => {
          dispatch(apiError(error));
          dispatch(loadRouteserverStatusError(routeserverId, error));
        });
  }
}


export function loadRouteserverProtocolRequest(routeserverId) {
  return {
    type: LOAD_ROUTESERVER_PROTOCOL_REQUEST,
    payload: {
      routeserverId: routeserverId,
    }
  }
}

export function loadRouteserverProtocolSuccess(routeserverId, protocol, api) {
  return {
    type: LOAD_ROUTESERVER_PROTOCOL_SUCCESS,
    payload: {
      routeserverId: routeserverId,
      protocol: protocol,
      api: api
    }
  }
}

export function loadRouteserverProtocol(routeserverId) {
  return (dispatch) => {
    dispatch(loadRouteserverProtocolRequest(routeserverId));
    axios.get(`/api/v1/routeservers/${routeserverId}/neighbors`)
      .then(
        ({data}) => {
          dispatch(loadRouteserverProtocolSuccess(
            routeserverId,
            data.neighbors,
            data.api,
          ));
        },
        (error) => dispatch(apiError(error)));
  }
}


export function selectGroup(group) {
  return {
    type: SELECT_GROUP,
    payload: group,
  }
}

