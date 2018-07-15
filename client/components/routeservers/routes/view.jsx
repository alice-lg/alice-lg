
import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'

/*
 * Render a RoutesView:
 * The routes view is a composit of:
 *  - A header
 *  - The Routes Table
 *  - A Paginator
 */

class RoutesView extends React.Component {

  componentDidMount() {


  }

  render() {

    return (
      <div className="routes-view">
        [HEADER]<br />

        [TABLE]<br />

        [Paginator]
      </div>
    );
  }

}


export default connect(
  (state) => ({
      routesFilterValue: state.routeservers.routesFilterValue
  })
)(RoutesView);

