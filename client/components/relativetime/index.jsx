

import moment from 'moment'

import React from 'react'


export default class RelativeTime extends React.Component {

  render() {
    let time = moment.utc(this.props.value);

    return (
      <span>{time.fromNow(this.props.suffix)}</span>
    )
  }
}



