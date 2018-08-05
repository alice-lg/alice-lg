
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
    let statusInfo = [];

    let rsStatus = this.props.details[this.props.routeserverId];
    if (rsStatus) {
      statusInfo.push(
        <div className="bird-version" key="status-version">
          Bird {rsStatus.version}
        </div>
      );
    }

    // Check for errors
    let rsError = this.props.errors[this.props.routeserverId];
    if (rsError) {
      if (rsError.code >= 100 && rsError.code < 200) {
        statusInfo.push(
          <div className="api-error" key="status-error">
            Unreachable
          </div>
        );
      } else {
        statusInfo.push(
          <div className="api-error" key="status-error">
            {rsError.tag}
          </div>
        );
      }
    }

    return (
      <div className="routeserver-status">
        {statusInfo}
      </div>
    );
  }
}

export default connect(
  (state) => {
    return {
      details: state.routeservers.details,
      errors: state.routeservers.statusErrors
    }
  }
)(Status);


