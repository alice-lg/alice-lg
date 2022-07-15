import { useContext
       , createContext
       , useMemo
       }
  from 'react';

import { useNeighbors }
  from 'app/components/neighbors/Provider';


const NeighborContext = createContext();

export const useNeighbor = () => useContext(NeighborContext);

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


export default NeighborProvider;
