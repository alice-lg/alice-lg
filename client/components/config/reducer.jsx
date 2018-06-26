import {LOAD_CONFIG_SUCCESS} from './actions'

const initialState = {
  routes_columns: {
    Gateway: "gateway",
    Interface: "interface",
    Metric: "metric",
  },
  prefix_lookup_enabled: false,
  content: {}
};


export default function reducer(state = initialState, action) {
  switch(action.type) {
    case LOAD_CONFIG_SUCCESS:
      return {
        routes_columns: action.payload.routes_columns,
        prefix_lookup_enabled: action.payload.prefix_lookup_enabled
       };
  }
  return state;
}



