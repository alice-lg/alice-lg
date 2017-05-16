
/*
 * Bird status
 */

import React from 'react'
import {connect} from 'react-redux'

// Actions
import {loadRouteserverStatus}
  from 'components/routeservers/actions'

class Status extends React.Component {
  componentDidMount() {
    this.props.dispatch(
      loadRouteserverStatus(this.props.routeserverId)
    );
  }

  render() {
    let rsStatus = this.props.details[this.props.routeserverId];
    if (!rsStatus) {
      return null;
    }

    return (
      <div className="routeserver-status">
        <div className="bird-version">
          Bird {rsStatus.version}
        </div>
      </div>
    );
  }
}

export default connect(
  (state) => {
    return {
      details: state.routeservers.details
    }
  }
)(Status);


