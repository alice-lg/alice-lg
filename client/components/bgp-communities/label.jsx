
import React from 'react'
import {connect} from 'react-redux'


class Label extends React.Component {
  render() {
    // Lookup communities
    const readableCommunity = this.props.communities[this.props.community];
    let cls = 'label label-bgp-community ';
    if (!readableCommunity) {
      cls += "label-bgp-unknown";
      // Default label
      return (
        <span className={cls}>{this.props.community}</span>
      );
    }

    cls += "label-success ";
    // Split community into components
    return (<span className={cls}>{readableCommunity} ({this.props.community})</span>);
  }
}

export default connect(
  (state) => ({
    communities: state.config.bgp_communities,
  })
)(Label);

