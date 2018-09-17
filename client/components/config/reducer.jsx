import {LOAD_CONFIG_SUCCESS} from './actions'
import {LOAD_ROUTESERVERS_SUCCESS} from 'components/routeservers/actions'

const initialState = {
  asn: 0, // Our own ASN (might be abstracted in the future)

  routes_columns: {},
  routes_columns_order: [],
  neighbours_columns: {},
  neighbours_columns_order: [],
  lookup_columns: {},
  lookup_columns_order: [],
  prefix_lookup_enabled: false,
  content: {},
  noexport_load_on_demand: true, // we have to assume this
                                 // otherwise fetch will start.
  bgp_communities: {},

  blackholes: {}, // Map blackholes to routeservers
  asns: {}, // Map ASNs to routeservers (for future use)
};

const _handleRouteserversConfig = function(state, payload) {
  let blackholes = {};
  let asns = {};
  let asn = 0;
  for (const rs of payload.routeservers) {
    blackholes[rs.id] = rs.blackholes; 
    asns[rs.is] = rs.asn;
    if (!asn) {
      asn = rs.asn; // Just go with the first asn as our own
    }
  }

  return Object.assign({}, state, {
    asn: asn,
    blackholes: blackholes,
    asns: asns,
  });
}

export default function reducer(state = initialState, action) {
  switch(action.type) {
    case LOAD_CONFIG_SUCCESS:
      return Object.assign({}, state, {
        routes_columns:       action.payload.routes_columns,
        routes_columns_order: action.payload.routes_columns_order,

        neighbours_columns:       action.payload.neighbours_columns,
        neighbours_columns_order: action.payload.neighbours_columns_order,

        lookup_columns: action.payload.lookup_columns,
        lookup_columns_order: action.payload.lookup_columns_order,

        prefix_lookup_enabled: action.payload.prefix_lookup_enabled,

        bgp_communities: action.payload.bgp_communities,
        noexport_load_on_demand: action.payload.noexport.load_on_demand
       });

    case LOAD_ROUTESERVERS_SUCCESS:
      return _handleRouteserversConfig(state, action.payload);
  }
  return state;
}

