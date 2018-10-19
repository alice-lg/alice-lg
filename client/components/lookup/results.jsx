
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'
import {Link} from 'react-router'
import {replace} from 'react-router-redux'

import {filtersEqual} from './filter-groups'

import FilterReason
  from 'components/routeservers/communities/filter-reason'

import NoexportReason
  from 'components/routeservers/communities/noexport-reason'

import BgpAttributesModal
  from 'components/routeservers/routes/bgp-attributes-modal'

import LoadingIndicator
	from 'components/loading-indicator/small'

import ResultsTable from './table'

import {loadResults, reset} from './actions'

import {RoutesPaginator,
        RoutesPaginationInfo} from './pagination'

import {RoutesHeader}
  from 'components/routeservers/routes/view'



const ResultsView = function(props) {
  if(!props.routes) {
    return null;
  }
  if(props.routes.length == 0) {
    return null;
  }

  const type = props.type;

  return (
    <div className={`card routes-view routes-${type}`}>
      <div className="row">
        <div className="col-md-6 routes-header-container">
          <RoutesHeader type={type} />
        </div>
        <div className="col-md-6">
          <RoutesPaginationInfo page={props.page}
                                pageSize={props.pageSize}
                                totalPages={props.totalPages}
                                totalResults={props.totalResults} />
         </div>
      </div>
      <ResultsTable routes={props.routes}
                    displayReasons={props.displayReasons} />
      <center>
        <RoutesPaginator page={props.page} totalPages={props.totalPages}
                         queryParam={props.query}
                         anchor={type} />
      </center>
    </div>
  );
}

class NoResultsView extends React.Component {
  render() {
    if (!this.props.show) {
      return null;
    }
    return (
      <p className="lookup-no-results text-info card">
        No prefixes could be found for <b>{this.props.query}</b>
      </p>
    );
  }
}

const NoResultsFallback = connect(
  (state) => {
    let total = state.lookup.totalRoutes;
    let query = state.lookup.query;
    let isLoading = state.lookup.isLoading;

    let show = false;

    if (total == 0 && query != "" && isLoading == false) {
      show = true;
    }

    return {
      show: show,
      query: state.lookup.query
    }
  }
)(NoResultsView);



class LookupResults extends React.Component {

  dispatchLookup() {
    const query = this.props.query;
    const pageImported = this.props.pagination.imported.page;
    const pageFiltered = this.props.pagination.filtered.page;
    const filters = this.props.filtersApplied;

    if (query == "") {
      // Dispatch reset and transition to main page
      this.props.dispatch(reset());
      this.props.dispatch(replace("/"));
    } else {
      this.props.dispatch(
        loadResults(query, filters, pageImported, pageFiltered)
      );
    }
  }

  componentDidMount() {
    // Dispatch query
    this.dispatchLookup();
  }

  componentDidUpdate(prevProps) {
    console.log("filters equal:", filtersEqual(this.props.filtersApplied, prevProps.filtersApplied));

    if(this.props.query != prevProps.query ||
       this.props.pagination.filtered.page != prevProps.pagination.filtered.page ||
       this.props.pagination.imported.page != prevProps.pagination.imported.page ||
       !filtersEqual(this.props.filtersApplied, prevProps.filtersApplied)) {
        this.dispatchLookup();
    }
  }

  render() {
    if(this.props.isLoading) {
      return <LoadingIndicator />;
    }

    const ref = this.refs[this.props.anchor];
    if(ref) {
      ref.scrollIntoView();
    }

    const filteredRoutes = this.props.routes.filtered;
    const importedRoutes = this.props.routes.imported;

    return (
      <div className="lookup-results">
        <BgpAttributesModal />

        <NoResultsFallback />

        <a ref="filtered" name="routes-filtered" />
        <ResultsView type="filtered"
                     routes={filteredRoutes}

                     page={this.props.pagination.filtered.page}
                     pageSize={this.props.pagination.filtered.pageSize}
                     totalPages={this.props.pagination.filtered.totalPages}
                     totalResults={this.props.pagination.filtered.totalResults}

                     query={this.props.query}

                     displayReasons="filtered" />

        <a ref="received" name="routes-received" />
        <ResultsView type="received"

                     page={this.props.pagination.imported.page}
                     pageSize={this.props.pagination.imported.pageSize}
                     totalPages={this.props.pagination.imported.totalPages}
                     totalResults={this.props.pagination.imported.totalResults}

                     query={this.props.query}

                     routes={importedRoutes} />
      </div>
    );
  }
}

export default connect(
  (state) => {
    const filteredRoutes = state.lookup.routesFiltered;
    const importedRoutes = state.lookup.routesImported;

    return {
      anchor: state.lookup.anchor,
      routes: {
        filtered: filteredRoutes,
        imported: importedRoutes
      },
      pagination: {
        filtered: {
          page: state.lookup.pageFiltered,
          pageSize: state.lookup.pageSizeFiltered,
          totalPages: state.lookup.totalPagesFiltered,
          totalResults: state.lookup.totalRoutesFiltered,
        },
        imported: {
          page: state.lookup.pageImported,
          pageSize: state.lookup.pageSizeImported,
          totalPages: state.lookup.totalPagesImported,
          totalResults: state.lookup.totalRoutesImported,
        }
      },
      isLoading: state.lookup.isLoading,
      query: state.lookup.query,
      filtersApplied: state.lookup.filtersApplied,
    }
  }
)(LookupResults);

