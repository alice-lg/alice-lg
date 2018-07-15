

const LOCATION_CHANGE = '@@router/LOCATION_CHANGE'

const initialState = {
  received: [],
  filtered: [],
  notExported: [],

  receivedPage: 0,
  receivedTotalPages: 0,
  receivedTotalResults: 0,

  filteredPage: 0,
  filteredTotalPages: 0,
  filteredTotalResults: 0,

  notExportedPage: 0,
  notExportedTotalPages: 0,
  notExportedTotalResults: 0,

  receivedLoading: false,
  filteredLoading: false,
  notExportedLoading: false,

  filterQuery: "",
}


// Handlers:
function _handleLocationChange(state, payload) {
  // Check query payload
  let query = payload.query;

  let filterQuery = query["q"] || "";

  let receivedPage    = query["pr"] || 0;
  let filteredPage    = query["pf"] || 0;
  let notExportedPage = query["pn"] || 0;

  // Assert numeric
  receivedPage    = parseInt(receivedPage);
  filteredPage    = parseInt(filteredPage);
  notExportedPage = parseInt(notExportedPage);

  let nextState = Object.assign({}, state, {
    filterQuery: filterQuery,

    receivedPage:    receivedPage,
    filteredPage:    filteredPage,
    notExportedPage: notExportedPage,
  });

  return nextState;
}



export default function reducer(state=initialState, action) {

  switch(action.type) {
    case LOCATION_CHANGE:
      return _handleLocationChange(state, action.payload);
  }

  return state;
}


