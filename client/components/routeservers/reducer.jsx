
// Routeserver Reducer

import {LOAD_ROUTESERVERS_REQUEST,
        LOAD_ROUTESERVERS_SUCCESS,

        LOAD_ROUTESERVER_STATUS_SUCCESS,
        LOAD_ROUTESERVER_STATUS_ERROR,

        LOAD_ROUTESERVER_PROTOCOL_REQUEST,
        LOAD_ROUTESERVER_PROTOCOL_SUCCESS}
  from './actions'

import {LOAD_CONFIG_SUCCESS} from 'components/config/actions'

const initialState = {

  all: [],

  details: {},
  protocols: {},
  statusErrors: {},

  rejectReasons: {},
  rejectId: 0,
  rejectAsn: 0,

  noexportReasons: {},
  noexportId: 0,
  noexportAsn: 0,

  rejectCandidates: {
    communities: {}
  },

  isLoading: false,

  protocolsAreLoading: false
};


// == Handlers ==
const _importConfig = function(state, payload) {
  // Get reject and filter reasons from config
  const rejectReasons = payload.reject_reasons;
  const rejectId      = payload.rejection.reject_id;
  const rejectAsn     = payload.rejection.asn;

  const noexportReasons = payload.noexport_reasons;
  const noexportId      = payload.noexport.noexport_id;
  const noexportAsn     = payload.noexport.asn;

  // Get reject candidates from config
  const rejectCandidates = payload.reject_candidates;

  return Object.assign({}, state, {
    rejectReasons: rejectReasons,
    rejectAsn:     rejectAsn,
    rejectId:      rejectId,

    rejectCandidates: rejectCandidates, 

    noexportReasons: noexportReasons,
    noexportAsn:     noexportAsn,
    noexportId:      noexportId,
  });
};


const _updateStatus = function(state, payload) {
  const details = Object.assign({}, state.details, {
    [payload.routeserverId]: payload.status
  });
  const errors = Object.assign({}, state.statusErrors, {
    [payload.routeserverId]: null,
  });

  return Object.assign({}, state, {
    details: details,
    statusErrors: errors 
  });
}


const _updateStatusError = function(state, payload) {
  var info = {
    code: 42,
    tag: "UNKNOWN_ERROR",
    message: "Unknown error"
  };

  if (payload.error &&
      payload.error.response && 
      payload.error.response.data &&
      payload.error.response.data.code) {
        info = payload.error.response.data;
  }
  
  var errors = Object.assign({}, state.statusErrors, {
    [payload.routeserverId]: info
  });

  return Object.assign({}, state, {
    statusErrors: errors 
  });
}


const _updateProtocol = function(state, payload) {
  var protocols = Object.assign({}, state.protocols, {
    [payload.routeserverId]: payload.protocol
  });

  return Object.assign({}, state, {
    protocols: protocols,
    protocolsAreLoading: false
  });
}


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
      return _updateProtocol(state, action.payload);

    case LOAD_CONFIG_SUCCESS:
      return _importConfig(state, action.payload);

    case LOAD_ROUTESERVER_STATUS_SUCCESS:
      return _updateStatus(state, action.payload);

    case LOAD_ROUTESERVER_STATUS_ERROR:
      return _updateStatusError(state, action.payload);
  }
  return state;
}


