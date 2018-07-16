
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


// Constants
import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from './actions';


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
                value={this.props.filterQuery}
                placeholder="Filter by Network or BGP next-hop"
                onChange={(e) => this.setFilter(e.target.value)}  />
            </div>

            <RoutesView
                type={ROUTES_RECEIVED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            <RoutesView
                type={ROUTES_FILTERED}
                routeserverId={this.props.params.routeserverId}
                protocolId={this.props.params.protocolId} />

            <RoutesView
                type={ROUTES_NOT_EXPORTED}
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
    let received = {
      routes:       state.routes.received,
      loading:      state.routes.receivedLoading,
      page:         state.routes.receivedPage,
      totalPages:   state.routes.receivedTotalPages,
      totalResults: state.routes.receivedTotalResults,
    };
    let filtered = {
      routes:       state.routes.filtered,
      loading:      state.routes.filteredLoading,
      page:         state.routes.filteredPage,
      totalPages:   state.routes.filteredTotalPages,
      totalResults: state.routes.filteredTotalResults,
    };
    let notExported = {
      routes:       state.routes.notExported,
      loading:      state.routes.notExportedLoading,
      page:         state.routes.notExportedPage,
      totalPages:   state.routes.notExportedTotalPages,
      totalResults: state.routes.notExportedTotalResults,
    };
    return({
      filterQuery: state.routes.filterQuery,
      routes: {
          received:    received,
          filtered:    filtered,
          notExported: notExported
      }
    });
  }
)(RoutesPage);

