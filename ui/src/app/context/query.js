
import { useMemo
       , useCallback
       }
  from 'react';

import { useSearchParams
       , useLocation
       }
  from 'react-router-dom';

/**
 * Well-known parameters
 */
export const PARAM_QUERY = "q";
export const PARAM_PAGE_FILTERED = "pf";
export const PARAM_PAGE_RECEIVED = "pr";
export const PARAM_PAGE_NOT_EXPORTED = "pn";
export const PARAM_LOAD_NOT_EXPORTED = "ne";
export const PARAM_ORDER = "o";
export const PARAM_SORT = "s";

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

/**
 * cleanParams removes all parameters without a
 * value.
 */
export const cleanParams = (params) => {
  let filtered = {};
  for (const k in params) {
    if (!params[k]) {
      continue;
    }
    filtered[k] = params[k];
  }
  return filtered;
}

/**
 * encodeQuery makes URLSearchParams from
 * cleaned params.
 */
export const encodeQuery = (params) => {
  return new URLSearchParams(cleanParams(params));
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
    let next;
    if (typeof(q) === "function") {
      next = q(query); // Inject current parameters
    } else {
      next = q;
    }
    setParams(encodeQuery(next));
  }, [query, setParams]);

  return [query, setQuery];
}

/**
 * useQueryLink is an alternative to useQuery where
 * instead of a navigation function a location object
 * is created, which can be passed to a Link
 */
export const useMakeQueryLocation = () => {
  const location = useLocation();
  const makeLocation = useCallback((q, overrides={}) => {
    const next = encodeQuery(q);
    return {
      ...location,
      ...overrides,
      search: next.toString()
    };
  }, [location]);
  return makeLocation;
}

export const useQueryLocation = (query, overrides={}) => {
  const makeLocation = useMakeQueryLocation();
  return useMemo(() => makeLocation(query, overrides), [
    makeLocation, query, overrides,
  ]);
};

