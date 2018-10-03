
import React from 'react'
import {connect} from 'react-redux'

class NoExportReason extends React.Component {
  render() {
    const route = this.props.route;

    if (!this.props.noexportReasons || !route || !route.bgp ||
        !route.bgp.large_communities) {
        return null;
    }
    const reason = route.bgp.large_communities.filter(elem =>
      elem[0] == this.props.asn && elem[1] == this.props.rejectId
    );
    if (!reason.length) {
      return null;
    }
    const noexportReason = this.props.noexportReasons[reason[0][2]];
    const cls = `noexport-reason noexport-reason-${reason[0][2]}`;
    return (
      <p className={cls}>
        <a href={`http://irrexplorer.nlnog.net/search/${route.network}`}
           target="_blank" >{noexportReason}</a>
      </p>);
  }
}

export default connect(
  state => {
    return {
      noexportReasons: state.routeservers.noexportReasons,

      asn:             state.routeservers.noexportAsn,
      rejectId:        state.routeservers.noexportId
    }
  }
)(NoExportReason);
