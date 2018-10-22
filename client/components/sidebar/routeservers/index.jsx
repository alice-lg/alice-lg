
/**
 * Routeservers List component
 */


import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'
import {Link} from 'react-router'

import {loadRouteservers,
        selectGroup}
  from 'components/routeservers/actions'

// Components
import Status from './status'

const GroupSelect = (props) => {
  if (props.groups.length < 2) {
    return null; // why bother?
  }

  const options = props.groups.map((group) => (
    <option key={group} value={group}>{group}</option>
  ));

  return (
    <div className="routeservers-groups-select">
      <select value={props.selected}
              onChange={(e) => props.onChange(e.target.value)}>
        {options}
      </select>
    </div>
  );
}


class RouteserversList extends React.Component {
  componentDidMount() {
    this.props.dispatch(
      loadRouteservers()
    );
  }

  onSelectGroup(group) {
    this.props.dispatch(selectGroup(group));
  }

  render() {
    const rsGroup = _.where(this.props.routeservers, {
        group: this.props.selectedGroup,
    });

    const routeservers = rsGroup.map((rs) =>
      <li key={rs.id}>
        <Link to={`/routeservers/${rs.id}`} className="routeserver-id">{rs.name}</Link>
        <Status routeserverId={rs.id} />
      </li>
    );

    return (
      <div className="routeservers-list">
        <h2>route servers</h2>
        <GroupSelect groups={this.props.groups}
                     selected={this.props.selectedGroup}
                     onChange={(group) => this.onSelectGroup(group)} />
        <ul>
          {routeservers}
        </ul>
      </div>
    );
  }
}


export default connect(
  (state) => {
    return {
      routeservers: state.routeservers.all,

      groups: state.routeservers.groups,
      isGrouped: state.routeservers.isGrouped,
      selectedGroup: state.routeservers.selectedGroup,
    };
  }
)(RouteserversList);


