
import _ from 'underscore'
import {debounce} from "underscore"

import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'
import {push, replace} from 'react-router-redux'

import Details    from '../details'
import Status     from '../status'
import PageHeader from 'components/page-header'

import {apiCacheStatus} from 'components/api-status/cache'

import ProtocolName
  from 'components/routeservers/protocols/name'


import SearchInput from 'components/search-input'

import RoutesView   from './view'
import QuickLinks   from './quick-links'
import {RelatedPeersTabs,
        RelatedPeersCard} from './related-peers'

import BgpAttributesModal
  from './bgp-attributes-modal'

import RoutesLoadingIndicator from './loading-indicator'

import {filterableColumnsText} from './utils'

import FiltersEditor from 'components/filters/editor'
import {mergeFilters} from 'components/filters/state'

import {makeLinkProps} from './urls'

// Actions
import {setFilterQueryValue, fetchRelatedPeers}
  from './actions'
import {loadRouteserverProtocol}
  from 'components/routeservers/actions'


// Constants
import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';


const makeQueryLinkProps = function(routing, query, loadNotExported) {
  // Load not exported routes flag
  const ne = loadNotExported ? 1 : 0;

  // As we need to reset the pagination, we can just
  // ommit these other parameters and just use pathname + query + ne
  return {
    pathname: routing.pathname,
    search: `?ne=${ne}&q=${query}`
  };
}


/*
 * Check if the routes view is empty, (while nothing is,
 * loading) and show info screen.
 */
const RoutesViewEmpty = (props) => {
  const hasContent = props.routes.received.totalResults > 0 ||
                     props.routes.filtered.totalResults > 0 ||
                     props.routes.notExported.totalResults > 0;
  const isLoading = props.routes.received.loading ||
                    props.routes.filtered.loading ||
                    props.routes.notExported.loading;

  if (isLoading) {
    return null; // We are not a loading indicator.
  }
 
  // Maybe this has something to do with a filter
  if (!hasContent && props.hasQuery) {
      return (
        <div className="card info-result-empty">
          <h4>No routes  matching your query.</h4>
          <p>Please check if your query is too restrictive.</p>
        </div>
      );
  }

  if (hasContent) {
    return null; // Nothing to do then.
  }

  return (
    <div className="card info-result-empty">
        <p className="card-body">There are <b>no routes</b> to display for this neighbor.</p>
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

    this.debouncedDispatch(replace(makeQueryLinkProps(
      this.props.routing, value, this.props.loadNotExported
    )));
  }

  componentDidMount() {
    // Assert neighbors for RS are loaded
    this.props.dispatch(
      loadRouteserverProtocol(this.props.params.routeserverId)
    );

    if (this.props.neighbor) {
      this.props.dispatch(
        fetchRelatedPeers(this.props.neighbor.asn)
      );
    }
  }

  componentDidUpdate(prevProps) {
    if (this.props.neighbor && this.props.neighbor != prevProps.neighbor) {
      this.props.dispatch(
        fetchRelatedPeers(this.props.neighbor.asn)
      );
    }
  }

  render() {
    let cacheStatus = apiCacheStatus(this.props.routes.received.apiStatus);
    if (this.props.anyLoading) {
      cacheStatus = null;
    }

    // We have to shift the layout a bit, to make room for
    // the related peers tabs
    let pageClass = "routeservers-page";
    if (this.props.localRelatedPeers.length > 1) {
      pageClass += " has-related-peers";
    }

    // Make placeholder for filter input
    const filterPlaceholder = "Filter by " +
      filterableColumnsText(
        this.props.routesColumns,
        this.props.routesColumnsOrder);

    return(
      <div className={pageClass}>
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
          <div className="col-main col-lg-9 col-md-12">

            <div className="card">
              <RelatedPeersTabs
                peers={this.props.localRelatedPeers}
                protocolId={this.props.params.protocolId}
                routeserverId={this.props.params.routeserverId} />
              <SearchInput
                value={this.props.filterValue}
                placeholder={filterPlaceholder}
                onChange={(e) => this.setFilter(e.target.value)}  />
            </div>

            <QuickLinks routes={this.props.routes} />

            <RoutesViewEmpty routes={this.props.routes}
                             hasQuery={!!this.props.filterValue}
                             loadNotExported={this.props.loadNotExported} />
            <RoutesView
                type={ROUTES_FILTERED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            {this.props.receivedLoading && <RoutesLoadingIndicator />}

            <RoutesView
                type={ROUTES_RECEIVED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            {this.props.notExportedLoading && <RoutesLoadingIndicator />}

            <RoutesView
                type={ROUTES_NOT_EXPORTED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />


          </div>
          <div className="col-lg-3 col-md-12 col-aside-details">
            <div className="card">
              <Status routeserverId={this.props.params.routeserverId}
                      cacheStatus={cacheStatus} />
            </div>
            <FiltersEditor makeLinkProps={makeLinkProps}
                           linkProps={this.props.linkProps}
                           filtersApplied={this.props.filtersApplied}
                           filtersAvailable={this.props.filtersAvailable} />
            <RelatedPeersCard
              neighbors={this.props.allRelatedPeers}
              rsId={this.props.params.routeserverId} 
              protocolId={this.props.params.protocolId} />
          </div>
        </div>
      </div>
    );
  }

}




export default connect(
  (state, props) => {
    const query = props.params.query; 
    const protocolId = props.params.protocolId;
    const rsId = props.params.routeserverId;
    const neighbors = state.routeservers.protocols[rsId];
    const neighbor = _.findWhere(neighbors, {id: protocolId});

    // Find related peers. Peers belonging to the same AS.
    let localRelatedPeers = [];
    if (neighbor) {
      localRelatedPeers = _.where(
        neighbors, {asn: neighbor.asn, state: "up"});
    }

    const received = {
      loading:      state.routes.receivedLoading,
      totalResults: state.routes.receivedTotalResults,
      apiStatus:    state.routes.receivedApiStatus
    };
    const filtered = {
      loading:      state.routes.filteredLoading,
      totalResults: state.routes.filteredTotalResults,
      apiStatus:    state.routes.filteredApiStatus
    };
    const notExported = {
      loading:      state.routes.notExportedLoading,
      totalResults: state.routes.notExportedTotalResults,
      apiStatus:    state.routes.notExportedApiStatus
    };

    const totalResults = state.routes.receivedTotalResults +
      state.routes.filteredTotalResults +
      state.routes.notExportedTotalResults;

    const anyLoading = state.routes.receivedLoading ||
                     state.routes.filteredLoading ||
                     state.routes.notExportedLoading;

    const filtersApplied = mergeFilters(
        state.routes.receivedFiltersApplied,
        state.routes.filteredFiltersApplied,
        state.routes.notExportedFiltersApplied
    );

    const filtersAvailable = mergeFilters(
        state.routes.receivedFiltersAvailable,
        state.routes.filteredFiltersAvailable,
        state.routes.notExportedFiltersAvailable
    );

    return({
      neighbor: neighbor,
      filterValue: state.routes.filterValue,
      routes: {
        [ROUTES_RECEIVED]:     received,
        [ROUTES_FILTERED]:     filtered,
        [ROUTES_NOT_EXPORTED]: notExported
      },
      routesColumns: state.config.routes_columns,
      routesColumnsOrder: state.config.routes_columns_order,

      routing: state.routing.locationBeforeTransitions,
      loadNotExported: state.routes.loadNotExported ||
                       !state.config.noexport_load_on_demand,

      totalResults: totalResults,
      anyLoading: anyLoading,

      filterQuery: state.routes.filterQuery,
      filtersApplied: filtersApplied,
      filtersAvailable: filtersAvailable,

      linkProps: {
        routing: state.routing.locationBeforeTransitions,

        loadNotExported: state.routes.loadNotExported,

        page:            0,
        pageReceived:    0, // Reset pagination on filter change
        pageFiltered:    0,
        pageNotExported: 0,

        query: state.routes.filterValue,

        filtersApplied: filtersApplied,
      },

      localRelatedPeers: localRelatedPeers,
      allRelatedPeers: state.routes.allRelatedPeers,

      // Loding indicator helper
      receivedLoading:    state.routes.receivedLoading,
      filteredLoading:    state.routes.filteredLoading,
      notExportedLoading: state.routes.notExportedLoading
    });
  }
)(RoutesPage);

