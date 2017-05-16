

import { combineReducers } from 'redux'

import bgpAttributesModalReducer
  from 'components/routeservers/routes/bgp-attributes-modal-reducer'

export default combineReducers({
	bgpAttributes: bgpAttributesModalReducer
});


