
import React from 'react'
import {connect} from 'react-redux'

import {makeReadableCommunity} from './utils'

/*
 * Make style tags
 * Derive classes from community parts.
 */
function _makeStyleTags(community) {
  return community.map((part, i) => {
    return `label-bgp-community-${i}-${part}`;
  });
}


/*
 * Render community label
 */
class Label extends React.Component {
  render() {
    // Lookup communities
    const readableCommunity = makeReadableCommunity(
      this.props.communities,
      this.props.community);
    const key = this.props.community.join(":");

    let cls = 'label label-bgp-community ';
    if (!readableCommunity) {
      cls += "label-bgp-unknown";
      // Default label
      return (
        <span className={cls}>{key}</span>
      );
    }

    // Apply style
    cls += "label-info ";

    const styleTags = _makeStyleTags(this.props.community);
    cls += styleTags.join(" ");

    return (<span className={cls}>{readableCommunity} ({key})</span>);
  }
}

export default connect(
  (state) => ({
    communities: state.config.bgp_communities,
  })
)(Label);

