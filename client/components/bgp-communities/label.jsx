
import React from 'react'
import {connect} from 'react-redux'


class Label extends React.Component {
  render() {
    console.log(this.props.communities, this.props.community);
    // Lookup communities
    const readableCommunity = this.props.communities[this.props.community];
    if (!readableCommunity) {
      return null;
    }

    let cls = 'label label-success label-bgp-community ';
    return (<span className={cls}>{readableCommunity}</span>);
  }
}

export default connect(
  (state) => ({
    communities: state.config.bgp_communities,
  })
)(Label);

