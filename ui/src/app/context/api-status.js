
import { createContext
       , useContext
       }
  from 'react';

import { parseServerTime }
  from 'app/components/datetime/time';

const ApiStatusContext = createContext();

export const useApiStatus = () => useContext(ApiStatusContext);


/**
 * Provide API status like cache information
 * and version to downstream components
 */
export const ApiStatusProvider = ({children, api}) => {
  let ctx = {};
  const cachedAt = api?.cache_status?.cached_at;
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
      requestDurationMs: api.request_duration_ms,
      store: api.store_status,
    };
  }

  return (
    <ApiStatusContext.Provider value={ctx}>
      {children}
    </ApiStatusContext.Provider>
  );
}

