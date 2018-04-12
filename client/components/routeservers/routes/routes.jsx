import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'


import {loadRouteserverRoutes, loadRouteserverRoutesFiltered} from '../actions'
import {showBgpAttributes} from './bgp-attributes-modal-actions'

import LoadingIndicator
	from 'components/loading-indicator/small'


import FilterReason
  from 'components/routeservers/large-communities/filter-reason'

import NoexportReason
  from 'components/routeservers/large-communities/noexport-reason'


function _filteredRoutes(routes, filter) {
  let filtered = [];
  if(filter == "") {
    return routes; // nothing to do here
  }

  filter = filter.toLowerCase();

  // Filter protocols
  filtered = _.filter(routes, (r) => {
    return (r.network.toLowerCase().indexOf(filter) != -1 ||
            r.gateway.toLowerCase().indexOf(filter) != -1 ||
            r.interface.toLowerCase().indexOf(filter) != -1);
  });

  return filtered;
}

class RoutesTable extends React.Component {
  showAttributesModal(route) {
    this.props.dispatch(
      showBgpAttributes(route)
    );
  }


  render() {
    let routes = this.props.routes;
    const routes_columns = this.props.routes_columns;

    routes = _filteredRoutes(routes, this.props.filter);
    if (!routes || !routes.length) {
      return null;
    }

    const _lookup = (r, path) => {
      const split = path.split(".").reduce((acc, elem) => acc[elem], r);

      if (Array.isArray(split)) { // this is for the ASN column, dirty hack.
        var baseUrl = "http://irrexplorer.nlnog.net/search/"
				let asnLinks = split.map((asn, i) => {
					return (<a key={`${asn}_${i}`} href={baseUrl + asn} target="_blank">{asn} </a>);
				});
				return asnLinks;
      }
      return split;
    }

    let routesView = routes.map((r,i) => {
      return (
        <tr key={`${r.network}_${i}`}>
          <td>
            <a href="#" onClick={() => this.showAttributesModal(r)}> {r.network}</a>
            {this.props.display_reasons == "filtered" && <FilterReason route={r} />}
            {this.props.display_reasons == "noexport" && <NoexportReason route={r} />}
          </td>
          {Object.keys(routes_columns).map(col => <td key={col}
           onClick={() => this.showAttributesModal(r)}> {_lookup(r, col)}</td>)}
        </tr>
      );
    });

    return (
      <div className="card">
        {this.props.header}
        <table className="table table-striped table-routes">
          <thead>
            <tr>
              <th>Network</th>
              {Object.values(routes_columns).map(col => <th key={col}>{col}</th>)}
            </tr>
          </thead>
          <tbody>
            {routesView}
          </tbody>
        </table>
      </div>
    );
  }
}


RoutesTable = connect(
  (state) => {
    return {
      filter:         state.routeservers.routesFilterValue,
      reject_reasons: state.routeservers.reject_reasons,
      routes_columns: state.config.routes_columns,
    }
  }
)(RoutesTable);


class RoutesTables extends React.Component {
  componentDidMount() {
    this.props.dispatch(
      loadRouteserverRoutes(this.props.routeserverId, this.props.protocolId)
    );
  }

  render() {
    if(this.props.isLoading) {
      return (
				<LoadingIndicator />
      );
    }

    const routes = this.props.routes[this.props.protocolId];
    const filtered = this.props.filtered[this.props.protocolId] || [];
    const noexport = this.props.noexport[this.props.protocolId] || [];

    if((!routes || routes.length == 0) &&
			 (!filtered || filtered.length == 0)) {
      return(
        <p className="help-block">
          No routes matched your filter.
        </p>
      );
    }


    const received = routes.filter(r => filtered.indexOf(r) < 0);

    const mkHeader = (color, action) => (
        <p style={{"color": color, "textTransform": "uppercase"}}>
          Routes {action}
        </p>
    );

    const filtdHeader = mkHeader("orange", "filtered");
    const recvdHeader = mkHeader("green",  "accepted");
    const noexHeader  = mkHeader("red",    "not exported");


    return (
      <div>
        <RoutesTable header={filtdHeader} routes={filtered} display_reasons="filtered"/>
        <RoutesTable header={recvdHeader} routes={received} display_reasons={false}/>
        <RoutesTable header={noexHeader}  routes={noexport} display_reasons="noexport"/>
      </div>
    );

  }
}


export default connect(
  (state) => {
    return {
      isLoading: state.routeservers.routesAreLoading,
      routes:    state.routeservers.routes,
      filtered:  state.routeservers.filtered,
      noexport:  state.routeservers.noexport,
    }
  }
)(RoutesTables);
