

import _ from 'underscore'
import bigInt from 'big-integer';

import React from 'react'
import {connect} from 'react-redux'
import {Link} from 'react-router'
import {push} from 'react-router-redux'

import {loadRouteserverProtocol}
  from 'components/routeservers/actions'

import RelativeTimestamp
	from 'components/datetime/relative-timestamp'

import LoadingIndicator
	from 'components/loading-indicator/small'

import {ipToNumeric} from 'components/utils/ip'
import {urlEscape} from 'components/utils/query'


function _filteredProtocols(protocols, filter) {
  let filtered = [];
  if(filter == "") {
    return protocols; // nothing to do here
  }

  // We support different filter modes:
  // - Default: Try to match as much as possible
  // - AS$num: Try to match ASN only
  const filterAsn = _getFilterAsn(filter);
  if (filterAsn) {
    filtered = _.filter(protocols, (p) => {
      return (p.asn == filterAsn);
    });
  } else {
    filter = filter.toLowerCase();
    filtered = _.filter(protocols, (p) => {
      return (p.asn == filter ||
              p.address.toLowerCase().indexOf(filter) != -1 ||
              p.description.toLowerCase().indexOf(filter) != -1);
    });
  }

  return filtered;
}


function _getFilterAsn(filter) {
  const tokens = filter.split("AS", 2);
  if (tokens.length !== 2) {
    return false; // Not an ASN query
  }
  const asn = parseInt(tokens[1], 10);

  // Check if ASN is a valid number
  if (asn >= 0 == false) {
    return false;
  }

  return asn;
}

function _sortAnum(sort) {
  return (a, b) => {
    const va = a[sort];
    const vb = b[sort];
    if (va < vb ) { return -1; }
    if (va > vb ) { return 1;  }
    return 0;
  }
}

function _sortIpAddr(sort) {
  return (a, b) => {
    const va = ipToNumeric(a[sort]);
    const vb = ipToNumeric(b[sort]);

    // Handle ipv6 case
    if (va instanceof bigInt) {
      return va.compareTo(vb);
    }

    if (va < vb ) { return -1; }
    if (va > vb ) { return 1;  }
    return 0;
  }
}

function _sortOrder(cmp, order) {
  return (a, b) => {
    const res = cmp(a, b);
    if (order == 'desc') {
      return res * -1;
    }
    return res;
  }
}

function _sortNeighbors(neighbors, sort, order) {
  // Make compare function
  let cmp = _sortAnum(sort);
  if (sort == "address") {
    cmp = _sortIpAddr(sort);
  }
  return neighbors.sort(_sortOrder(cmp, order));
}


class RoutesLink extends React.Component {
  render() {
    let url = `/routeservers/${this.props.routeserverId}/protocols/${this.props.protocol}/routes`;
    if (this.props.state.toLowerCase() != 'up') {
      return (<span>{this.props.children}</span>);
    }
    return (
      <Link to={url}>
        {this.props.children}
      </Link>
    )
  }
}

class NeighborColumnHeader extends React.Component {
  render() {
    const baseUrl = `/routeservers/${this.props.rsId}`;
    const name = this.props.columns[this.props.column];
    const sortColumn = this.props.column.toLowerCase();
    const active = sortColumn == this.props.sort;
    const query = urlEscape(this.props.query);
    let cls = `col-neighbor-attr col-neighbor-${this.props.column} `;

    // Render link with sorting indicator
    if (active) {
      const nextOrder = (this.props.order == 'asc') ? 'desc' : 'asc';
      const url = `${baseUrl}?s=${sortColumn}&o=${nextOrder}&q=${query}`;
      let indicator = <i className="fa fa-arrow-circle-up"></i>;

      cls += 'col-neighbor-active ';

      if (this.props.order == 'desc') {
        indicator = <i className="fa fa-arrow-circle-down"></i>;
      }
      return (
        <th className={cls}>
          <Link to={url}>{name} {indicator}</Link>
        </th>
      );
    }
    
    // Column is not active, just present a link:
    const url = `${baseUrl}?s=${sortColumn}&o=${this.props.order}&q=${query}`
    return(
      <th className={cls}>
        <Link to={url}>{name}</Link>
      </th>
    );
  }
}



//
// Neighbours Columns Components
//
// * Render columums either with direct property
//   access, or
// * Use a "widget", a rendering function to which
//   the neighbour is passed.

// Helper:
const lookupProperty = function(obj, path) {
  let property = path.split(".").reduce((acc, part) => acc[part], obj);
  if (typeof(property) == "undefined") {
    property = `Property "${path}" not found in object.`;
  }

  return property;
}

// Widgets:
const ColDescription = function(props) {
  const neighbour = props.neighbour;
  return (
    <td>
      <RoutesLink routeserverId={props.rsId}
                  protocol={neighbour.id}
                  state={neighbour.state}>
        {neighbour.description}
        {neighbour.state.toLowerCase() != "up" && 
         neighbour.last_error &&
          <span className="protocol-state-error">
              {neighbour.last_error}
          </span>}
      </RoutesLink>
    </td>
  );
}

const ColUptime = function(props) {
  return (
    <td className="date-since">
      <RelativeTimestamp value={props.neighbour.uptime} suffix={true} />
    </td>
  );
}


const ColLinked = function(props) {
  // Access neighbour property by path
  const property = lookupProperty(props.neighbour, props.column);
  return (
    <td>
      <RoutesLink routeserverId={props.rsId}
                  protocol={props.neighbour.id}
                  state={props.neighbour.state}>
        {property}
      </RoutesLink>
    </td>
  );
}

const ColPlain = function(props) {
  // Access neighbour property by path
  const property = lookupProperty(props.neighbour, props.column);
  return (
    <td>{property}</td>
  );
}

// Column:
const NeighbourColumn = function(props) {
  const neighbour = props.neighbour;
  const column = props.column;

  const widgets = {
    // Special cases
    "asn": ColPlain,
    "state": ColPlain,

    "Uptime": ColUptime,
    "Description": ColDescription,
  };

  // Get render function
  let Widget = widgets[column] || ColLinked;
  return (
    <Widget neighbour={neighbour} column={column} rsId={props.rsId} />
  );
}



class NeighboursTableView extends React.Component {

  render() {
    const columns = this.props.neighboursColumns;
    const columnsOrder = this.props.neighboursColumnsOrder;

    const sortedNeighbors = _sortNeighbors(this.props.neighbours,
                                           this.props.sortColumn,
                                           this.props.sortOrder);

    let header = columnsOrder.map((col) => {
      return (
        <NeighborColumnHeader key={col}
                              rsId={this.props.routeserverId}
                              columns={columns} column={col}
                              sort={this.props.sortColumn}
                              order={this.props.sortOrder}
                              query={this.props.filterQuery} />
      );
    });

    let neighbours = sortedNeighbors.map((n) => {
      let neighbourColumns = columnsOrder.map((col) => {
        return <NeighbourColumn key={col}
                                rsId={this.props.routeserverId}
                                column={col}
                                neighbour={n} />
      });
      return <tr key={n.id}>{neighbourColumns}</tr>;
    });

    let uptimeTitle;
    let sectionTitle = '';
    let sectionAnchor = 'sessions-unknown';
    let sectionCls = 'card-header card-header-neighbors ';
    switch(this.props.state) {
      case 'up':
        uptimeTitle = 'Uptime'; 
        sectionAnchor = 'sessions-up';
        sectionTitle  = 'BGP Sessions Established';
        sectionCls  += 'established ';
        break;
      case 'down':
        uptimeTitle = 'Downtime'; 
        break;
      case 'start':
        sectionAnchor = 'sessions-up';
        uptimeTitle = 'Since';
        sectionAnchor = 'sessions-down';
        sectionTitle = 'BGP Sessions Down';
        sectionCls += 'down ';
    }


    return (
      <div className="card">
        <a name={sectionAnchor} />
        <p className={sectionCls}>{sectionTitle}</p>

        <table className="table table-striped table-protocols">
          <thead>
            <tr>
              {header}
            </tr>
          </thead>
          <tbody>
            {neighbours}
          </tbody>
        </table>
      </div>
    );
  }
}

const NeighboursTable = connect(
  (state) => ({
    neighboursColumns:      state.config.neighbours_columns,
    neighboursColumnsOrder: state.config.neighbours_columns_order,

    sortColumn: state.neighbors.sortColumn,
    sortOrder:  state.neighbors.sortOrder,
    filterQuery: state.neighbors.filterQuery,
  })
)(NeighboursTableView);


class Protocols extends React.Component {
  componentDidMount() {
    this.props.dispatch(
      loadRouteserverProtocol(parseInt(this.props.routeserverId))
    );
  }

  componentWillReceiveProps(nextProps) {
    if(this.props.routeserverId != nextProps.routeserverId) {
      this.props.dispatch(
        loadRouteserverProtocol(parseInt(nextProps.routeserverId))
      );
    }
  }

  render() {

    if(this.props.isLoading) {
      return (
        <div className="card">
					<LoadingIndicator />
        </div>
      );
    }


    let protocol = this.props.protocols[parseInt(this.props.routeserverId)];
    if(!protocol) {
      return null;
    }

    protocol = _filteredProtocols(protocol, this.props.filterQuery);
    if(!protocol || protocol.length == 0) {
      return (
        <div className="card">
          <p className="help-block">
            No neighbours could be found.
          </p>
        </div>
      );
    }

    // Filter neighbours
    let neighboursUp = [];
    let neighboursDown = [];
    let neighboursIdle = [];

    for (let id in protocol) {
      let n = protocol[id];
      switch(n.state.toLowerCase()) {
        case 'up':
          neighboursUp.push(n);
          break;
        case 'down':
          neighboursDown.push(n);
          break;
        case 'start':
          neighboursIdle.push(n);
          break;
        default:
          neighboursUp.push(n);
          console.error("Couldn't classify neighbour by state:", n);
      }
    }


    // Render tables
    let tables = [];
    if (neighboursUp.length) {
      tables.push(<NeighboursTable key="up" state="up"
                                   neighbours={neighboursUp}
                                   routeserverId={this.props.routeserverId} />);
    }
    if (neighboursIdle.length) {
      tables.push(<NeighboursTable key="start" state="start"
                                   neighbours={neighboursIdle}
                                   routeserverId={this.props.routeserverId} />);
    }
    if (neighboursDown.length) {
      tables.push(<NeighboursTable key="down" state="down"
                                   neighbours={neighboursDown}
                                   routeserverId={this.props.routeserverId} />);
    }

    return (
      <div>{tables}</div>
    );
  }
}


export default connect(
  (state) => {
    return {
      isLoading: state.routeservers.protocolsAreLoading,
      protocols: state.routeservers.protocols,

      filterQuery: state.neighbors.filterQuery,
      routing: state.routing.locationBeforeTransitions,
    }
  }
)(Protocols);

