
import { combineReducers } from 'redux'


// Library Reducers
import { routerReducer } from 'react-router-redux'

// Application Reducers
import routeserversReducer
  from 'components/routeservers/reducer'

import modalsReducer
  from 'components/modals/reducer'

import errorsReducer
  from 'components/errors/reducer'

import configReducer
  from 'components/config/reducer'

import contentReducer
  from 'components/content/reducer'

import lookupReducer
	from 'components/lookup/reducer'

import routesReducer
  from 'components/routeservers/routes/reducer'

import neighborsReducer
  from 'components/routeservers/protocols/reducer'

export default combineReducers({
  routeservers:  routeserversReducer,
  neighbors:     neighborsReducer,
  routes:        routesReducer,
  modals:        modalsReducer,
  routing:       routerReducer,
	lookup:				 lookupReducer,
  errors:        errorsReducer,
  config:        configReducer,
  content:       contentReducer,
});

