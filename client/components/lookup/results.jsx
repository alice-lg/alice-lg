
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'
import {Link} from 'react-router'

import FilterReason
  from 'components/routeservers/large-communities/filter-reason'

import NoexportReason
  from 'components/routeservers/large-communities/noexport-reason'

import {showBgpAttributes}
  from 'components/routeservers/routes/bgp-attributes-modal-actions'

import BgpAttributesModal
  from 'components/routeservers/routes/bgp-attributes-modal'

import LoadingIndicator
	from 'components/loading-indicator/small'

class ResultsTableView extends React.Component {

  showAttributesModal(route) {
    this.props.dispatch(
      showBgpAttributes(route)
    );
  }

  render() {
    if (this.props.routes.length == 0) {
      return null;
    }

    const routes = this.props.routes.map((route) => (
      <tr key={route.id + '_' + route.neighbour.id + '_' + route.routeserver.id}>
        <td onClick={() => this.showAttributesModal(route)}>{route.network}
            {this.props.display_reasons == "filtered" && <FilterReason route={route} />}
        </td>
        <td onClick={() => this.showAttributesModal(route)}>{route.bgp.as_path.join(" ")}</td>
        <td onClick={() => this.showAttributesModal(route)}>
          {route.gateway}
        </td>
        <td>
          <Link to={`/routeservers/${route.routeserver.id}/protocols/${route.neighbour.id}/routes`}>
            {route.neighbour.description}
          </Link>
        </td>
        <td>
          <Link to={`/routeservers/${route.routeserver.id}/protocols/${route.neighbour.id}/routes`}>
            {route.neighbour.asn}
          </Link>
        </td>
        <td>
          <Link to={`/routeservers/${route.routeserver.id}`}>
            {route.routeserver.name}
          </Link>
        </td>
      </tr>
    ));

    return (
      <div className="card">
        {this.props.header}
        <table className="table table-striped table-routes">
          <thead>
            <tr>
              <th>Network</th>
              <th>AS Path</th>
              <th>Gateway</th>
              <th>Neighbour</th>
              <th>ASN</th>
              <th>RS</th>
            </tr>
          </thead>
          <tbody>
            {routes}
          </tbody>
        </table>
      </div>
    );
  }
}

const ResultsTable = connect()(ResultsTableView);


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

const NoResults = connect(
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
    const noexHeader  = mkHeader("red",    "not exported");

    let filteredRoutes = this.props.routes.filtered;
    let importedRoutes = this.props.routes.imported;

    return (
      <div className="lookup-results">

        <BgpAttributesModal />

        <NoResults />

        <ResultsTable header={filtdHeader}
                      routes={filteredRoutes}
                      display_reasons="filtered" />
        <ResultsTable header={recvdHeader}
                      routes={importedRoutes} />
      </div>
    )
  }

}

function selectRoutes(routes, state) {
  return _.where(routes, {state: state});
}

export default connect(
  (state) => {
    let routes = state.lookup.results;
    let filteredRoutes = selectRoutes(routes, 'filtered');
    let importedRoutes = selectRoutes(routes, 'imported');
    return {
      routes: {
        filtered: filteredRoutes,
        imported: importedRoutes
      },
    }
  }
)(LookupResults);

