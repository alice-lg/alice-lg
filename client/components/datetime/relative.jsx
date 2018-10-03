
import moment from 'moment'

import React from 'react'

export default class RelativeTime extends React.Component {

  render() {
    if (!this.props.value) {
      return null;
    }

    let time = false;
    if (this.props.value instanceof moment) {
      time = this.props.value;
    } else {
      time = moment.utc(this.props.value);
    }

    // A few seconds ago / in a few seconds can be replaced 
    // with 'just now'.
    // fuzzyNow can be set as a threshold of seconds
    if (this.props.fuzzyNow) {
      const now = moment.utc();
      if (Math.abs(now - time) / 1000.0 < this.props.fuzzyNow) {
        return (
          <span>just now</span>
        );
      }
    }

    return (
      <span>{time.fromNow(this.props.suffix)}</span>
    );
  }
}

