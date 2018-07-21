
import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'

import RoutesTable from './table'
import {RoutesPaginator,
        RoutesPaginationInfo} from './pagination'

import {fetchRoutesReceived,
        fetchRoutesFiltered,
        fetchRoutesNotExported} from './actions'

// Constants
import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';


const RoutesHeader = (props) => {
  const type = props.type;
  const color = {
    [ROUTES_RECEIVED]: "green",
    [ROUTES_FILTERED]: "orange",
    [ROUTES_NOT_EXPORTED]: "red"
  }[type];
  const rtype = {
    [ROUTES_RECEIVED]: "accepted",
    [ROUTES_FILTERED]: "filtered",
    [ROUTES_NOT_EXPORTED]: "not exported"
  }[type];
  return (<p style={{"color": color, "textTransform": "uppercase"}}>
            Routes {rtype}
          </p>);
};

/*
 * Render a RoutesView:
 * The routes view is a composit of:
 *  - A header
 *  - The Routes Table
 *  - A Paginator
 */

class RoutesView extends React.Component {

  dispatchFetchRoutes() {
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

  componentDidMount() {
    this.dispatchFetchRoutes();
  }

  componentDidUpdate(prevProps) {
    console.log("Component did update -- this.props:", this.props, "prevProps:", prevProps);
  }

  render() {
    const type = this.props.type;
    const state = this.props.routes[type];
    const queryParam = {
      [ROUTES_RECEIVED]:     "pr",
      [ROUTES_FILTERED]:     "pf",
      [ROUTES_NOT_EXPORTED]: "pn",
    }[type];
    const name = {
      [ROUTES_RECEIVED]:     "routes-received",
      [ROUTES_FILTERED]:     "routes-filtered",
      [ROUTES_NOT_EXPORTED]: "routes-not-exported",
    }[type];

    if (state.loading) {
      return null;
    }

    if (state.totalResults == 0) {
      return null;
    }

    return (
      <div className={`card routes-view ${name}`} id={name}>
        <div className="row">
          <div className="col-md-6">
            <RoutesHeader type={type} />
          </div>
          <div className="col-md-6">
            <RoutesPaginationInfo page={state.page}
                                  pageSize={state.pageSize}
                                  totalPages={state.totalPages}
                                  totalResults={state.totalResults} />

          </div>
        </div>
        <RoutesTable routes={state.routes} />

        <center>
            <RoutesPaginator page={state.page} totalPages={state.totalPages}
                             queryParam={queryParam}
                             anchor={name} />
        </center>
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
      pageSize:     state.routes.receivedPageSize,
      totalPages:   state.routes.receivedTotalPages,
      totalResults: state.routes.receivedTotalResults,
    };
    let filtered = {
      routes:       state.routes.filtered,
      loading:      state.routes.filteredLoading,
      page:         state.routes.filteredPage,
      pageSize:     state.routes.filteredPageSize,
      totalPages:   state.routes.filteredTotalPages,
      totalResults: state.routes.filteredTotalResults,
    };
    let notExported = {
      routes:       state.routes.notExported,
      loading:      state.routes.notExportedLoading,
      page:         state.routes.notExportedPage,
      pageSize:     state.routes.notExportedPageSize,
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

