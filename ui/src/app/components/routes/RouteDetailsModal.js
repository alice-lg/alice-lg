
import { useCallback }
  from 'react';

import { useRouteDetails }
  from 'app/context/routes';

import { Modal
       , ModalHeader
       , ModalBody
       , ModalFooter
       }
  from 'app/components/modal/Modal';
import BgpCommunitiyLabel
  from 'app/components/routes/BgpCommunityLabel';


const RouteDetailsModal = () => {
  const [ route, setRoute ] = useRouteDetails();

  const onDismiss = useCallback(() => setRoute(null), [setRoute]);

  const attrs = route?.bgp;
  if (!attrs) {
    return null;
  }

  const communities = attrs.communities || [];
  const extCommunities = attrs.ext_communities || [];
  const largeCommunities = attrs.large_communities || [];

  // As communities can be repeated, we can not use them
  // directly as keys, but may have to prepend a suffix.
  const communityKeyCnt = {};
  const communityKey = (community) => {
    const k = community.join(":");
    communityKeyCnt[k] = (communityKeyCnt[k]||0) + 1;
    return `${k}+${communityKeyCnt[k]}`;
  };

  return (
    <Modal
      className="bgp-attributes-modal"
      onDismiss={onDismiss}>
        <ModalHeader onDismiss={onDismiss}>
          <p>BGP Attributes for Network:</p>
          <h4>{route.network}</h4>
        </ModalHeader>
        <ModalBody>
          <table className="table table-nolines">
           <tbody>
            <tr>
              <th>Origin:</th><td>{attrs.origin}</td>
            </tr>
            <tr>
              <th>Local Pref:</th><td>{attrs.local_pref}</td>
            </tr>
            <tr>
             <th>Next Hop:</th><td>{attrs.next_hop}</td>
            </tr>
            <tr>
                <th>MED:</th>
                <td>{attrs.med}</td>
            </tr>
            {attrs.as_path &&
                <tr>
                  <th>AS Path:</th><td>{attrs.as_path.join(' ')}</td>
                </tr>}
            {communities.length > 0 &&
                <tr>
                  <th>Communities:</th>
                  <td>
                    {communities.map((c) => <BgpCommunitiyLabel community={c} key={communityKey(c)} />)}
                  </td>
                </tr>}
            {extCommunities.length > 0 &&
              <tr>
                <th>Ext. Communities:</th>
                <td>
                    {extCommunities.map((c) => <BgpCommunitiyLabel community={c} key={communityKey(c)} />)}
                </td>
              </tr>}
            {largeCommunities.length > 0 &&
                <tr>
                  <th>Large Communities:</th>
                  <td>
                    {largeCommunities.map((c) => <BgpCommunitiyLabel community={c} key={communityKey(c)} />)}
                  </td>
                </tr>}
           </tbody>
          </table>
        </ModalBody>

        <ModalFooter>
          <button className="btn btn-default"
                  onClick={onDismiss}>Close</button>
        </ModalFooter>
    </Modal>
  );
}

export default RouteDetailsModal;
