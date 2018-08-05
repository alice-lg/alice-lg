import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'

import {resetApiError} from './actions'
import {infoFromError} from './utils'

class ErrorsPage extends React.Component {

  resetApiError() {
    this.props.dispatch(resetApiError());
  }

  render() {
    if (!this.props.error) {
      return null;
    }

    let status = null;
    if (this.props.error.response) {
      status = this.props.error.response.status;
    }

    if (!status || (status != 429 && status < 500)) {
      return null;
    }

    let body = null;


    // Find affected routeserver
    const errorInfo = infoFromError(this.props.error);
    const rsId = errorInfo.routeserver_id; 
    let rs = null;
    if (rsId !== null) {
      rs = _.findWhere(this.props.routeservers, { id: rsId });
    }

    if (status == 429) {
      body = (
        <div className="error-message">
          <p>Alice reached the request limit.</p>
          <p>We suggest you try at a less busy time.</p>
        </div>
      );
    } else {
      let errorStatus = "";
      if (this.props.error.response) {
        errorStatus = " (got HTTP " + this.props.error.response.status + ")";
      }
      if (errorInfo) {
        errorStatus = ` (got ${errorInfo.tag})`;
      }

      body = (
        <div className="error-message">
          <p>
            Alice has trouble connecting to the API 
            {rs && 
              <span> of <b>{rs.name}</b></span>}
              {errorStatus}
            .
          </p>
          <p>If this problem persist, we suggest you try again later.</p>
        </div>
      );
    }

    return (
      <div className="error-notify">
        <div className="error-dismiss">
          <i className="fa fa-times-circle" aria-hidden="true"
             onClick={() => this.resetApiError()}></i>
        </div>
        <div className="error-icon">
          <i className="fa fa-times-circle" aria-hidden="true"></i>
        </div>
        {body}
      </div>
    );
  }
}

export default connect(
  (state) => ({
      error: state.errors.error,
      routeservers: state.routeservers.all,
  })
)(ErrorsPage);
