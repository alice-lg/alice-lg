import {LOAD_CONFIG_SUCCESS} from './actions'

const initialState = {
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
};


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

        noexport_load_on_demand: action.payload.noexport.load_on_demand
       });
  }
  return state;
}



