
/*
 * This will migrate to become the neighbors
 * reducer. Currently neihgbors are stored in
 * the routeserver reducer.
 */

import {SET_FILTER_VALUE} from './actions'

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE';

const DEFAULT_SORT_COLUMN = "asn";
const DEFAULT_SORT_ORDER = "asc";

const initialState = {
  sortColumn: DEFAULT_SORT_COLUMN,
  sortOrder: DEFAULT_SORT_ORDER,

  filterQuery: "",
  filterValue: ""
};


// Reducer functions

function _handleLocationChange(state, payload) {
  const query = payload.query;
  const sort = query["s"] || DEFAULT_SORT_COLUMN;
  const order = query["o"]  || DEFAULT_SORT_ORDER; 
  const filterQuery = query["q"] || "";

  return Object.assign({}, state, {
    sortColumn: sort,
    sortOrder: order,
    
    filterQuery: filterQuery,
    filterValue: filterQuery
  });
}


export default function(state=initialState, action) {
  switch (action.type) {
    case LOCATION_CHANGE:
      return _handleLocationChange(state, action.payload);

    case SET_FILTER_VALUE:
      return Object.assign({}, state, {
        filterValue: action.payload.value
      });

    default:
  }

  return state;
}

