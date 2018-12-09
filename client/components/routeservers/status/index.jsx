
import React from 'react'
import {connect} from 'react-redux'

import Datetime from 'components/datetime'
import RelativeTime from 'components/datetime/relative'


class Details extends React.Component {

  render() {

    let rsStatus = this.props.details[this.props.routeserverId];
    if (!rsStatus) {
      return null;
    }

    // Get routeserver name
    let rs = this.props.routeservers[this.props.routeserverId];
    if (!rs) {
      return null;
    }

    let lastReboot = rsStatus.last_reboot;
    if (lastReboot == "0001-01-01T00:00:00Z") {
        lastReboot = null;
    }

    let cacheStatus = null;
    if (this.props.cacheStatus && 
        this.props.cacheStatus.ttl &&
        this.props.cacheStatus.generatedAt) {
      const cs = this.props.cacheStatus;
      cacheStatus = [
         <tr key="cache-status-cached-at">
           <td><i className="fa fa-refresh"></i></td>
           <td>
             Generated <b><RelativeTime value={cs.generatedAt}
                                        fuzzyNow={5}
                                        pastEvent={true} /></b>.<br />
             Next refresh <b><RelativeTime futureEvent={true}
                                           fuzzyNow={5}
                                           value={cs.ttl} /></b>.
           </td>
         </tr>,
      ];
    };

    return (
      <table className="routeserver-status">
        <tbody>
        {lastReboot &&
          <tr>
            <td><i className="fa fa-clock-o"></i></td>
            <td>Last Reboot: <b><Datetime value={lastReboot} /></b></td>
          </tr>}
        <tr>
          <td><i className="fa fa-clock-o"></i></td>
          <td>Last Reconfig: <b><Datetime value={rsStatus.last_reconfig} /></b></td>
        </tr>

        <tr>
          <td><i className="fa fa-thumbs-up"></i></td>
          <td><b>{rsStatus.message}</b></td>
        </tr>

        {cacheStatus}
        </tbody>
      </table>
    );
  }
}

export default connect(
  (state) => {
    return {
      routeservers: state.routeservers.byId,
      details: state.routeservers.details
    }
  }
)(Details);

