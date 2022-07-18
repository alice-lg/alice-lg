/**
 * RouteServers provider
 *
 * This provider fetches all route servers and makes them
 * available through a context
 */

import axios from 'axios';

import { useState
       , useEffect
       , useContext
       , createContext
       , useRef
       }
  from 'react';
import { useParams }
  from 'react-router-dom';

import { useErrorHandler }
  from 'app/context/errors';


// Contexts
const RouteServersContext      = createContext([]);
const RouteServerStatusContext = createContext();

export const useRouteServers      = () => useContext(RouteServersContext);
export const useRouteServerStatus = () =>  useContext(RouteServerStatusContext);

/**
 * Use selected route server uses the route server context
 * in combination with the navigation to return the current
 * route server.
 */
export const useRouteServer = () => {
  const { routeServerId } = useParams();
  const routeServers      = useRouteServers();
  return routeServers.find((rs) => rs.id === routeServerId)
}

/**
 * RouteServersProvider loads the route servers from the
 * backend and uses these as provider value.
 */
export const RouteServersProvider = ({children}) => {
  const init          = useRef();
  const handleError   = useErrorHandler();
  const [rs, setRs]   = useState([]);
  
  // Load route servers from backend
  useEffect(() => {
    if (init.current) {
      return;
    }
    axios.get('/api/v1/routeservers')
      .then(
        ({data}) => setRs(data.routeservers),
        (error) => handleError(error)
      );
      init.current = true;
  }, [handleError]);

  return (
    <RouteServersContext.Provider value={rs}>
      {children}
    </RouteServersContext.Provider>
  );
}

/**
 * RouteServerStatusProvider loads the route server status
 * and provides it through the context
 */
export const RouteServerStatusProvider = ({children, routeServerId}) => {
  const handleError         = useErrorHandler();
  const [status, setStatus] = useState({});

  useEffect(() => {
    axios.get(`/api/v1/routeservers/${routeServerId}/status`)
      .then(
        ({data}) => setStatus(data.status),
        (error)  => handleError(error),
      );
  }, [routeServerId, handleError]);

  return (
    <RouteServerStatusContext.Provider value={status}>
      {children}
    </RouteServerStatusContext.Provider>
  );
}
