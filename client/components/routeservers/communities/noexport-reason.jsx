
import React from 'react'
import {connect} from 'react-redux'

import {resolveCommunities} from './utils'

class NoExportReason extends React.Component {
  render() {
    const route = this.props.route;
  
    console.log(this.props.noexportReasons);

    if (!this.props.noexportReasons || !route || !route.bgp ||
        !route.bgp.large_communities) {
        return null;
    }

    const reasons = resolveCommunities(
      this.props.noexportReasons, route.bgp.large_communities,
    );

    const reasonsView = reasons.map(([community, reason], key) => {
      const cls = `noexport-reason noexport-reason-${community[2]}`;
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
      noexportReasons: state.routeservers.noexportReasons,
    }
  }
)(NoExportReason);

