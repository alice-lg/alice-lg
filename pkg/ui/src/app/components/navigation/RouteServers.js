
import axios
  from 'axios';

import { useEffect
       , useState
       }
  from 'react';
import { Link }
  from 'react-router-dom';

import { useErrors }
  from 'app/components/errors/Provider';
import { useRouteServers }
  from 'app/components/routeservers/Provider';


const Status = ({routeServerId}) => {
  const [status, setStatus] = useState({
    backend: "",
    version: "",
  });
  const [error, setError] = useState(null);
  const [handleError] = useErrors();

  useEffect(() => {
    axios.get(`/api/v1/routeservers/${routeServerId}/status`)
      .then(
        ({data}) => setStatus(data.status),
        (error)  => {
          handleError(error);
          setError(error); // Local error display
        });
  }, [routeServerId, handleError, setError]);


  if (error && error.code >= 100 && error.code < 200) {
    return (
      <div className="routeserver-status">
        <div className="api-error">
          Unreachable
        </div>
      </div>
    );
  } else if (error) {
    return (
      <div className="routeserver-status">
        <div className="api-error">
          {error.response?.data?.tag}
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
                   onChange={(group) => setSelectedGroup(group)} />
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
