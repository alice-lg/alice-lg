
import moment from 'moment'

import React from 'react'

export default class RelativeTime extends React.Component {

  // Local state updates, to trigger a rerender
  // every second for time updates.
  componentDidMount() {
    this.timer = setInterval(() => {
      this.setState({
        now: Date.now()
      })
    }, 1000);
  }

  // Stop timer
  componentWillUnmount() {
    clearInterval(this.timer);
  }

  // Helper: Assert time is an instance of moment
  getTime() {
    if (!this.props.value) {
      return false;
    }

    let time = false;
    if (this.props.value instanceof moment) {
      time = this.props.value;
    } else {
      time = moment.utc(this.props.value);
    }
    return time 
  }


  // Time can be capped, if we are handling a past
  // or future event:
  capTime(time) {
    const now = moment.utc();
    if (this.props.pastEvent && time.isAfter(now)) {
      return now;
    }

    if (this.props.futureEvent && time.isBefore(now)) {
      return now;
    }

    return time;
  }
  

  render() {
    let time = this.getTime();
    if (!time) {
      return null; // Well, nothing to do here
    }

    time = this.capTime(time);

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

