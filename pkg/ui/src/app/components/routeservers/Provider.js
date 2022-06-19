/**
 * RouteServers provider
 *
 * This provider fetches all routeservers and makes them
 * available through a context
 */

import axios from 'axios';

import { useState
       , useEffect
       , useContext
       , createContext
       }
  from 'react';

import { useErrors }
  from 'app/components/errors/Provider';


const RouteserversContext = createContext([]);

export const useRouteservers = () => useContext(RouteserversContext);


/**
 * RouteserversProvider loads the routeservers from the
 * backend and uses these as provider value.
 */
const RouteserversProvider = ({children}) => {
  const [handleError] = useErrors();
  const [rs, setRs]   = useState([]);
  
  // Load routeservers from backend
  useEffect(() => {
    axios.get('/api/v1/routeservers')
      .then(
        ({data}) => setRs(data),
        (error) => handleError(error)
      );
  }, [handleError]);

  return (
    <RouteserversContext.Provider value={rs}>
      {children}
    </RouteserversContext.Provider>
  );
}

export default RouteserversProvider;
