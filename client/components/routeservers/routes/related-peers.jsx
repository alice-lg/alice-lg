
import React from 'react'

import {Link} from 'react-router'

import {makePeerLinkProps} from './urls'

/*
 * Render related peers as tabs
 */
export function RelatedPeersTabs(props) {
  if (props.peers.length < 2) {
    return null; // Nothing to do here.
  }

  const peers = props.peers.map((p) => (
    <li key={p.id} 
        className={props.protocolId === p.id ? "active" : ""} >
      <Link to={makePeerLinkProps(props.routeserverId, p.id)}>
        {p.address}
      </Link>
    </li>
  ));
  
  return (
    <ul className="related-peers">
      {peers}
    </ul>
  );

}

/*
 * Render a card with related peers for the sidebar.
 *
 * This provides quick links to the same peer on other
 * routeservers.
 */
export function RelatedPeersCard(props) {
  let neighbors = props.neighbors;
  if (!neighbors || neighbors.length < 2) {
    return null; // nothing to render here.
  }

  // Exclude own neighbor
  let related = [];
  for neighbor in neighbors {
    if (neighbor.routeserver_id == props.rsId &&
        neighbor.id == props.protocolId) {
          continue; // Skip current peer.
    }
  }

  return (
    <div className="card card-related-peers">
      <h2 className="card-header">Related Neighbors</h2>
    </div>
  );
}


