
import { useRouteservers }
  from 'app/components/routeservers/Provider';

/**
 * GroupSelect shows a drop down for selecting a
 * group of routeservers.
 */
const GroupSelect = ({groups, selected, onChange}) => {
  if (groups.length < 2) {
    return null; // why bother?
  }

  const options = groups.map((group) => (
    <li key={group}>
      <button className="btn btn-link btn-option"
              onClick={() => onChange(group)}>
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
           {selected}
           <span className="caret"></span>
        </button>
        <ul className="dropdown-menu"
            aria-labelledby="select-routeservers-group">
          {options}
        </ul>
      </div>
    </div>
  );
}

/**
 * Routeservers shows a list of routeservers for navigation
 */
const RouteServers = () => {
  const routeservers = useRouteservers();
  
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
};


export default RouteServers;
