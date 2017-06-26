
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'

import FilterReason
  from 'components/routeservers/large-communities/filter-reason'

import NoexportReason
  from 'components/routeservers/large-communities/noexport-reason'


class ResultsTable extends React.Component {

  render() {
    if (this.props.routes.length == 0) {
      return null;
    }

    const routes = this.props.routes.map((route) => (
      <tr key={route.id + '_' + route.neighbour.id + '_' + route.routeserver.id}>
        <td>{route.network}</td>
        <td>{route.bgp.as_path.join(" ")}</td>
        <td>{route.gateway}</td>
        <td>{route.neighbour.description}</td>
        <td>{route.neighbour.asn}</td>
        <td>{route.routeserver.name}</td>
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


class LookupResults extends React.Component {

  render() {
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
        <ResultsTable header={filtdHeader} routes={filteredRoutes} />
        <ResultsTable header={recvdHeader} routes={importedRoutes} />
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
      }
    }
  }
)(LookupResults);

