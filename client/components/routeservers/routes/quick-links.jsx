
import React from 'react'
import {connect} from 'react-redux'


/*
 * Quick links:
 * Jump to anchors for: not exported, filtered and received
 */

const QuickLinks = function(props) {

  const isLoading = props.routes.received.loading ||
                    props.routes.filtered.loading;
  
  // Do no display some dangleing "go to:" text
  if (isLoading) {
    return null;
  }

  // Handle special not exported: Default just works like
  // filtered or received. When loaded on demand, we override
  // this.
  let showNotExported = (!props.routes.notExported.loading &&
                          props.routes.notExported.totalResults > 0);
  if (props.loadNotExportedOnDemand) {
    // Show the link when nothing else is loading anymore
    showNotExported = !isLoading;
  }

  return (
    <div className="quick-links routes-quick-links">
      <span>Go to:</span>
      <ul>
        {(!props.routes.filtered.loading && 
           props.routes.filtered.totalResults > 0) &&
          <li className="filtered">
            <a href="#routes-filtered">Filtered</a></li>}
        {(!props.routes.received.loading &&
           props.routes.received.totalResults > 0) &&
          <li className="received">
            <a href="#routes-received">Accepted</a></li>}
        {showNotExported &&
          <li className="not-exported">
            <a href="#routes-not-exported">Not Exported</a></li>}
      </ul>
    </div>
  );
}

export default connect(
  (state) => ({
    "loadNotExportedOnDemand": state.config.noexport_load_on_demand, 
  })
)(QuickLinks);

