

import moment from 'moment'

import React from 'react'


export default class RelativeTime extends React.Component {

  render() {
    if (!this.props.value) {
      return null;
    }

    let time = moment.utc(this.props.value);
    return (
      <span>{time.fromNow(this.props.suffix)}</span>
    );
  }
}



