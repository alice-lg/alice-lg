import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'


import {loadRouteserverRoutes, loadRouteserverRoutesFiltered} from '../actions'
import {showBgpAttributes} from './bgp-attributes-modal-actions'

import LoadingIndicator
	from './loading-indicator'

import PrimaryIndicator
  from './primary-indicator'

import FilterReason
  from 'components/routeservers/large-communities/filter-reason'

import NoexportReason
  from 'components/routeservers/large-communities/noexport-reason'

import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';

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



/*
 * Rendering Components
 * ====================
 */

const ColDefault = function(props) {
  return (
    <td>
      <span onClick={props.onClick}>{_lookup(props.route, props.column)}</span>
    </td>
  )
}

// Include filter and noexport reason in this column.
const ColNetwork = function(props) {
  return (
    <td className="col-route-network">
      <span className="route-network" onClick={props.onClick}>
        <PrimaryIndicator route={props.route} />
        {props.route.network}
      </span>
      {props.displayReasons == ROUTES_FILTERED && <FilterReason route={props.route} />}
      {props.displayReasons == ROUTES_NOT_EXPORTED && <NoexportReason route={props.route} />}
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
                                                       displayReasons={this.props.type} />)
          )}
        </tr>
      );
    });

    return (
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
    );
  }
}

export default connect(
  (state) => ({
    routesColumns:      state.config.routes_columns,
    routesColumnsOrder: state.config.routes_columns_order,
  })
)(RoutesTable);

