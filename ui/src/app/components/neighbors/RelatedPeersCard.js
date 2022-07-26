import { useParams }
  from 'react-router-dom';

import { isUpState }
  from 'app/components/neighbors/state';

import { useRelatedNeighbors }
  from 'app/context/neighbors';
import { useRouteServersMap }
  from 'app/context/route-servers';

import RelativeTimestamp
  from 'app/components/datetime/RelativeTimestamp';

const RoutesStats = ({peer}) => {
  if (!isUpState(peer.state)) {
    return null; // Nothing to render 
  }
  return (
    <div className="related-peers-routes-stats">
      <span className="atooltip routes-received">
        {peer.routes_received}
        <i>Routes Received</i>
      </span> / <span className="atooltip routes-accepted">
        {peer.routes_accepted}
        <i>Routes Accepted</i>
      </span> / <span className="atooltip routes-filtered">
        {peer.routes_filtered}
        <i>Routes Filtered</i>
      </span> / <span className="atooltip routes-exported">
        {peer.routes_exported}
        <i>Routes Exported</i>
      </span>
    </div>
  );
}

/*
 * Display a link to a peer. If the peer state is up.
 */
const PeerLink = ({to, children}) => {
  const neighbor = to;
  if (!neighbor) {
    return null;
  }

  const pid = neighbor.id;
  const rid = neighbor.routeserver_id;

  let peerUrl;
  if (isUpState(neighbor.state)) {
    peerUrl = `/routeservers/${rid}/protocols/${pid}/routes`;
  } else {
    peerUrl = `/routeservers/${rid}#sessions-down`;
  }
  // Render link
  return (
    <a href={peerUrl}>{children}</a>
  );
}

const normalizePeerState = (state) =>
  isUpState(state) ? 'up' : state;


/*
 * Render a card with related peers for the sidebar.
 *
 * This provides quick links to the same peer on other
 * routeservers.
 */
const RelatedPeersCard = () => {
  const { routeServerId, neighborId } = useParams();
  const { neighbors } = useRelatedNeighbors();
  const routeServersMap = useRouteServersMap();

  if (neighbors.length < 2) {
    return null; // nothing to render here.
  }

  // Exclude own neighbor and group peers by routeserver
  let related = {};
  for (let neighbor of neighbors) {
    if (neighbor.routeserver_id === routeServerId &&
        neighbor.id === neighborId) {
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
    relatedRs.push(routeServersMap[rsId]); 
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
                  <td className="peer-address">
                    <PeerLink to={peer}>{peer.address}</PeerLink>
                  </td>
                  <td className="peer-stats">
                    <RoutesStats peer={peer} />
                  </td>
                  <td className="uptime">
                    {normalizePeerState(peer.state)} {peer.uptime > 0 && <>
                        for <RelativeTimestamp value={peer.uptime} suffix={true} />
                      </>
                    }
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ))}
    </div>
  );
};

export default RelatedPeersCard;
