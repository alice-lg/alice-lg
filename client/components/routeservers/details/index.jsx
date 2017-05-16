
import React from 'react'
import {connect} from 'react-redux'

class Details extends React.Component {
  render() {
    let rsStatus = this.props.details[this.props.routeserverId];
    if (!rsStatus) {
      return null;
    }

    // Get routeserver name
    let rs = this.props.routeservers[parseInt(this.props.routeserverId)];
    if (!rs) {
      return null;
    }

    return (
       <span className="status-name">{rs.name}</span>
    );
  }
}

export default connect(
  (state) => {
    return {
      routeservers: state.routeservers.all,
      details: state.routeservers.details
    }
  }
)(Details);

