
import React from 'react'
import moment from 'moment'

export default class RelativeTimestamp extends React.Component {
  render() {
    const tsMs = this.props.value / 1000.0 / 1000.0; // nano -> micro -> milli
    const now = moment.utc()
    const rel = now.subtract(tsMs, 'ms');

    return (
      <span>{rel.fromNow(this.props.suffix)}</span>
    );
  }
}

