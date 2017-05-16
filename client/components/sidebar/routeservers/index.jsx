
/**
 * Routeservers List component
 */


import React from 'react'
import { connect } from 'react-redux'

import{ push } from 'react-router-redux'

import { loadRouteservers } from 'components/routeservers/actions'

// Components 
import Status from './status'


class RouteserversList extends React.Component {

  componentDidMount() {
    this.props.dispatch(
      loadRouteservers()
    );
  }

  showRouteserver(id) {
    this.props.dispatch(
      push(`/routeservers/${id}`)
    );
  }

  render() {
    let routeservers = this.props.routeservers.map((rs) =>
      <li key={rs.id} onClick={() => this.showRouteserver(rs.id)}>
        <span className="routeserver-id">{rs.name}</span>
        <Status routeserverId={rs.id} />
      </li>
    );

    return (
      <div className="routeservers-list">
        <h2>Routeservers</h2>
        <ul> 
          {routeservers}
        </ul>
      </div>
    );
  }
}


export default connect(
  (state) => {
    return {
      routeservers: state.routeservers.all
    };
  }
)(RouteserversList);


