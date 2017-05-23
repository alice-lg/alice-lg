
// Routeserver Reducer

import {LOAD_ROUTESERVERS_REQUEST,
        LOAD_ROUTESERVERS_SUCCESS,
        LOAD_ROUTESERVER_STATUS_SUCCESS,

        LOAD_ROUTESERVER_PROTOCOL_REQUEST,
        LOAD_ROUTESERVER_PROTOCOL_SUCCESS,

        LOAD_ROUTESERVER_ROUTES_REQUEST,
        LOAD_ROUTESERVER_ROUTES_SUCCESS,

        LOAD_ROUTESERVER_ROUTES_FILTERED_REQUEST,
        LOAD_ROUTESERVER_ROUTES_FILTERED_SUCCESS,

        SET_PROTOCOLS_FILTER_VALUE,
        SET_ROUTES_FILTER_VALUE}
  from './actions'

import {LOAD_REJECT_REASONS_SUCCESS,
        LOAD_NOEXPORT_REASONS_SUCCESS}
  from './large-communities/actions'


const initialState = {
  all: [],
  filtered: {},
  details: {},
  protocols: {},
  routes: {},

  reject_reasons: {},
  reject_id: 0,
  reject_asn: 0,

  noexport_reasons: {},
  noexport_id: 0,
  noexport_asn: 0,

  protocolsFilterValue: "",
  routesFilterValue: "",

  isLoading: false,

  routesAreLoading: false,
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

    case LOAD_ROUTESERVER_ROUTES_REQUEST:
    case LOAD_ROUTESERVER_ROUTES_FILTERED_REQUEST:
      return Object.assign({}, state, {
        routesAreLoading: true
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

    case LOAD_ROUTESERVER_ROUTES_SUCCESS:
      var routes = Object.assign({}, state.routes, {
        [action.payload.protocolId]: action.payload.routes
      });
      return Object.assign({}, state, {
        routes: routes,
        routesAreLoading: false
      });

    case LOAD_ROUTESERVER_ROUTES_FILTERED_SUCCESS:
      var filtered = Object.assign({}, state.filtered, {
        [action.payload.protocolId]: action.payload.routes
      });
      return Object.assign({}, state, {
        filtered: filtered,
        routesAreLoading: false
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

    case SET_PROTOCOLS_FILTER_VALUE:
    case SET_ROUTES_FILTER_VALUE:
      return Object.assign({}, state, action.payload);

  }
  return state;
}



