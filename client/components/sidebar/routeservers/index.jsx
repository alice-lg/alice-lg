
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
    <li key={group}>
      <button className="btn btn-link btn-option" onClick={() => props.onChange(group)}>
        {group}
      </button>
    </li>
  ));

  return (
    <div className="routeservers-groups-select">
        <div className="dropdown">
            <button className="btn btn-default dropdown-toggle btn-select"
                    type="button"
                    id="select-routeservers-group"
                    data-toggle="dropdown"
                    aria-haspopup="true"
                    aria-expanded="true">
               {props.selected}
               <span className="caret"></span>
            </button>
            <ul className="dropdown-menu" aria-labelledby="select-routeservers-group">
              {options}
            </ul>
        </div>
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


