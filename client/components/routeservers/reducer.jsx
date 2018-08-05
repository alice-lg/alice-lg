
// Routeserver Reducer

import {LOAD_ROUTESERVERS_REQUEST,
        LOAD_ROUTESERVERS_SUCCESS,

        LOAD_ROUTESERVER_STATUS_SUCCESS,
        LOAD_ROUTESERVER_STATUS_ERROR,

        LOAD_ROUTESERVER_PROTOCOL_REQUEST,
        LOAD_ROUTESERVER_PROTOCOL_SUCCESS,

        SET_PROTOCOLS_FILTER_VALUE,
        SET_PROTOCOLS_FILTER,

        SET_ROUTES_FILTER_VALUE}
  from './actions'

import {LOAD_REJECT_REASONS_SUCCESS,
        LOAD_NOEXPORT_REASONS_SUCCESS}
  from './large-communities/actions'


const initialState = {

  all: [],

  errors: {},
  details: {},
  protocols: {},

  reject_reasons: {},
  reject_id: 0,
  reject_asn: 0,

  noexport_reasons: {},
  noexport_id: 0,
  noexport_asn: 0,

  protocolsFilterValue: "",
  protocolsFilter: "",

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
      return Object.assign({}, state, {
        details: details
      });

    case LOAD_ROUTESERVER_STATUS_ERROR:
      console.log("ROUTESERVER STATUS ERROR:", action);
      
      return state;

    case SET_PROTOCOLS_FILTER_VALUE:
      return Object.assign({}, state, {
        protocolsFilterValue: action.payload.value
      });

    case SET_PROTOCOLS_FILTER:
      return Object.assign({}, state, {
        protocolsFilter: action.payload.value
      });

    case SET_ROUTES_FILTER_VALUE:
      return Object.assign({}, state, action.payload);

  }
  return state;
}


