

import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'

import {loadRouteserverProtocol}
  from 'components/routeservers/actions'

import {Link} from 'react-router'

import RelativeTime
	from 'components/relativetime'

import LoadingIndicator
	from 'components/loading-indicator/small'

function _filteredProtocols(protocols, filter) {
  let filtered = [];
  if(filter == "") {
    return protocols; // nothing to do here
  }

  filter = filter.toLowerCase();

  // Filter protocols
  filtered = _.filter(protocols, (p) => {
    return (p.address.toLowerCase().indexOf(filter) != -1 ||
            p.description.toLowerCase().indexOf(filter) != -1);
  });

  return filtered;
}


class RoutesLink extends React.Component {
  render() {
    let url = `/routeservers/${this.props.routeserverId}/protocols/${this.props.protocol}/routes`;
    if (this.props.state != 'up') {
      return (<span>{this.props.children}</span>);
    }
    return (
      <Link to={url}>
        {this.props.children}
      </Link>
    )
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
        {neighbour.state != "up" && neighbour.last_error &&
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
      <RelativeTime value={props.neighbour.details.state_changed}
                    suffix={true} />
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

    let header = columnsOrder.map((col) => {
      return (
        <th key={col}>{columns[col]}</th>
      );
    });

    let neighbours = this.props.neighbours.map( (n) => {
      let neighbourColumns = columnsOrder.map((col) => {
        return <NeighbourColumn key={col}
                                rsId={this.props.routeserverId}
                                column={col}
                                neighbour={n} />
      });
      return (
        <tr key={n.id}>
          {neighbourColumns}
        </tr>
      );
      /*
      return (
        <tr key={n.id}>
          <td>
            <RoutesLink routeserverId={this.props.routeserverId}
                        protocol={n.id}
                        state={n.state}>
              {n.routes_received}
            </RoutesLink>
          </td>
        <td>
            <RoutesLink routeserverId={this.props.routeserverId}
                        protocol={n.id}
                        state={n.state}>
              {n.routes_filtered}
            </RoutesLink>
          </td>
        </tr>
      );
      */
    });

    let uptimeTitle;
    switch(this.props.state) {
      case 'up':
        uptimeTitle = 'Uptime'; break;
      case 'down':
        uptimeTitle = 'Downtime'; break;
      case 'start':
        uptimeTitle = 'Since'; break;
    }

    return (
      <div className="card">
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

    protocol = _filteredProtocols(protocol, this.props.filter);
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
      switch(n.state) {
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
    if (neighboursDown.length) {
      tables.push(<NeighboursTable key="down" state="down"
                                   neighbours={neighboursDown}
                                   routeserverId={this.props.routeserverId} />);
    }
    if (neighboursIdle.length) {
      tables.push(<NeighboursTable key="start" state="start"
                                   neighbours={neighboursIdle}
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
      filter: state.routeservers.protocolsFilterValue
    }
  }
)(Protocols);

