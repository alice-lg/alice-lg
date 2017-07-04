
import React from 'react'
import {connect} from 'react-redux'

import PageHeader from 'components/page-header'

import Lookup from 'components/lookup'
import LookupSummary from 'components/lookup/results-summary'

class LookupView extends React.Component {
  render() {
    if (this.props.enabled == false) {
      return null;
    }

    return (
      <div className="lookup-container">
       <div className="col-md-8">
         <Lookup />
       </div>
       <div className="col-md-4">
         <LookupSummary />
       </div>
      </div>
    );
  }
}

const LookupPage = connect(
  (state) => {
    return {
      enabled: state.config.prefix_lookup_enabled
    }
  }
)(LookupView);


export default class Welcome extends React.Component {
  render() {
    return (
      <div className="welcome-page">
       <PageHeader></PageHeader>

       <div className="jumbotron">
         <h1>Welcome to Alice!</h1>
         <p>Your friendly bird looking glass</p>
       </div>

       <LookupPage />

      </div>
    );
  }
}


