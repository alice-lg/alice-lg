
import React from 'react'
import {connect} from 'react-redux'
import {replace} from 'react-router-redux'

import PageHeader from 'components/page-header'

import Lookup from 'components/lookup'
import LookupSummary from 'components/lookup/results-summary'

import Content from 'components/content'

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
      </div>
    );
  }
}

const LookupWidget = connect(
  (state) => {
    return {
      enabled: state.config.prefix_lookup_enabled
    }
  }
)(LookupView);


class Welcome extends React.Component {
  componentDidMount() {
    // Check if there is a query already set
    if (this.props.query != "") {
      // We should redirect to the search page
      const destination = {
        pathname: "/search",
        search: `?q=${this.props.query}`
      };
      this.props.dispatch(replace(destination));
    }
  }

  render() {
    return (
      <div className="welcome-page">
       <PageHeader></PageHeader>

       <div className="jumbotron">
         <h1><Content id="welcome.title">Welcome to Alice!</Content></h1>
         <p><Content id="welcome.tagline">Your friendly bird looking glass</Content></p>
       </div>

       <LookupWidget />

      </div>
    );
  }
}

export default connect(
  (state) => ({
    query: state.lookup.query,
  })
)(Welcome);

