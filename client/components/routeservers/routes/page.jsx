
import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'

import Details    from '../details'
import Status     from '../status'
import PageHeader from 'components/page-header'

import ProtocolName
  from 'components/routeservers/protocols/name'

import RoutesView  from './view'

import SearchInput from 'components/search-input'

import BgpAttributesModal
  from './bgp-attributes-modal'

// Actions
import {setRoutesFilterValue}
  from '../actions'
import {loadRouteserverProtocol}
  from 'components/routeservers/actions'

class RoutesPage extends React.Component {

  setFilter(value) {
    this.props.dispatch(
      setRoutesFilterValue(value)
    );
  }

  componentDidMount() {
    // Assert neighbors for RS are loaded
    this.props.dispatch(
      loadRouteserverProtocol(parseInt(this.props.params.routeserverId))
    );
  }

  render() {
    return(
      <div className="routeservers-page">
        <PageHeader>
          <Link to={`/routeservers/${this.props.params.routeserverId}`}>
            <Details routeserverId={this.props.params.routeserverId} />
          </Link>
          <span className="spacer">&raquo;</span>
          <ProtocolName routeserverId={this.props.params.routeserverId}
                        protocolId={this.props.params.protocolId} />
        </PageHeader>

        <BgpAttributesModal />

        <div className="row details-main">
          <div className="col-md-8">

            <div className="card">
              <SearchInput
                value={this.props.routesFilterValue}
                placeholder="Filter by Network or BGP next-hop"
                onChange={(e) => this.setFilter(e.target.value)}  />
            </div>

            <RoutesView
                routes="received"
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            <RoutesView
                routes="filtered"
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            <RoutesView
                routes="not-exported"
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

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
      routesFilterValue: state.routeservers.routesFilterValue
    }
  }
)(RoutesPage);

