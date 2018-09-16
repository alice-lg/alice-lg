
// Routeserver Reducer

import {LOAD_ROUTESERVERS_REQUEST,
        LOAD_ROUTESERVERS_SUCCESS,

        LOAD_ROUTESERVER_STATUS_SUCCESS,
        LOAD_ROUTESERVER_STATUS_ERROR,

        LOAD_ROUTESERVER_PROTOCOL_REQUEST,
        LOAD_ROUTESERVER_PROTOCOL_SUCCESS}
  from './actions'

import {LOAD_REJECT_REASONS_SUCCESS,
        LOAD_NOEXPORT_REASONS_SUCCESS}
  from './large-communities/actions'


const initialState = {

  all: [],

  details: {},
  protocols: {},
  statusErrors: {},

  reject_reasons: {},
  reject_id: 0,
  reject_asn: 0,

  noexport_reasons: {},
  noexport_id: 0,
  noexport_asn: 0,

  isLoading: false,

  protocolsAreLoading: false
};


export default function reducer(state = initialState, action) {
  switch(action.type) {
    case LOAD_ROUTESERVERS_REQUEST:
      return Object.assign({}, state, {
        isLoading: true
      });

    case LOAD_ROUTESERVERS_SUCCESS:
      return Object.assign({}, state, {
        all: action.payload.routeservers,
        isLoading: false
      });

    case LOAD_ROUTESERVER_PROTOCOL_REQUEST:
      return Object.assign({}, state, {
        protocolsAreLoading: true
      })

    case LOAD_ROUTESERVER_PROTOCOL_SUCCESS:
      var protocols = Object.assign({}, state.protocols, {
        [action.payload.routeserverId]: action.payload.protocol
      });
      return Object.assign({}, state, {
        protocols: protocols,
        protocolsAreLoading: false
      });

    case LOAD_NOEXPORT_REASONS_SUCCESS:
    case LOAD_REJECT_REASONS_SUCCESS:
      return Object.assign({}, state, action.payload);


    case LOAD_ROUTESERVER_STATUS_SUCCESS:
      var details = Object.assign({}, state.details, {
        [action.payload.routeserverId]: action.payload.status
      });
      var errors = Object.assign({}, state.statusErrors, {
        [action.payload.routeserverId]: null,
      });
      return Object.assign({}, state, {
        details: details,
        statusErrors: errors 
      });

    case LOAD_ROUTESERVER_STATUS_ERROR:
      var info = {
        code: 42,
        tag: "UNKNOWN_ERROR",
        message: "Unknown error"
      };

      if (action.payload.error &&
          action.payload.error.response && 
          action.payload.error.response.data &&
          action.payload.error.response.data.code) {
            info = action.payload.error.response.data;
      }
      
      var errors = Object.assign({}, state.statusErrors, {
        [action.payload.routeserverId]: info
      });
      return Object.assign({}, state, {
        statusErrors: errors 
      });
      return state;
  }
  return state;
}


