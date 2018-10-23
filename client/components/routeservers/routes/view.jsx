
import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'

import RoutesTable from './table'
import {RoutesPaginator,
        RoutesPaginationInfo} from './pagination'

import {fetchRoutesReceived,
        fetchRoutesFiltered,
        fetchRoutesNotExported} from './actions'

import {makeLinkProps} from './urls'

import {filtersEqual} from 'components/filters/groups'

// Constants
import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';


export const RoutesHeader = (props) => {
  const type = props.type;
  const rtype = {
    [ROUTES_RECEIVED]: "accepted",
    [ROUTES_FILTERED]: "filtered",
    [ROUTES_NOT_EXPORTED]: "not exported"
  }[type];
  let stype = type;
  if (stype == ROUTES_NOT_EXPORTED) {
    stype = 'not-exported'; // *sigh*
  }
  const cls = `card-header card-header-routes ${stype}`;
  return (<p className={cls}>Routes {rtype}</p>);
};

/*
 * Render a RoutesView:
 * The routes view is a composit of:
 *  - A header
 *  - The Routes Table
 *  - A Paginator
 */

class RoutesView extends React.Component {

  constructor(props) {
    super(props);
  }

  dispatchFetchRoutes() {
    const type = this.props.type;

    // Get filters
    const filters = this.props.filtersApplied;

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

    // Handle special case, when on demand loading is enabled,
    // we defer this dispatch, until an user interaction.
    if (!params.loadRoutes) {
      return; // We are done here.
    }

    // Otherwise, just dispatch the request:
    this.props.dispatch(fetchRoutes(rsId, pId, params.page, query, filters));
  }

  /*
   * Diff props and this.props to check if we need to
   * dispatch another fetch routes
   */
  routesNeedFetch(props) {
    const type = this.props.type;
    const nextParams = this.props.routes[type];
    const params = props.routes[type]; // Previous props

    if (this.props.filterQuery != props.filterQuery || // Query
        params.page != nextParams.page || // Pagination
        !filtersEqual(this.props.filtersApplied, props.filtersApplied) || // Filters
        params.loadRoutes != nextParams.loadRoutes || // Defered loading
        props.protocolId != this.props.protocolId // Switch related peers
        ) {
          return true;
    }

    return false;
  }

  componentDidMount() {
    this.dispatchFetchRoutes();
  }

  componentDidUpdate(prevProps) {
    const scrollAnchor = this.refs.scrollAnchor;

    if (this.routesNeedFetch(prevProps)) {
      this.dispatchFetchRoutes();

      if (scrollAnchor) {
        scrollAnchor.scrollIntoView({
          behaviour: "smooth",
          block: "start",
        });
      }
    }
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

    if (!state.loadRoutes) {
      // In case it was not yet requested, render a trigger
      // and defer routesFetching until a user interaction has
      // occured.
      return this.renderLoadTrigger();
    }

    if (state.loading) {
      return null;
    }

    if (state.totalResults == 0) {
      return null;
    }

    // Render the routes card
    return (
      <div className={`card routes-view ${name}`}>
        <div className="row">
          <div className="col-md-6 routes-header-container">
            <a name={name} id={name} ref="scrollAnchor">
              <RoutesHeader type={type} />
            </a>
          </div>
          <div className="col-md-6">
            <RoutesPaginationInfo page={state.page}
                                  pageSize={state.pageSize}
                                  totalPages={state.totalPages}
                                  totalResults={state.totalResults} />
           </div>
        </div>
        <RoutesTable type={type}
                     routes={state.routes}
                     routeserverId={this.props.routeserverId} />
        <center>
          <RoutesPaginator page={state.page} totalPages={state.totalPages}
                           queryParam={queryParam}
                           anchor={name} />
        </center>
      </div>
    );
  }

  renderLoadTrigger() {
    const type = this.props.type;
    const state = this.props.routes[type];
    const name = {
      [ROUTES_RECEIVED]:     "routes-received",
      [ROUTES_FILTERED]:     "routes-filtered",
      [ROUTES_NOT_EXPORTED]: "routes-not-exported",
    }[type];


    // This is an artificial delay, to make the user wait until
    // filtered and recieved routes are fetched
    if (!state.otherLoaded) {
      return null;
    }

    const linkProps = makeLinkProps({
      loadNotExported: true,

      anchor: "routes-not-exported",
      page: this.props.routes.notExported.page,

      pageReceived:    this.props.routes.received.page,
      pageFiltered:    this.props.routes.filtered.page,
      pageNotExported: this.props.routes.notExported.page,

      routing: this.props.routing
    });

    return (
      <div className={`card routes-view ${name}`}>
        <div className="row">
          <div className="col-md-6">
            <a name={name} id={name} ref="scrollAnchor">
              <RoutesHeader type={type} />
            </a>
          </div>
        </div>
        <p className="help">
          Due to the potentially high amount of routes not exported,
          they are only fetched on demand.
        </p>

        <Link to={linkProps} className="btn btn-block btn-danger">
           Load Routes Not Exported
        </Link>
      </div>
    );
  }
}

export default connect(
  (state) => {
    const received = {
      routes:       state.routes.received,
      requested:    state.routes.receivedRequested,
      loading:      state.routes.receivedLoading,
      page:         state.routes.receivedPage,
      pageSize:     state.routes.receivedPageSize,
      totalPages:   state.routes.receivedTotalPages,
      totalResults: state.routes.receivedTotalResults,
      loadRoutes:   true,
    };
    const filtered = {
      routes:       state.routes.filtered,
      loading:      state.routes.filteredLoading,
      requested:    state.routes.filteredRequested,
      page:         state.routes.filteredPage,
      pageSize:     state.routes.filteredPageSize,
      totalPages:   state.routes.filteredTotalPages,
      totalResults: state.routes.filteredTotalResults,
      loadRoutes:   true,
    };
    const notExported = {
      routes:       state.routes.notExported,
      requested:    state.routes.notExportedRequested,
      loading:      state.routes.notExportedLoading,
      page:         state.routes.notExportedPage,
      pageSize:     state.routes.notExportedPageSize,
      totalPages:   state.routes.notExportedTotalPages,
      totalResults: state.routes.notExportedTotalResults,

      loadRoutes:    state.routes.loadNotExported ||
                     !state.config.noexport_load_on_demand,

      otherLoaded:  state.routes.receivedRequested &&
                    !state.routes.receivedLoading  &&
                    state.routes.filteredRequested &&
                    !state.routes.filteredLoading
    };
    return({
      filterQuery: state.routes.filterQuery,
      routes: {
          [ROUTES_RECEIVED]:     received,
          [ROUTES_FILTERED]:     filtered,
          [ROUTES_NOT_EXPORTED]: notExported
      },
      filtersApplied: state.routes.filtersApplied,
      routing: state.routing.locationBeforeTransitions
    });
  }
)(RoutesView);

