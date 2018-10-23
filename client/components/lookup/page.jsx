
import React from 'react'
import {connect} from 'react-redux'

import PageHeader from 'components/page-header'

import Lookup from 'components/lookup'
import LookupSummary from 'components/lookup/results-summary'
import FiltersEditor from 'components/filters/editor'

import Content from 'components/content'

import {makeLinkProps} from './state'

class _LookupView extends React.Component {
  render() {
    if (this.props.enabled == false) {
      return null;
    }

    return (
      <div className="lookup-container details-main">
       <div className="col-main col-lg-9 col-md-12">
         <Lookup />
       </div>
       <div className="col-aside-details col-lg-3 col-md-12">
         <LookupSummary />
         <FiltersEditor makeLinkProps={makeLinkProps} />
       </div>
      </div>
    );
  }
}

const LookupView = connect(
  (state) => {
    return {
      enabled: state.config.prefix_lookup_enabled
    }
  }
)(_LookupView);


export default class LookupPage extends React.Component {
  render() {
    return (
      <div className="welcome-page">
       <PageHeader></PageHeader>
       <p></p>
       <LookupView />
      </div>
    );
  }
}

