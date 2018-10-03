
import React from 'react'
import {connect} from 'react-redux'


class FilterReason extends React.Component {
  render() {
    const route = this.props.route;

    if (!this.props.rejectReasons || !route || !route.bgp ||
        !route.bgp.large_communities) {
        return null;
    }

    const reason = route.bgp.large_communities.filter(elem =>
      elem[0] == this.props.asn && elem[1] == this.props.rejectId
    );
    if (!reason.length) {
      return null;
    }
    const filterReason = this.props.rejectReasons[reason[0][2]];
    const cls = `reject-reason reject-reason-${reason[0][2]}`;
    return (
      <p className={cls}>
        <a href={`http://irrexplorer.nlnog.net/search/${route.network}`}
           target="_blank" >{filterReason}</a>
      </p>
    );
  }
}

export default connect(
  state => {
    return {
      rejectReasons: state.routeservers.rejectReasons,
      asn:           state.routeservers.rejectAsn,
      rejectId:      state.routeservers.rejectId,
    }
  }
)(FilterReason);

