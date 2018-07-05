import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'


import {loadRouteserverRoutes, loadRouteserverRoutesFiltered} from '../actions'
import {showBgpAttributes} from './bgp-attributes-modal-actions'

import LoadingIndicator
	from 'components/loading-indicator/small'


import FilterReason
  from 'components/routeservers/large-communities/filter-reason'

import NoexportReason
  from 'components/routeservers/large-communities/noexport-reason'


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
    const routesColumns = this.props.routesColumns;
    const routesColumnsOrder = this.props.routesColumnsOrder;

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

    let routesView = routes.map((r,i) => {
      return (
        <tr key={`${r.network}_${i}`} onClick={() => this.showAttributesModal(r)}>
          <td>
            {r.network}
            {this.props.displayReasons == "filtered" && <FilterReason route={r} />}
            {this.props.displayReasons == "noexport" && <NoexportReason route={r} />}
          </td>
          {routesColumnsOrder.map(col => <td key={col}>{_lookup(r, col)}</td>)}
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
              {routesColumnsOrder.map(col => <th key={col}>{routesColumns[col]}</th>)}
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
      filter:             state.routeservers.routesFilterValue,
      routesColumns:      state.config.routes_columns,
      routesColumnsOrder: state.config.routes_columns_order,
    }
  }
)(RoutesTable);


class RoutesTables extends React.Component {
  componentDidMount() {
    this.props.dispatch(
      loadRouteserverRoutes(this.props.routeserverId, this.props.protocolId)
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
    const noexport = this.props.noexport[this.props.protocolId] || [];

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
    const recvdHeader = mkHeader("green",  "accepted");
    const noexHeader  = mkHeader("red",    "not exported");


    return (
      <div>
        <RoutesTable header={filtdHeader} routes={filtered} displayReasons="filtered"/>
        <RoutesTable header={recvdHeader} routes={received} displayReasons={false}/>
        <RoutesTable header={noexHeader}  routes={noexport} displayReasons="noexport"/>
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
      noexport:  state.routeservers.noexport,
    }
  }
)(RoutesTables);
