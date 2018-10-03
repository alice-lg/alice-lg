
import React from 'react'
import {connect} from 'react-redux'

import {resolveCommunities} from './utils'

class FilterReason extends React.Component {
  render() {
    const route = this.props.route;

    if (!this.props.rejectReasons || !route || !route.bgp ||
        !route.bgp.large_communities) {
        return null;
    }

    const reasons = resolveCommunities(
      this.props.rejectReasons, route.bgp.large_communities,
    );

    const reasonsView = reasons.map(([community, reason], key) => {
      const cls = `reject-reason reject-reason-${community[2]}`;
      return (
        <p key={key} className={cls}>
          <a href={`http://irrexplorer.nlnog.net/search/${route.network}`}
             target="_blank" >{reason}</a>
        </p>
      );
    });

    return (<div className="reject-reasons">{reasonsView}</div>);
  }
}

export default connect(
  state => {
    return {
      rejectReasons: state.routeservers.rejectReasons,
    }
  }
)(FilterReason);

