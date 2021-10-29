
/*
 * Lookup Results Table
 * --------------------
 */

import React from 'react'
import {connect} from 'react-redux'
import {Link} from 'react-router'
import {push} from 'react-router-redux'


import {_lookup,
        ColDefault,
        ColNetwork,
        ColFlags,
        ColAsPath} from 'components/routeservers/routes/route/column'

import {showBgpAttributes}
  from 'components/routeservers/routes/bgp-attributes-modal-actions'


// Link Wrappers:
const ColLinkedNeighbor = function(props) {
  const route = props.route;
  const to = `/routeservers/${route.routeserver.id}/protocols/${route.neighbor.id}/routes`;
  
  return (
    <td>
      <Link to={to}>{_lookup(props.route, props.column)}</Link>
    </td>
  );
}

const ColLinkedRouteserver = function(props) {
  const route = props.route;
  const to = `/routeservers/${route.routeserver.id}`;
  
  return (
    <td>
      <Link to={to}>{_lookup(props.route, props.column)}</Link>
    </td>
  );
}


// Custom RouteColumn
const RouteColumn = function(props) {
  const widgets = {
    "network": ColNetwork,

    "flags": ColFlags,

    "bgp.as_path": ColAsPath,
    "ASPath": ColAsPath,

    "neighbor.description": ColLinkedNeighbor,
    "neighbor.asn": ColLinkedNeighbor,
    
    "routeserver.name": ColLinkedRouteserver
  };

  const rsId = props.route.routeserver.id;
  const blackholes = props.blackholesMap[rsId] || [];

  let Widget = widgets[props.column] || ColDefault;
  return (
    <Widget column={props.column} route={props.route}
            displayReasons={props.displayReasons}
            blackholes={blackholes}
            onClick={props.onClick} />
  );
}


class LookupRoutesTable extends React.Component {
  showAttributesModal(route) {
    this.props.dispatch(showBgpAttributes(route));
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
        <tr key={i}>
          {routesColumnsOrder.map(col => {
            return (<RouteColumn key={col}
                                 onClick={() => this.showAttributesModal(r)}
                                 blackholesMap={this.props.blackholesMap}
                                 column={col}
                                 route={r}
                                 displayReasons={this.props.displayReasons} />);
            }
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
    blackholesMap:      state.config.blackholes,
    routesColumns:      state.config.lookup_columns,
    routesColumnsOrder: state.config.lookup_columns_order,
  })
)(LookupRoutesTable);


