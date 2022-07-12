
import { useMemo
       , useCallback
       }
  from 'react';

import { useSearchParams
       }
  from 'react-router-dom';


/**
 * useQuery is an extension to useLocation to handle
 * query parameters. Internally this uses URLSearchParams
 * for decoding but returns an object merged with the defaults.
 * To prevent loops, the search parameters are only updated
 * if they differ.
 */
export const useQuery = (defaults={}) => {
  const [query, setQuery] = useSearchParams(defaults);
  const params = useMemo(() => {
    // For convenient access convert params to object
    let q = {};
    for (const [k, v] of query) {
      q[k] = v;
    }
    return q;
  }, [query]);

  const update = useCallback((q) => {
    // Only update if query differs
    const next = new URLSearchParams({...params, ...q});
    if (next.toString() !== query.toString()) {
      setQuery(next);
    }
  }, [params, query, setQuery]);
  return [params, update];
}

