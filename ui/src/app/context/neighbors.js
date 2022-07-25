
import axios from 'axios';

import { createContext
       , useContext
       , useState
       , useEffect
       , useMemo
       }
  from 'react';

import { useErrorHandler }
  from 'app/context/errors';
import { ApiStatusProvider }
  from 'app/context/api-status';

import { isUpState } 
  from 'app/components/neighbors/state';

const initialState = {
  neighbors: [],
  api: {},
  isLoading: true,
};

// Contexts
const NeighborsContext        = createContext();
const NeighborContext         = createContext();
const RelatedNeighborsContext = createContext();


export const useNeighbors        = () => useContext(NeighborsContext);
export const useNeighbor         = () => useContext(NeighborContext);
export const useRelatedNeighbors = () => useContext(RelatedNeighborsContext);


/**
 * useLocalRelatedPeers returns all neighbors on an rs
 * sharing the same ASN and are in state 'up'
 */
export const useLocalRelatedPeers = () => {
  const { neighbors } = useNeighbors();
  const neighbor  = useNeighbor();
  return useMemo(() => {
    if (!neighbor) {
      return [];
    }
    return neighbors.filter((n) =>
      (n.asn === neighbor.asn && isUpState(n.state)));
  }, [neighbors, neighbor])
}

/**
 * NeighborsProvider loads the neighbors for a selected
 * route server identified by id
 */
export const NeighborsProvider = ({children, routeServerId}) => {
  const handleError = useErrorHandler();
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
      <ApiStatusProvider api={state.api}>
        {children}
      </ApiStatusProvider>
    </NeighborsContext.Provider>
  );
}

/**
 * NeighborProvider provides a single neighbor context
 */
export const NeighborProvider = ({neighborId, children}) => {
  const { neighbors } = useNeighbors();
  const neighbor = useMemo(
    () => neighbors.find((n) => n.id === neighborId),
    [neighbors, neighborId]);
  return (
    <NeighborContext.Provider value={neighbor}>
      {children}
    </NeighborContext.Provider>
  );
};

/**
 * RelatedNeighborsProvider fetches all related neighbors
 * from the backend identified the neighbor's ASN.
 */
export const RelatedNeighborsProvider = ({children}) => {
  const neighbor = useNeighbor();
  const handleError = useErrorHandler();
  const [state, setState] = useState(initialState);

  useEffect(() => {
    if (!neighbor) {
      return;
    }
    setState((s) => ({...s, isLoading: true}));
    // Load related neighbors
    const queryUrl = `/api/v1/lookup/neighbors?asn=${neighbor.asn}`;
    axios.get(queryUrl).then(
      ({data}) => {
        setState({
          isLoading: false,
          neighbors: data.neighbors,
        });
      },
      (error) => {
        handleError(error);
        setState((s) => ({...s, isLoading: false}));
      }
    );
  }, [neighbor, handleError]);

  return (
    <RelatedNeighborsContext.Provider value={state}>
      {children}
    </RelatedNeighborsContext.Provider>
  );
}

