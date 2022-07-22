
import { useMemo
       , useCallback
       }
  from 'react';

import { useSearchParams
       , useLocation
       }
  from 'react-router-dom';


/**
 * paramsToQuery creates an object with query params
 * as key / value pairs for convenient access.
 */
const paramsToQuery = (params) => {
  let q = {};
  for (const [k, v] of params) {
    q[k] = v;
  }
  return q;
}


export const encodeQuery = (params) => {
  let filtered = {};
  for (const k in params) {
    if (params[k] === "") {
      continue;
    }
    filtered[k] = params[k];
  }
  return new URLSearchParams(filtered);
}

/**
 * useQuery is an extension to useLocation to handle
 * query parameters. Internally this uses URLSearchParams
 * for decoding but returns an object merged with the defaults.
 * To prevent loops, the search parameters are only updated
 * if they differ.
 */
export const useQuery = (defaults={}) => {
  const [params, setParams] = useSearchParams(defaults);
  const query = useMemo(() => paramsToQuery(params), [params]);
  const setQuery = useCallback((q) => {
    let next = {...query, ...q};
    const nextParams = encodeQuery(next);
    if (nextParams.toString() !== params.toString()) {
      setParams(nextParams);
    }
  }, [params, query, setParams]);
  return [query, setQuery];
}

/**
 * useQueryParams returns the parameters only
 */
export const useQueryParams = (defaults={}) => {
  const [params] = useSearchParams(defaults);
  const query = useMemo(() => paramsToQuery(params), [params]);
  return query;
}

/**
 * useQueryLink is an alternative to useQuery where
 * instead of a navigation function a location object
 * is created, which can be passed to a Link
 */
export const useQueryLink = (defaults={}) => {
  const location = useLocation();
  const [params] = useSearchParams(defaults);
  const query = useMemo(() => paramsToQuery(params), [params]);

  const makeLocation = useCallback((q) => {
    const next = new URLSearchParams({...query, ...q});
    return {...location, search: next.toString()};
  }, [location, query]);

  return [query, makeLocation];
}

