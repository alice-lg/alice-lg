
import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'

import {fetchRoutesReceived,
        fetchRoutesFiltered,
        fetchRoutesNotExported} from './actions'

// Constants
import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';


/*
 * Render a RoutesView:
 * The routes view is a composit of:
 *  - A header
 *  - The Routes Table
 *  - A Paginator
 */

class RoutesView extends React.Component {

  componentDidMount() {
    const type = this.props.type;

    // Depending on the component's configuration, dispatch
    // routes fetching
    const fetchRoutes = {
      [ROUTES_RECEIVED]:     fetchRoutesReceived,
      [ROUTES_FILTERED]:     fetchRoutesFiltered,
      [ROUTES_NOT_EXPORTED]: fetchRoutesNotExported,
    }[type];

    // Gather required params
    const params = this.props.routes[type];
    const rsId = this.props.routeserverId;
    const pId = this.props.protocolId;
    const query = this.props.filterQuery;

    // Make request
    this.props.dispatch(fetchRoutes(rsId, pId, params.page, query));
  }

  render() {

    return (
      <div className="routes-view">
        [HEADER]<br />

        [TABLE]<br />

        [Paginator]
      </div>
    );
  }

}




export default connect(
  (state) => {
    let received = {
      routes:       state.routes.received,
      loading:      state.routes.receivedLoading,
      page:         state.routes.receivedPage,
      totalPages:   state.routes.receivedTotalPages,
      totalResults: state.routes.receivedTotalResults,
    };
    let filtered = {
      routes:       state.routes.filtered,
      loading:      state.routes.filteredLoading,
      page:         state.routes.filteredPage,
      totalPages:   state.routes.filteredTotalPages,
      totalResults: state.routes.filteredTotalResults,
    };
    let notExported = {
      routes:       state.routes.notExported,
      loading:      state.routes.notExportedLoading,
      page:         state.routes.notExportedPage,
      totalPages:   state.routes.notExportedTotalPages,
      totalResults: state.routes.notExportedTotalResults,
    };
    return({
      filterQuery: state.routes.filterQuery,
      routes: {
          [ROUTES_RECEIVED]:     received,
          [ROUTES_FILTERED]:     filtered,
          [ROUTES_NOT_EXPORTED]: notExported
      },
    });
  }
)(RoutesView);

