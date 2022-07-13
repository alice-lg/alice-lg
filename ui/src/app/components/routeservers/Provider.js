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

import { useErrors }
  from 'app/components/errors/Provider';


const RouteServersContext = createContext([]);

export const useRouteServers = () => useContext(RouteServersContext);

/**
 * Use selected route server uses the route server context
 * in combination with the navigation to return the current
 * route server.
 */
export const useSelectedRouteServer = () => {
  const { routeServerId } = useParams();
  const routeServers      = useRouteServers();
  return routeServers.find((rs) => rs.id === routeServerId)
}


/**
 * RouteServersProvider loads the route servers from the
 * backend and uses these as provider value.
 */
const RouteServersProvider = ({children}) => {
  const init          = useRef();
  const [handleError] = useErrors();
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

export default RouteServersProvider;
