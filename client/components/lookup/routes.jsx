import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'


import {loadRouteserverRoutes, loadRouteserverRoutesFiltered} from '../actions'
import {showBgpAttributes} from './bgp-attributes-modal-actions'

import LoadingIndicator
	from 'components/loading-indicator/small'


class FilterReason extends React.Component {
  render() {
    const route = this.props.route;

    if (!this.props.reject_reasons || !route || !route.bgp ||
        !route.bgp.large_communities) {
        return null;
    }

    const reason = route.bgp.large_communities.filter(elem =>
      elem[0] == this.props.asn && elem[1] == this.props.reject_id
    );
    if (!reason.length) {
      return null;
    }

    return <p className="reject-reason">{this.props.reject_reasons[reason[0][2]]}</p>;
  }
}

FilterReason = connect(
  state => {
    return {
      reject_reasons: state.routeservers.reject_reasons,
      asn:            state.routeservers.asn,
      reject_id:      state.routeservers.reject_id,
    }
  }
)(FilterReason);


function _filteredRoutes(routes, filter) {
  let filtered = [];
  if(filter == "") {
    return routes; // nothing to do here
  }

  filter = filter.toLowerCase();

  // Filter protocols
  filtered = _.filter(routes, (r) => {
    return (r.network.toLowerCase().indexOf(filter) != -1 ||
            r.gateway.toLowerCase().indexOf(filter) != -1 ||
            r.interface.toLowerCase().indexOf(filter) != -1);
  });

  return filtered;
}

class RoutesTable extends React.Component {
  showAttributesModal(route) {
    this.props.dispatch(
      showBgpAttributes(route)
    );
  }


  render() {
    let routes = this.props.routes;
    const routes_columns = this.props.routes_columns;

    routes = _filteredRoutes(routes, this.props.filter);
    if (!routes || !routes.length) {
      return null;
    }

    const _lookup = (r, path) => {
      const split = path.split(".").reduce((acc, elem) => acc[elem], r);

      if (Array.isArray(split)) {
        return split.join(" ");
      }
      return split;
    }

    let routesView = routes.map((r) => {
      return (
        <tr key={r.network} onClick={() => this.showAttributesModal(r)}>
          <td>{r.network}{this.props.display_filter && <FilterReason route={r}/>}</td>
          {Object.keys(routes_columns).map(col => <td key={col}>{_lookup(r, col)}</td>)}
        </tr>
      );
    });

    return (
      <div className="card">
        {this.props.header}
        <table className="table table-striped table-routes">
          <thead>
            <tr>
              <th>Network</th>
              {Object.values(routes_columns).map(col => <th key={col}>{col}</th>)}
            </tr>
          </thead>
          <tbody>
            {routesView}
          </tbody>
        </table>
      </div>
    );
  }
}


RoutesTable = connect(
  (state) => {
    return {
      filter:         state.routeservers.routesFilterValue,
      reject_reasons: state.routeservers.reject_reasons,
      routes_columns: state.config.routes_columns,
    }
  }
)(RoutesTable);


class RoutesTables extends React.Component {
  componentDidMount() {
    this.props.dispatch(
      loadRouteserverRoutes(this.props.routeserverId, this.props.protocolId)
    );
    this.props.dispatch(
      loadRouteserverRoutesFiltered(this.props.routeserverId,
                                    this.props.protocolId)
    );
  }

  render() {
    if(this.props.isLoading) {
      return (
				<LoadingIndicator />
      );
    }

    const routes = this.props.routes[this.props.protocolId];
    const filtered = this.props.filtered[this.props.protocolId] || [];

    if((!routes || routes.length == 0) &&
			 (!filtered || filtered.length == 0)) {
      return(
        <p className="help-block">
          No routes matched your filter.
        </p>
      );
    }


    const received = routes.filter(r => filtered.indexOf(r) < 0);

    const mkHeader = (color, action) => (
        <p style={{"color": color, "textTransform": "uppercase"}}>
          Routes {action}
        </p>
    );

    const filtdHeader = mkHeader("orange", "filtered");
    const recvdHeader = mkHeader("green", "accepted");


    return (
      <div>
        <RoutesTable header={filtdHeader} routes={filtered} display_filter={true}/>
        <RoutesTable header={recvdHeader} routes={received} display_filter={false}/>
      </div>
    );

  }
}


export default connect(
  (state) => {
    return {
      isLoading: state.routeservers.routesAreLoading,
      routes:    state.routeservers.routes,
      filtered:  state.routeservers.filtered,
    }
  }
)(RoutesTables);
