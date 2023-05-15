
import axios
  from 'axios';

import { useEffect
       , useState
       , useCallback
       }
  from 'react';
import { Link }
  from 'react-router-dom';

import { useErrorHandler }
  from 'app/context/errors';
import { useRouteServers }
  from 'app/context/route-servers';


const Status = ({routeServerId}) => {
  const [status, setStatus] = useState({
    backend: "",
    version: "",
  });
  const [error, setError] = useState(null);
  const handleError = useErrorHandler();

  useEffect(() => {
    axios.get(`/api/v1/routeservers/${routeServerId}/status`)
      .then(
        ({data}) => setStatus(data.status),
        (error)  => {
          setError(error); // Local error display
        });
  }, [routeServerId, handleError]);

  const errorInfo = error?.response?.data;

  if (errorInfo && errorInfo.tag === "CONNECTION_REFUSED") {
    return (
      <div className="routeserver-status">
        <div className="api-error">
          route server unreachable
        </div>
      </div>
    );
  } else if (errorInfo && errorInfo.tag === "GENERIC_ERROR") {
    return (
      <div className="routeserver-status">
        <div className="api-error">
          did not respond
        </div>
      </div>
    );
  } else if (errorInfo) {
    return (
      <div className="routeserver-status">
        <div className="api-error">
          {errorInfo.tag}
        </div>
      </div>
    );
  }

  return (
    <div className="routeserver-status">
      <div className="bird-version">
        {status.backend} {status.version}
      </div>
    </div>
  );
};


/**
 * Select a routeserver button
 */
const GroupSelectOption = ({group, onSelect}) => {
  const selectGroup = useCallback(() => onSelect(group), [
    group, onSelect,
  ]);
  return (
    <li>
      <button className="btn btn-link btn-option" onClick={selectGroup}>
        {group}
      </button>
    </li>
  );
}

/**
 * GroupSelect shows a drop down for selecting a
 * group of routeservers.
 */
const GroupSelect = ({groups, selected, onSelect}) => {
  const [expanded, setExpanded] = useState(false);

  const toggleDropdown = useCallback(() => {
    setExpanded((state) => !state);
  }, []);

  const selectGroup = useCallback((group) => {
    onSelect(group);
    setExpanded(false);
  }, [onSelect]);

  if (groups.length < 2) {
    return null; // why bother?
  }

  const options = groups.map((group) =>
    <GroupSelectOption key={group} group={group} onSelect={selectGroup} />
  );

  const dropdownClass = `dropdown ${expanded && 'open'}`;

  return (
    <div className="routeservers-groups-select">
      <div className={dropdownClass}>
        <button className="btn btn-default dropdown-toggle btn-select"
                type="button"
                id="select-routeservers-group"
                onClick={toggleDropdown}
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
  const routeServers = useRouteServers();
  const [selectedGroup, setSelectedGroup] = useState(null);

  let groups = [];
  for (const rs of routeServers) {
    if (groups.indexOf(rs.group) === -1) {
      groups.push(rs.group);
    }
  }

  useEffect(() => {
    setSelectedGroup(routeServers[0]?.group);
  }, [routeServers])

  if (selectedGroup === null) {
    return null; // nothing to display yet
  }
  
  const groupRs = routeServers.filter((rs) => rs.group === selectedGroup);
  
  return (
    <div className="routeservers-list">
      <h2>route servers</h2>
      <GroupSelect groups={groups}
                   selected={selectedGroup}
                   onSelect={setSelectedGroup} />
      <ul>
        {groupRs.map((rs) =>
          <li key={rs.id}>
            <Link to={`/routeservers/${rs.id}`}
                  className="routeserver-id">{rs.name}</Link>
            <Status routeServerId={rs.id} />
          </li>)}
      </ul>
    </div>
  );
};


export default RouteServers;
