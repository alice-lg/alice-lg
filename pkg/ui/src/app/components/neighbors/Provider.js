
import axios from 'axios';

import { createContext
       , useContext
       , useState
       , useEffect
       }
  from 'react';

import { useErrors }
  from 'app/components/errors/Provider';

import CacheStatusProvider
  from 'app/components/cache/StatusProvider';

const initialState = {
  neighbors: [],
  api: {},
  isLoading: true,
};

const NeighborsContext = createContext();

export const useNeighbors = () => useContext(NeighborsContext);

/**
 * NeighborsProvider loads the neighbors for a selected
 * route server identified by id
 */
const NeighborsProvider = ({children, routeServerId}) => {
  const [handleError]     = useErrors();
  const [state, setState] = useState(initialState);

  useEffect(() => {
    setState((s) => ({...s, isLoading: true}));
    // Load RouteServer's neighbors
    axios.get(`/api/v1/routeservers/${routeServerId}/neighbors`).then(
      ({data}) => {
        setState({
          isLoading: false,
          neighbors: data.neighbors,
          api: data.api,
        });
      },
      (error) => {
        handleError(error);
        setState((s) => ({...s, isLoading: false}));
      }
    );
  }, [routeServerId, handleError]);

  return (
    <NeighborsContext.Provider value={state}>
      <CacheStatusProvider api={state.api}>
        {children}
      </CacheStatusProvider>
    </NeighborsContext.Provider>
  );
}

export default NeighborsProvider;
