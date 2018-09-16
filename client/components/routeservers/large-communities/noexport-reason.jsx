
import React from 'react'
import {connect} from 'react-redux'

class NoExportReason extends React.Component {
  render() {
    const route = this.props.route;

    if (!this.props.noexport_reasons || !route || !route.bgp ||
        !route.bgp.large_communities) {
        return null;
    }
    const reason = route.bgp.large_communities.filter(elem =>
      elem[0] == this.props.asn && elem[1] == this.props.reject_id
    );
    if (!reason.length) {
      return null;
    }
    const noexport_reason = this.props.noexport_reasons[reason[0][2]];
    const cls = `noexport-reason noexport-reason-${reason[0][2]}`;
    return (
      <p className={cls}>
        <a className={cls}
           href={`http://irrexplorer.nlnog.net/search/${route.network}`}
           target="_blank" >{noexport_reason}</a>
      </p>);
  }
}

export default connect(
  state => {
    return {
      noexport_reasons: state.routeservers.noexport_reasons,
      asn:              state.routeservers.noexport_asn,
      reject_id:        state.routeservers.noexport_id
    }
  }
)(NoExportReason);
