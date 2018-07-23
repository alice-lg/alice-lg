
import {debounce} from "underscore"

import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'
import {push} from 'react-router-redux'

import Details    from '../details'
import Status     from '../status'
import PageHeader from 'components/page-header'

import ProtocolName
  from 'components/routeservers/protocols/name'

import RoutesView  from './view'

import SearchInput from 'components/search-input'

import BgpAttributesModal
  from './bgp-attributes-modal'


import RoutesLoadingIndicator from './loading-indicator'

// Actions
import {setFilterQueryValue}
  from './actions'
import {loadRouteserverProtocol}
  from 'components/routeservers/actions'


// Constants
import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';


const makeQueryLinkProps = function(routing, query) {
  // As we need to reset the pagination, we can just
  // ommit these parameters and just use pathname + query
  return {
    pathname: routing.pathname,
    search: `?q=${query}`
  };
}


/*
 * Check if the routes view is empty, (while nothing is,
 * loading) and show info screen.
 */
const RoutesViewEmpty = (props) => {
  const isLoading = props.routes.received.loading ||
                    props.routes.filtered.loading ||
                    props.routes.notExported.loading;

  if (isLoading) {
    return null; // We are not a loading indicator.
  }
  
  const hasContent = props.routes.received.totalResults > 0 ||
                     props.routes.filtered.totalResults > 0 ||
                     props.routes.notExported.totalResults > 0;
  if (hasContent) {
    return null; // Nothing to do then.
  }


  // Show info screen
  return (
    <div className="card info-result-empty">
      <h4>No routes found matching your query.</h4>
      <p>Please check if your query is too restrictive.</p>
    </div>
  );
}


class RoutesPage extends React.Component {
  constructor(props) {
    super(props);
    
    // Create debounced dispatch, as we don't want to flood
    // the server with API queries
    this.debouncedDispatch = debounce(this.props.dispatch, 350);
  }


  setFilter(value) {
    this.props.dispatch(
      setFilterQueryValue(value)
    );

    this.debouncedDispatch(push(makeQueryLinkProps(
      this.props.routing, value
    )));
  }

  componentDidMount() {
    // Assert neighbors for RS are loaded
    this.props.dispatch(
      loadRouteserverProtocol(parseInt(this.props.params.routeserverId))
    );
  }

  render() {
    console.log("render props filter value:", this.props.filterValue);
    return(
      <div className="routeservers-page">
        <PageHeader>
          <Link to={`/routeservers/${this.props.params.routeserverId}`}>
            <Details routeserverId={this.props.params.routeserverId} />
          </Link>
          <span className="spacer">&raquo;</span>
          <ProtocolName routeserverId={this.props.params.routeserverId}
                        protocolId={this.props.params.protocolId} />
        </PageHeader>

        <BgpAttributesModal />

        <div className="row details-main">
          <div className="col-md-8">

            <div className="card">
              <SearchInput
                value={this.props.filterValue}
                placeholder="Filter by Network or BGP next-hop"
                onChange={(e) => this.setFilter(e.target.value)}  />
            </div>

            <RoutesViewEmpty routes={this.props.routes} />

            <RoutesView
                type={ROUTES_FILTERED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            <RoutesView
                type={ROUTES_RECEIVED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            <RoutesView
                type={ROUTES_NOT_EXPORTED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            <RoutesLoadingIndicator />

          </div>
          <div className="col-md-4">
            <div className="card">
              <Status routeserverId={this.props.params.routeserverId} />
            </div>
          </div>
        </div>
      </div>
    );
  }

}


export default connect(
  (state) => {
    let received = {
      loading:      state.routes.receivedLoading,
      totalResults: state.routes.receivedTotalResults,
    };
    let filtered = {
      loading:      state.routes.filteredLoading,
      totalResults: state.routes.filteredTotalResults,
    };
    let notExported = {
      loading:      state.routes.notExportedLoading,
      totalResults: state.routes.notExportedTotalResults,
    };
    return({
      filterValue: state.routes.filterValue,
      routes: {
          [ROUTES_RECEIVED]:     received,
          [ROUTES_FILTERED]:     filtered,
          [ROUTES_NOT_EXPORTED]: notExported
      },
      routing: state.routing.locationBeforeTransitions
    });
  }
)(RoutesPage);

