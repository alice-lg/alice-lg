
/*
 * The ConfigProvider fetches the runtime configuration
 * from the backend and provides it through useContext.
 */

import axios
  from 'axios';

import { createContext
       , useContext
       , useState
       , useEffect
       }
  from 'react';

import { useErrorHandler }
  from 'app/context/errors';

const initialState = {
  routes_columns: {},
  routes_columns_order: [],
  neighbors_columns: {},
  neighbors_columns_order: [],
  lookup_columns: {},
  lookup_columns_order: [],
  prefix_lookup_enabled: false,
  content: {},
  noexport_load_on_demand: true, // we have to assume this
                                 // otherwise fetch will start.
  rpki: { 
    enabled: false,
  },

  bgp_communities: {},

  blackholes: {}, // Map blackholes to routeservers
  asns: {}, // Map ASNs to routeservers (for future use)
};

export const ConfigContext = createContext(null);
export const useConfig = () => useContext(ConfigContext);

export const ConfigProvider = ({children}) => {
  const [config, setConfig] = useState(initialState);
  const handleError = useErrorHandler();
  
  // OnLoad: once
  useEffect(() => {
    // Fetch config from backend
    axios.get('/api/v1/config').then(
      ({data}) => setConfig(data),
      (error) => handleError(error)
    );
  }, [handleError]);

  return (
    <ConfigContext.Provider value={config}>
      {children}
    </ConfigContext.Provider>
  );
}


/**
 * RoutesTableConfigProvider
 */
const RoutesTableConfigContext = createContext();

export const useRoutesTableConfig = () => useContext(RoutesTableConfigContext);

/**
 * Configure routes columns and columns oder
 */
export const RoutesTableConfigProvider = ({
  children,
  columns,
  columnsOrder
}) => {
  const context = {
    columns,
    columnsOrder,
  };
  return (
    <RoutesTableConfigContext.Provider value={context}>
      {children}
    </RoutesTableConfigContext.Provider>
  );
}


