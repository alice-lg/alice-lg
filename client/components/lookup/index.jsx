
/*
 * Alice (Prefix-)Lookup
 */

import React from 'react'
import {connect} from 'react-redux'

import {loadResults} from './actions'

import LookupResults from './results'
import SearchInput from 'components/search-input/debounced'


class LookupHelp extends React.Component {
  render() {
    if(this.props.query != '') {
      return null;
    }

    return (
      <div className="lookup-help">
        <h3>Did you know?</h3>
        <p>You can search for</p>
        <ul>
          <li><b>Network Addresses</b>,</li>
          <li><b>Peers</b> by entering their name and</li>
          <li><b>ASNs</b> by prefixing them with 'AS'</li>
        </ul>
        <p>Just start typing!</p>
      </div>
    );
  }
}


class Lookup extends React.Component {
  doLookup(q) {
    this.props.dispatch(loadResults(q));
  }

  componentDidMount() {
    // this is yucky but the debounced
    // search input seems to kill the ref=
    let input = document.getElementById('lookup-search-input');
    input.focus();
  }

  render() {
    return (
      <div className="lookup-container">
        <div className="card">
          <SearchInput
            id="lookup-search-input"
            placeholder="Search for prefixes on all routeservers"
            onChange={(e) => this.doLookup(e.target.value)}  />
        </div>

        <LookupHelp query={this.props.query} />

        <LookupResults />
      </div>
    )
  }
}

export default connect(
  (state) => {
    return {
        query: state.lookup.query,
        isLoading: state.lookup.isLoading,
        error: state.lookup.error
    }
  }
)(Lookup);


