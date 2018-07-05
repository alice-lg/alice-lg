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

// Helper: Lookup value in route path
const _lookup = (r, path) => {
  const split = path.split(".").reduce((acc, elem) => acc[elem], r);
  return split;
}

const ColDefault = function(props) {
  return (
    <td onClick={props.onClick}>{_lookup(props.route, props.column)}</td>
  )
}

// Include filter and noexport reason in this column.
const ColNetwork = function(props) {
  return (
    <td onClick={props.onClick}>
      {props.route.network}
      {props.displayReasons == "filtered" && <FilterReason route={props.route} />}
      {props.displayReasons == "noexport" && <NoexportReason route={props.route} />}
    </td>
  );
}

// Special AS Path Widget
const ColAsPath = function(props) {
    const asns = _lookup(props.route, "bgp.as_path");
    const baseUrl = "http://irrexplorer.nlnog.net/search/"

    let asnLinks = asns.map((asn, i) => {
      return (<a key={`${asn}_${i}`} href={baseUrl + asn} target="_blank">{asn} </a>);
    });

    return (
        <td>
          {asnLinks}
        </td>
    );
}

const RouteColumn = function(props) {
  const widgets = {
    "network": ColNetwork,
    "bgp.as_path": ColAsPath,

    "ASPath": ColAsPath,
  };

  let Widget = widgets[props.column] || ColDefault;
  return (
    <Widget column={props.column} route={props.route}
            displayReasons={props.displayReasons}
            onClick={props.onClick} />
  );
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

    let routesView = routes.map((r,i) => {
      return (
        <tr key={`${r.network}_${i}`}>
          {routesColumnsOrder.map(col => (<RouteColumn key={col}
                                                       onClick={() => this.showAttributesModal(r)}
                                                       column={col}
                                                       route={r}
                                                       displayReasons={this.props.displayReasons} />)
          )}
        </tr>
      );
    });

    return (
      <div className="card">
        {this.props.header}
        <table className="table table-striped table-routes">
          <thead>
            <tr>
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
