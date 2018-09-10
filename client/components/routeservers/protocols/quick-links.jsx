
import React from 'react'
import {connect} from 'react-redux'


class QuickLinks extends React.Component {
  
  render() {
    if (this.props.isLoading) {
      return null; // nothing to go to.
    }
    
    return (
      <div className="quick-links neighbors-quick-links">
        <span>Go to:</span>
        <ul>
          <li className="established">
            <a href="#sessions-up">Established</a>
          </li>
          <li className="down">
            <a href="#sessions-down">Down</a>
          </li>
        </ul>
      </div>
    );
  }

}

export default connect(
  (state) => ({
    isLoading: state.routeservers.protocolsAreLoading,
  })
)(QuickLinks);

