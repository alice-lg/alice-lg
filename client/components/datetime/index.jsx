
/**
 * Datetime Component
 *
 * @author Matthias Hannig <mha@ecix.net>
 */


import React from 'react'

import moment from 'moment'


export default class Datetime extends React.Component {
  render() {
    let timefmt = this.props.format;
    if (!timefmt) {
      timefmt = 'LLLL';
    }

    let time = moment(this.props.value);
    return (
      <span>{time.format(timefmt)}</span>
    );
  }
}

