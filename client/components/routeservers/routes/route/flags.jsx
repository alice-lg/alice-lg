
import React from 'react'
import {connect} from 'react-redux'

import {isRejectCandidate}
  from 'components/routeservers/communities/utils'

/*
 * Primary Route Indicator
 */
export const PrimaryIndicator = function(props) {
  if (props.route.primary) {
    return(
      <span className="route-prefix-flag primary-route is-primary-route"><i className="fa fa-star"></i>
        <div>Best Route</div>
      </span>
    );
  }

  // Default
  return (
    <span className="route-prefix-flag primary-route not-primary-route"></span>
  );
}

/*
 * Blackhole Route Indicator
 */
export const BlackholeIndicator = function(props) {
  const blackholes = props.blackholes || [];
  const communities = props.route.bgp.communities;
  const nextHop = props.route.bgp.next_hop;

  // Check if next hop is a known blackhole
  let isBlackhole = blackholes.includes(nextHop);

  // Check if BGP community 65535:666 is set
  for (const c of communities) {
    if (c[0] == 65535 && c[1] == 666) {
      isBlackhole = true;
      break;
    }
  }

  if (isBlackhole) {
    return(
      <span className="route-prefix-flag blackhole-route is-blackhole-route"><i className="fa fa-circle"></i>
        <div>Blackhole</div>
      </span>
    );
  }

  return (
    <span className="route-prefix-flag blackhole-route not-blackhole-route"></span>
  );
}

/*
 * RPKI Status Indicators
 */
const _RpkiIndicator = function(props) {
  // Check if indicator is enabled
  if (props.rpki.enabled == false) { return null; }

  // Check for BGP large communities as configured in the alice.conf
  const rpkiValid      = props.rpki.valid;
  const rpkiUnknown    = props.rpki.unknown;
  const rpkiNotChecked = props.rpki.not_checked;
  const rpkiInvalid    = props.rpki.invalid;

  const communities = props.route.bgp.large_communities;
  for (const com of communities) {

    // RPKI VALID
    if (com[0].toFixed() === rpkiValid[0] &&
        com[1].toFixed() === rpkiValid[1] &&
        com[2].toFixed() === rpkiValid[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-valid">
          <i className="fa fa-check-circle" /> 
          <div>RPKI Valid</div>
        </span>
      );
    }

    // RPKI UNKNOWN
    if (com[0].toFixed() === rpkiUnknown[0] &&
        com[1].toFixed() === rpkiUnknown[1] &&
        com[2].toFixed() === rpkiUnknown[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-unknown">
          <i className="fa fa-question-circle" />
          <div>RPKI Unknown</div>
        </span>
      );
    }

    // RPKI NOT CHECKED
    if (com[0].toFixed() === rpkiNotChecked[0] &&
        com[1].toFixed() === rpkiNotChecked[1] &&
        com[2].toFixed() === rpkiNotChecked[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-not-checked">
          <i className="fa fa-circle-o" />
          <div>RPKI not checked</div>
        </span>
      );
    }

    // RPKI INVALID
    // Depending on the configration this can either be a
    // single flag or a range with a given reason.
    let rpkiInvalidReason = 0;
    if (com[0].toFixed() === rpkiInvalid[0] &&
        com[1].toFixed() === rpkiInvalid[1]) {

      // This needs to be considered invalid, now try to detect why
      if (rpkiInvalid.length > 3 && rpkiInvalid[3] == "*") {
        // Check if token falls within range
        const start = parseInt(rpkiInvalid[2], 10);
        if (com[2] >= start) {
          rpkiInvalidReason = com[2];
        }
      } else {
        if (com[2].toFixed() === rpkiInvalid[2]) {
          rpkiInvalidReason = 1;
        }
      }
    }

    // This in invalid, render it!
    if (rpkiInvalidReason > 0) {
      const cls = `route-prefix-flag rpki-route rpki-invalid rpki-invalid-${rpkiInvalidReason}`
      return (
        <span className={cls}>
          <i className="fa fa-minus-circle" />
          <div>RPKI Invalid</div>
        </span>
      );
    }
  }

  return null;
}

export const RpkiIndicator = connect(
  (state) => ({
    rpki: state.config.rpki,
    asn: state.config.asn,
  })
)(_RpkiIndicator);


/*
 * Reject Candidate Indicator
 */

class _RejectCandidateIndicator extends React.Component {

  render() {
    if (!this.props.candidateCommunities) {
      return null;
    }
    if (!isRejectCandidate(this.props.candidateCommunities, this.props.route)) {
      return null;
    }

    const cls = `route-prefix-flag reject-candidate-route`;
    return (
      <span className={cls}>
        <i className="fa fa-exclamation-circle" />
        <div>Reject Candidate</div>
      </span>
    );
  }

}

export const RejectCandidateIndicator = connect(
  (state) => ({
    candidateCommunities: state.routeservers.rejectCandidates.communities,
  })
)(_RejectCandidateIndicator);

