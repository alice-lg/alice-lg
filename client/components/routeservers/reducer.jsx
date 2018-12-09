
// Routeserver Reducer

import {LOAD_ROUTESERVERS_REQUEST,
        LOAD_ROUTESERVERS_SUCCESS,

        LOAD_ROUTESERVER_STATUS_SUCCESS,
        LOAD_ROUTESERVER_STATUS_ERROR,

        LOAD_ROUTESERVER_PROTOCOL_REQUEST,
        LOAD_ROUTESERVER_PROTOCOL_SUCCESS,

        SELECT_GROUP}
  from './actions'

import {LOAD_CONFIG_SUCCESS} from 'components/config/actions'

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE'

const initialState = {

  all: [],
  byId: {},

  selectedRsId: "",

  groups: [],
  isGrouped: false,
  selectedGroup: "",

  details: {},
  protocols: {},
  statusErrors: {},

  rejectReasons: {},
  noexportReasons: {},

  rejectCandidates: {
    communities: {}
  },

  isLoading: false,

  protocolsAreLoading: false
};

// == Helpers ==
const _groupForRsId = function(routeservers, rsId) {
  const rs = routeservers[rsId]||{group:""};
  return rs.group;
}


// == Handlers ==
const _importConfig = function(state, payload) {
  // Get reject and filter reasons from config
  const rejectReasons   = payload.reject_reasons;
  const noexportReasons = payload.noexport_reasons;

  // Get reject candidates from config
  const rejectCandidates = payload.reject_candidates;

  return Object.assign({}, state, {
    rejectReasons:    rejectReasons,
    rejectCandidates: rejectCandidates,

    noexportReasons: noexportReasons
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

const _loadRouteservers = function(state, routeservers) {
  // Caclulate grouping
  const byId = {};
  let groups = [];
  for (const rs of routeservers) {
    byId[rs.id] = rs;
    if (groups.indexOf(rs.group) == -1) {
      groups.push(rs.group);
    }
  }

  const selectedGroup = _groupForRsId(byId, state.selectedRsId);

  return Object.assign({}, state, {
    all:  routeservers,
    byId: byId,

    groups: groups,
    isGrouped: groups.length > 1,
    selectedGroup: selectedGroup,
    isLoading: false
  });
}


const _restoreStatefromLocation = function(state, location) {
  const path = location.pathname.split("/");
  if(path.length < 2) {
    return state; // nothing to do here
  }
  const [resource, routeserverId, ...rest] = path;

  if (resource != "routeservers") {
    return state;
  }

  let selectedGroup = state.selectedGroup;
  if (state.all.length > 0) {
    selectedGroup = _groupForRsId(state.byId, routeserverId);
  }

  return Object.assign({}, state, {
    selectedRsId: routeserverId,
    selectedGroup: selectedGroup,
  });
}


export default function reducer(state = initialState, action) {
  switch(action.type) {
    case LOAD_ROUTESERVERS_REQUEST:
      return Object.assign({}, state, {
        isLoading: true
      });

    case LOAD_ROUTESERVERS_SUCCESS:
      return _loadRouteservers(state, action.payload.routeservers);

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

    case SELECT_GROUP:
      return Object.assign({}, state, {
        selectedGroup: action.payload,
      });

    case LOCATION_CHANGE:
      return _restoreStatefromLocation(state, action.payload);
  }
  return state;
}


