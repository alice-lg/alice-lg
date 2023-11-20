
import { Link
       , useParams
       }
  from 'react-router-dom';

import { useLocalRelatedPeers }
  from 'app/context/neighbors';

/**
 * RelatedPeersTabs show locally related peers as tabs
 */
const LocalRelatedPeersTabs = () => {
  const peers = useLocalRelatedPeers();
  const { neighborId, routeServerId } = useParams();

  if (peers.length < 2) {
    return null; // Nothing to do here.
  }
  
  const peerUrl = (n) =>
    `/routeservers/${routeServerId}/neighbors/${n.id}/routes`;

  const relatedPeers = peers.map((p) => (
    <li key={p.id} 
        className={neighborId === p.id ? "active" : ""} >
      <Link to={peerUrl(p)}>
        {p.address}
      </Link>
    </li>
  ));
  
  return (
    <ul className="related-peers">
      {relatedPeers}
    </ul>
  );
}

export default LocalRelatedPeersTabs;
