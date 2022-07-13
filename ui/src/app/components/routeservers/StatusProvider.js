
import axios from 'axios';

import { useEffect
       , useState
       , useContext
       , createContext
       }
  from 'react';

import { useErrorHandler }
  from 'app/components/errors/Provider';

const RouteServerStatusContext = createContext();

export const useRouteServerStatus = () =>  useContext(RouteServerStatusContext);


/**
 * RouteServerStatusProvider loads the route server status
 * and provides it through the context
 */
const RouteServerStatusProvider = ({children, routeServerId}) => {
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

export default RouteServerStatusProvider;
