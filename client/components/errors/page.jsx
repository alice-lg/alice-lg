import React from 'react'
import {connect} from 'react-redux'

import {resetApiError} from './actions'

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

    if (status == 429) {
      body = (
        <div className="error-message">
          <p>Alice reached the request limit.</p>
          <p>We suggest you try at a less busy time.</p>
        </div>
      );
    } else {
      body = (
        <div className="error-message">
          <p>
            Alice has trouble connecting to the API
            {this.props.error.response &&
              " (got HTTP " + this.props.error.response.status + ")"}
            .
          </p>
          <p>If this problem persist, we suggest you try again later.</p>
        </div>
      );
    }

    return (
      <div className="error-notify">
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
      error: state.errors.error
  })
)(ErrorsPage);
