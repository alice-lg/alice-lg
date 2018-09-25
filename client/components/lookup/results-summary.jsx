
import React from 'react'
import {connect} from 'react-redux'
import RelativeTime from 'components/relativetime'

class ResultsBox extends React.Component {

  render() {
    if (this.props.query == '') {
      return null;
    }

    const queryDuration = this.props.queryDuration.toFixed(2);
    const cachedAt = this.props.cachedAt;
    const cacheTtl = this.props.cacheTtl;

    return (
      <div className="card">
        <div className="lookup-result-summary">
          <ul>
            <li>
              Found <b>{this.props.totalImported}</b> accepted 
              and <b>{this.props.totalFiltered}</b> filtered routes.
            </li>
            <li>Query took <b>{queryDuration} ms</b> to complete.</li>
            <li>Routes cache was built <b><RelativeTime value={cachedAt} /> </b>
                and will be refreshed <b><RelativeTime value={cacheTtl} /></b>.
            </li>
          </ul>
        </div>
      </div>
    );
  }
}


export default connect(
  (state) => {
    return {
      totalImported: state.lookup.totalRoutesImported,
      totalFiltered: state.lookup.totalRoutesFiltered, 

      cachedAt: state.lookup.cachedAt,
      cacheTtl: state.lookup.cacheTtl,

      queryDuration: state.lookup.queryDurationMs
    }
  }
)(ResultsBox)

