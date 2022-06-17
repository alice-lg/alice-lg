
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

import { useErrors }
  from 'app/components/errors/Provider';

const initialState = {
  asn: 0, // Our own ASN (might be abstracted in the future)

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

const ConfigContext = createContext(null);
export const useConfig = () => useContext(ConfigContext);

const ConfigProvider = ({children}) => {
  const [config, setConfig] = useState(initialState);
  const [handleError] = useErrors();
  
  // OnLoad: once
  useEffect(() => {
    // Fetch config from backend
    axios.get('/api/v1/config').then(
      ({data}) => setConfig(data),
      (error) => handleError(error)
    );
  }, []);

  return (
    <ConfigContext.Provider value={config}>
      {children}
    </ConfigContext.Provider>
  );
}

export default ConfigProvider;
