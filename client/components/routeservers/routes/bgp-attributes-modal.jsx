/**
 * Show BGP attributes as a modal dialog
 *
 * @author Matthias Hannig <mha@ecix.net>
 */

import React from 'react'
import {connect} from 'react-redux'

import Modal, {Header, Body, Footer} from 'components/modals/modal'

import BgpCommunitiyLabel
  from 'components/routeservers/communities/label'

import {hideBgpAttributesModal}
  from './bgp-attributes-modal-actions'



class BgpAttributesModal extends React.Component {
  closeModal() {
    this.props.dispatch(
      hideBgpAttributesModal()
    );
  }

  render() {
    let attrs = this.props.bgpAttributes;
    if (!attrs.bgp) {
      return null;
    }

    const communities = attrs.bgp.communities;
    const extCommunities = attrs.bgp.ext_communities;
    const largeCommunities = attrs.bgp.large_communities;

    // As communities can be repeated, we can not use them
    // directly as keys, but may have to prepend a suffix.
    const communityKeyCnt = {};
    const communityKey = (community) => {
      const k = community.join(":");
      communityKeyCnt[k] = (communityKeyCnt[k]||0) + 1;
      return `${k}+${communityKeyCnt[k]}`;
    };

    return (
      <Modal className="bgp-attributes-modal"
             show={this.props.show}
             onClickBackdrop={() => this.closeModal()}>

        <Header onClickClose={() => this.closeModal()}>
          <p>BGP Attributes for Network:</p>
          <h4>{attrs.network}</h4>
        </Header>

        <Body>
          <table className="table table-nolines">
           <tbody>
            <tr>
              <th>Origin:</th><td>{attrs.bgp.origin}</td>
            </tr>
            <tr>
              <th>Local Pref:</th><td>{attrs.bgp.local_pref}</td>
            </tr>
            <tr>
             <th>Next Hop:</th><td>{attrs.bgp.next_hop}</td>
            </tr>
            <tr>
                <th>MED</th>
                <td>{attrs.bgp.med}</td>
            </tr>
            {attrs.bgp && attrs.bgp.as_path &&
                <tr>
                  <th>AS Path:</th><td>{attrs.bgp.as_path.join(' ')}</td>
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
        </Body>

        <Footer>
          <button className="btn btn-default"
                  onClick={() => this.closeModal()}>Close</button>
        </Footer>

      </Modal>
    );
  }
}

export default connect(
  (state) => {
    return {
      show: state.modals.bgpAttributes.show,
      bgpAttributes: state.modals.bgpAttributes.bgpAttributes
    }
  }
)(BgpAttributesModal);

