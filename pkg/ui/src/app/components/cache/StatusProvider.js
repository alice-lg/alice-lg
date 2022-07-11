
import { createContext
       , useContext
       }
  from 'react';

import { parseServerTime }
  from 'app/components/datetime/time';

const CacheStatusContext = createContext();

export const useCacheStatus = () => useContext(CacheStatusContext);


const CacheStatusProvider = ({children, api}) => {
  let ctx = null;

  const cachedAt = api.cache_status?.cached_at;
  if (cachedAt) {
    const ttl = parseServerTime(api.ttl);
    const generatedAt = parseServerTime(cachedAt);
    const age = ttl.diff(generatedAt); // ms

    // Create cache status from API metadata
    ctx = {
      resultFromCache: api.result_from_cache,
      ttl: ttl,
      ttlTime: api.ttl,
      version: api.version,
      cachedAt: api.cache_status?.cached_at,
      origTtl: api.cache_status?.orig_ttl,
      generatedAt: generatedAt,
      age: age,
    };
  }

  return (
    <CacheStatusContext.Provider value={ctx}>
      {children}
    </CacheStatusContext.Provider>
  );
}

export default CacheStatusProvider;
