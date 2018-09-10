
/**
 * Alice (formerly known as Birdseye) v.2.0.0
 * ------------------------------------------
 *
 * @author Matthias Hannig <mha@ecix.net>
 */

import axios     from 'axios'

import React     from 'react'
import ReactDOM  from 'react-dom'

import { Component } from 'react'

// Config
import { configureAxios } from './config'

// Content
import { contentUpdate } from './components/content/actions'

// Redux
import { createStore, applyMiddleware } from 'redux'
import { Provider } from 'react-redux'

// Router
import { createHistory } from 'history'
import { Router,
         Route,
         IndexRoute,
         IndexRedirect,
         useRouterHistory } from 'react-router'

import { syncHistoryWithStore } from 'react-router-redux'


// Components
import LayoutMain from 'layouts/main'


import WelcomePage
  from 'components/welcome'
import RouteserverPage
  from 'components/routeservers/page'
import RoutesPage
  from 'components/routeservers/routes/page'

// Middlewares
import thunkMiddleware from 'redux-thunk'
import createLogger from 'redux-logger'
import { routerMiddleware as createRouterMiddleware }
  from 'react-router-redux'

// Reducer
import combinedReducer from './reducer/app-reducer'

// Setup routing
const browserHistory = useRouterHistory(createHistory)({
  basename: '/alice'
});


// Setup application
let store;
const routerMiddleware = createRouterMiddleware(browserHistory);
if (window.NO_LOG) {
  store = createStore(combinedReducer, applyMiddleware(
    routerMiddleware,
    thunkMiddleware
  ));
} else {
  const loggerMiddleware = createLogger();
  store = createStore(combinedReducer, applyMiddleware(
    routerMiddleware,
    thunkMiddleware,
    loggerMiddleware
  ));
}


// Create extension endpoint:
window.Alice = {
  updateContent: (content) => {
    store.dispatch(contentUpdate(content));    
  }
};

const history = syncHistoryWithStore(browserHistory, store);

// Setup axios
configureAxios(axios);

// Create App
class Birdseye extends Component {
  render() {
    return (
      <Provider store={store}>
        <Router history={history}>
          <Route path="/" component={LayoutMain}>
            <IndexRoute component={WelcomePage}/>
            <Route path="/routeservers">
              <Route path=":routeserverId" component={RouteserverPage} />
              <Route path=":routeserverId/protocols/:protocolId/routes" component={RoutesPage} />
            </Route>
          </Route>
        </Router>
      </Provider>
    );
  }
}

var mount = document.getElementById('app');
ReactDOM.render(<Birdseye />, mount);

