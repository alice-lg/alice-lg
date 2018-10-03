
import React from 'react'
import {connect} from 'react-redux'

import {resolveCommunity} from './utils'



/*
 * Expand variables in string:
 *    "Test AS$0 rejects $2"
 * will expand with [23, 42, 123] to
 *    "Test AS23 rejects 123"
 */
function _expandVars(str, vars) {
  if (!str) {
    return str; // We don't have to do anything.
  }

  var res = str;
  vars.map((v, i) => {
    res = res.replace(`$${i}`, v); 
  });

  return res;
}

/*
 * Make style tags
 * Derive classes from community parts.
 */
function _makeStyleTags(community) {
  return community.map((part, i) => {
    return `label-bgp-community-${i}-${part}`;
  });
}


class Label extends React.Component {
  render() {

    // Lookup communities
    const readableCommunityLabel = resolveCommunity(this.props.communities, this.props.community);
    const readableCommunity = _expandVars(readableCommunityLabel, this.props.community);
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

