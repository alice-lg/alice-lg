
/*
 * Primar Route Indicator
 */

import React from 'react'


export default PrimaryIndicator = function(props) {
  if (props.route.details && props.route.primary) {
    return(
      <span className="primary-route is-primary-route">&gt;
        <div>Best Route</div>
      </span>
    );
  }

  // Default
  return (
    <span className="primary-route not-primary-route"></span>
  )
}


