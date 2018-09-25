
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'
import {Link} from 'react-router'
import {replace} from 'react-router-redux'

import FilterReason
  from 'components/routeservers/large-communities/filter-reason'

import NoexportReason
  from 'components/routeservers/large-communities/noexport-reason'

import BgpAttributesModal
  from 'components/routeservers/routes/bgp-attributes-modal'

import LoadingIndicator
	from 'components/loading-indicator/small'

import ResultsTable from './table'

import {loadResults, reset} from './actions'


const ResultsView = function(props) {
  if(props.routes.length == 0) {
    return null;
  }

  return (
    <div className="card">
      {props.header}
      <ResultsTable routes={props.routes}
                    displayReasons={props.displayReasons} />
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
    let total = state.lookup.results;
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

  dispatchLookup(query) {
    if (query == "") {
      // Dispatch reset and transition to main page
      this.props.dispatch(reset());
      this.props.dispatch(replace("/"));
    } else {
      this.props.dispatch(
        loadResults(query)
      );
    }
  }

  componentDidMount() {
    // Dispatch query
    this.dispatchLookup(this.props.query);
  }

  componentDidUpdate(prevProps) {
    if(this.props.query != prevProps.query) {
      this.dispatchLookup(this.props.query);
    }
  }

  render() {
    if(this.props.isLoading) {
      return (
				<LoadingIndicator />
      );
    }

    const mkHeader = (color, action) => (
        <p style={{"color": color, "textTransform": "uppercase"}}>
          Routes {action}
        </p>
    );

    const filtdHeader = mkHeader("orange", "filtered");
    const recvdHeader = mkHeader("green",  "accepted");

    const filteredRoutes = this.props.routes.filtered;
    const importedRoutes = this.props.routes.imported;

    return (
      <div className="lookup-results">
        <BgpAttributesModal />

        <NoResultsFallback />

        <ResultsView header={filtdHeader}
                     routes={filteredRoutes}
                     displayReasons="filtered" />

        <ResultsView header={recvdHeader}
                     routes={importedRoutes} />
      </div>
    )
  }

}

export default connect(
  (state) => {
    const filteredRoutes = state.lookup.routesFiltered;
    const importedRoutes = state.lookup.routesImported; 

    return {
      routes: {
        filtered: filteredRoutes,
        imported: importedRoutes
      },
      pagination: {
        filtered: {
          page: state.lookup.pageFiltered,
          totalPages: state.lookup.totalPagesFiltered,
        },
        imported: {
          page: state.lookup.pageImported,
          totalPages: state.lookup.totalPagesImported,
        }
      },
      isLoading: state.lookup.isLoading,
      query: state.lookup.query,
    }
  }
)(LookupResults);

