
import React from 'react'
import {connect} from 'react-redux'

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
  // Check for BGP large community RS:1000:[1..3]
  // as defined in https://www.euro-ix.net/en/forixps/large-bgp-communities/
  const communities = props.route.bgp.large_communities;
  const ownAsn = props.asn;

  let rpkiState = 0; // Not set

  for (const c of communities) {
    if (c[0] == ownAsn && c[1] == 1000) {
      rpkiState = c[2];  
    }
  }

  switch(rpkiState) {
    case 1:
      return (
        <span className="route-prefix-flag rpki-route rpki-valid">
          <i className="fa fa-check-circle" /> 
          <div>RPKI Valid</div>
        </span>
      );

    case 2:
      return (
        <span className="route-prefix-flag rpki-route rpki-unknown">
          <i className="fa fa-question-circle" />
          <div>RPKI Unknown</div>
        </span>
      );

    case 3:
      return (
        <span className="route-prefix-flag rpki-route rpki-not-checked"></span>
      );
  }

  if (rpkiState >= 4) { // Invalid
    const cls = `route-prefix-flag rpki-route rpki-invalid rpki-invalid-${rpkiState}`
    return (
      <span className={cls}>
        <i className="fa fa-minus-circle" />
        <div>RPKI Invalid</div>
      </span>
    );
  }

  return null;
}

export const RpkiIndicator = connect(
  (state) => ({
    asn: state.config.asn,
  })
)(_RpkiIndicator);


