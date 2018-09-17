import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'


import {loadRouteserverRoutes, loadRouteserverRoutesFiltered} from '../actions'
import {showBgpAttributes} from './bgp-attributes-modal-actions'

import LoadingIndicator from './loading-indicator'

import RouteColumn from './column'


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
    const blackholes = this.props.blackholes;

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
                                                       blackholes={blackholes}
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
  (state, props) => {
    const rsId = parseInt(props.routeserverId, 10);
    const blackholes = state.config.blackholes[rsId];
    return {
      blackholes:         blackholes,
      routesColumns:      state.config.routes_columns,
      routesColumnsOrder: state.config.routes_columns_order,
    }
  }
)(RoutesTable);

