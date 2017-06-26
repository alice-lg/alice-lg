
/*
 * Alice (Prefix-)Lookup
 */

import React from 'react'
import {connect} from 'react-redux'

import {loadResults} from './actions'

import LookupResults from './results'
import SearchInput from 'components/search-input/debounced'

class Lookup extends React.Component {
  doLookup(q) {
    this.props.dispatch(loadResults(q));
  }

  render() {
    return (
      <div className="lookup-container">
        <div className="card">
          <SearchInput
            placeholder="Search for prefixes by entering a network address"
            onChange={(e) => this.doLookup(e.target.value)}  />
        </div>

        <LookupResults />
      </div>
    )
  }
}

export default connect(
  (state) => {
    return {
        isLoading: state.lookup.isLoading,
        error: state.lookup.error
    }
  }
)(Lookup);


