
/**
 * Routeservers Actions
 */

import axios from 'axios'

import {apiError} from 'components/errors/actions'

export const LOAD_ROUTESERVERS_REQUEST = '@birdseye/LOAD_ROUTESERVERS_REQUEST';
export const LOAD_ROUTESERVERS_SUCCESS = '@birdseye/LOAD_ROUTESERVERS_SUCCESS';
export const LOAD_ROUTESERVERS_ERROR   = '@birdseye/LOAD_ROUTESERVERS_ERROR';

export const LOAD_ROUTESERVER_STATUS_REQUEST = '@birdseye/LOAD_ROUTESERVER_STATUS_REQUEST';
export const LOAD_ROUTESERVER_STATUS_SUCCESS = '@birdseye/LOAD_ROUTESERVER_STATUS_SUCCESS';
export const LOAD_ROUTESERVER_STATUS_ERROR   = '@birdseye/LOAD_ROUTESERVER_STATUS_ERROR';

export const LOAD_ROUTESERVER_PROTOCOL_REQUEST = '@birdseye/LOAD_ROUTESERVER_PROTOCOL_REQUEST';
export const LOAD_ROUTESERVER_PROTOCOL_SUCCESS = '@birdseye/LOAD_ROUTESERVER_PROTOCOL_SUCCESS';
export const LOAD_ROUTESERVER_PROTOCOL_ERROR   = '@birdseye/LOAD_ROUTESERVER_PROTOCOL_ERROR';


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

    axios.get('/api/routeservers')
      .then(({data}) => {
        dispatch(loadRouteserversSuccess(data["routeservers"]));
      })
      .catch((error) => {
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
    axios.get(`/api/routeservers/${routeserverId}/status`)
      .then(({data}) => {
        dispatch(loadRouteserverStatusSuccess(routeserverId, data.status));
      })
      .catch((error) => {
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
    axios.get(`/api/routeservers/${routeserverId}/neighbours`)
      .then(({data}) => {
        console.log("LRS:", data);
        dispatch(loadRouteserverProtocolSuccess(
          routeserverId,
          data.neighbours,
          data.api,
        ));
      })
      .catch((error) => dispatch(apiError(error)));
  }
}

