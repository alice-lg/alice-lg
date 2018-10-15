
import {debounce} from 'underscore'

import React from 'react'
import {connect} from 'react-redux'
import {replace} from 'react-router-redux'

import PageHeader from 'components/page-header'
import Details from './details'
import Status from './status'

import SearchInput from 'components/search-input'

import Protocols from './protocols'
import QuickLinks from './protocols/quick-links'

import {setFilterValue} from './protocols/actions'
import {makeQueryLinkProps} from './protocols/routing'


class RouteserversPage extends React.Component {

  constructor(props) {
    super(props);
    this.dispatchDebounced = debounce(this.props.dispatch, 350);
  }


  setFilter(value) {
    // Set filter value (for input rendering)
    this.props.dispatch(setFilterValue(value));

    // Update location delayed
    this.dispatchDebounced(replace(
      makeQueryLinkProps(
        this.props.routing,
        value,
        this.props.sortColumn,
        this.props.sortOrder)));
  }

  render() {
    return(
      <div className="routeservers-page">
        <PageHeader>
          <Details routeserverId={this.props.params.routeserverId} />
        </PageHeader>

        <div className="row details-main">
          <div className="col-main col-lg-9 col-md-12">
            <div className="card">
              <SearchInput
                value={this.props.filterValue}
                placeholder="Filter by Neighbour, ASN or Description"
                onChange={(e) => this.setFilter(e.target.value)}
              />
            </div>
            <QuickLinks />

            <Protocols protocol="bgp" routeserverId={this.props.params.routeserverId} />
          </div>
          <div className="col-lg-3 col-md-12 col-aside-details">
            <div className="card">
              <Status routeserverId={this.props.params.routeserverId}
                      cacheStatus={this.props.cacheStatus} />
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
      routing: state.routing.locationBeforeTransitions,

      filterValue: state.neighbors.filterValue,
      sortColumn:  state.neighbors.sortColumn,
      sortOrder:   state.neighbors.sortOrder,

      cacheStatus: {
        generatedAt: state.neighbors.cachedAt,
        ttl: state.neighbors.cacheTtl,
      }

    };
  }
)(RouteserversPage);


