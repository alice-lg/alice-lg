/**
 * Show BGP attributes as a modal dialog
 *
 * @author Matthias Hannig <mha@ecix.net>
 */

import React from 'react'
import {connect} from 'react-redux'

import Modal, {Header, Body, Footer} from 'components/modals/modal'

import BgpCommunitiyLabel from 'components/bgp-communities/label'

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

    let communities = [];
    if (attrs.bgp.communities) {
      communities = attrs.bgp.communities.map((c) => c.join(':'));
    }

    let large_communities = [];
    if (attrs.bgp.large_communities) {
      large_communities = attrs.bgp.large_communities.map((c) => c.join(':'));
    }

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
            <tr>
              <th>Communities:</th>
              <td>{communities.join(' ')}</td>
            </tr>
            <tr>
              <th></th>
              <td>
                {communities.map((c) => <BgpCommunitiyLabel community={c} key={c} />)}
              </td>
            </tr>
            {large_communities.length > 0 &&
                <tr>
                  <th>Large Communities:</th>
                  <td>{large_communities.join(' ')}</td>
                </tr>
                <tr>
                  <th></th>
                  <td>
                    {large_communities.map((c) => <BgpCommunitiyLabel community={c} key={c} />)}
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

