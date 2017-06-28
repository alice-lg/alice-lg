
import React from 'react'
import {connect} from 'react-redux'


class ResultsBox extends React.Component {

  render() {
    if (this.props.query == '') {
      return null;
    }

    const queryDuration = this.props.queryDuration.toFixed(2);

    return (
      <div className="card">
        <div className="lookup-result-summary">
          <ul>
            <li>
            Displaying <b>{this.props.resultsCount}</b> of <b>{this.props.total}</b> routes
            </li>
            <li>Query took <b>{queryDuration} ms</b> to complete</li>
          </ul>
        </div>
      </div>
    );
  }
}


export default connect(
  (state) => {
    return {
      query: state.lookup.query,
      resultsCount: state.lookup.results.length,
      start: state.lookup.offset,
      end: state.lookup.limit + state.lookup.offset,
      total: state.lookup.totalRoutes,
      queryDuration: state.lookup.queryDurationMs
    }
  }
)(ResultsBox)

