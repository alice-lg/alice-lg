
import React from 'react'
import {connect} from 'react-redux'

import PageHeader from 'components/page-header'
import Details from './details'
import Status from './status'

import SearchInput from 'components/search-input'

import Protocols from './protocols'

import {setProtocolsFilterValue} from './actions'

class RouteserversPage extends React.Component {

  setFilter(value) {
    this.props.dispatch(
      setProtocolsFilterValue(value)
    );
  }

  render() {
    return(
      <div className="routeservers-page">
        <PageHeader>
          <Details routeserverId={this.props.params.routeserverId} />
        </PageHeader>

        <div className="row details-main">
          <div className="col-md-8">
            <div className="card">
              <SearchInput
                value={this.props.protocolsFilterValue}
                placeholder="Filter by Neighbour, ASN or Description"
                onChange={(e) => this.setFilter(e.target.value)}
              />
            </div>

            <Protocols protocol="bgp" routeserverId={this.props.params.routeserverId} />
          </div>
          <div className="col-md-4">
            <div className="card">
              <Status routeserverId={this.props.params.routeserverId} />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default connect(
  (state) => {
    return {
      protocolsFilterValue: state.routeservers.protocolsFilterValue
    };
  }
)(RouteserversPage);


