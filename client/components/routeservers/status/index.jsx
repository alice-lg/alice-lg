
import React from 'react'
import {connect} from 'react-redux'

import Datetime from 'components/datetime'


class Details extends React.Component {

  render() {
    let rsStatus = this.props.details[this.props.routeserverId];
    if (!rsStatus) {
      return null;
    }

    // Get routeserver name
    let rs = this.props.routeservers[parseInt(this.props.routeserverId)];
    if (!rs) {
      return null;
    }


    let lastReboot = rsStatus.last_reboot;
    if (lastReboot == "0001-01-01T00:00:00Z") {
        lastReboot = null;
    }

    return (
      <div className="routeserver-status">
        <ul>
          {lastReboot &&
            <li><i className="fa fa-clock-o"></i>
              Last Reboot: <b><Datetime value={lastReboot} /></b>
            </li>}
          <li><i className="fa fa-clock-o"></i>
            Last Reconfig: <b><Datetime value={rsStatus.last_reconfig} /></b>
          </li>
          <li><i className="fa fa-battery-full"></i>
            <b>{rsStatus.message}</b></li>
        </ul>
      </div>
    );
  }
}

export default connect(
  (state) => {
    return {
      routeservers: state.routeservers.all,
      details: state.routeservers.details
    }
  }
)(Details);

