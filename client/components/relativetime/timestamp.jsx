
import React from 'react'
import moment from 'moment'


window.momnet = moment;

export default class RelativeTimestamp extends React.Component {
  render() {

    let now = moment.utc()
    let rel = moment(now._d.getTime() - (this.props.value / 1000.0 / 1000.0))

    return (
      <span>{rel.fromNow(this.props.suffix)}</span>
    );
  }
}


