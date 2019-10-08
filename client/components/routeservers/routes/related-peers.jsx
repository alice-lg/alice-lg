
import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'

import {makePeerLinkProps} from './urls'

import RelativeTimestamp
	from 'components/datetime/relative-timestamp'

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
 * Display a link to a peer. If the peer state is up.
 */
function PeerLink(props) {
  const neighbor = props.to;
  if (!neighbor) {
    return null;
  }

  const pid = neighbor.id;
  const rid = neighbor.routeserver_id;
  let peerUrl = `/routeservers/${rid}/protocols/${pid}/routes`;

  if (neighbor.state == "up") {
    // Render link
    return (
      <a href={peerUrl}>{props.children}</a>
    );
  } else {
    // Only display the content
    return (<span>{props.children}</span>);
  }
}



/*
 * Render a card with related peers for the sidebar.
 *
 * This provides quick links to the same peer on other
 * routeservers.
 */
function RelatedPeersCardView(props) {
  let neighbors = props.neighbors;
  if (!neighbors || neighbors.length < 2) {
    return null; // nothing to render here.
  }

  // Exclude own neighbor and group peers by routeserver
  let related = {};
  for (let neighbor of neighbors) {
    if (neighbor.routeserver_id == props.rsId &&
        neighbor.id == props.protocolId) {
          continue; // Skip current peer.
    }

    if (!related[neighbor.routeserver_id]) {
      related[neighbor.routeserver_id] = [];
    }
    related[neighbor.routeserver_id].push(neighbor);
  }

  // Get routeserver info for routeserver id as key in object.
  let relatedRs = [];
  for (let rsId in related) {
    relatedRs.push(props.routeservers[rsId]); 
  }


  return (
    <div className="card card-related-peers">
      <h2 className="card-header">Related Neighbors</h2>
      {relatedRs.map(rs => (
        <div key={rs.id} className="related-peers-rs-group">
          <h3>{rs.name}</h3>
          <table className="related-peers-rs-peer">
            <tbody>
              {related[rs.id].map(peer => (
                <tr key={peer.id}>
                  <td>
                    <PeerLink to={peer}>{peer.address}</PeerLink>
                  </td>
                  <td>
                    {peer.description}
                  </td>
                  <td>
                    {peer.state}
                  </td>
                  <td>
                    <RelativeTimestamp 
                      value={peer.uptime}
                      suffix={true} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ))}
    </div>
  );
}

export let RelatedPeersCard = connect(
  (state) => ({
    routeservers: state.routeservers.byId
  })
)(RelatedPeersCardView);
