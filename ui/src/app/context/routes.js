
import axios from 'axios';

import { useState
       , useMemo
       , useEffect
       , useContext
       , createContext
       }
  from 'react';


import { useErrorHandler }
  from 'app/context/errors';
import { ApiStatusProvider }
  from 'app/context/api-status';


export const ROUTES_RECEIVED = "received";
export const ROUTES_FILTERED = "filtered";
export const ROUTES_NOT_EXPORTED = "not-exported";

const ROUTES_PROPERTIES = {
  [ROUTES_RECEIVED]: "imported",
  [ROUTES_FILTERED]: "filtered",
  [ROUTES_NOT_EXPORTED]: "not_exported",
};

// Contexts
const RoutesReceivedContext    = createContext();
const RoutesFilteredContext    = createContext();
const RoutesNotExportedContext = createContext();

export const useRoutesReceived    = () => useContext(RoutesReceivedContext);
export const useRoutesFiltered    = () => useContext(RoutesFilteredContext);
export const useRoutesNotExported = () => useContext(RoutesNotExportedContext);

// Providers

/**
 * Encode routes query url params for a given 
 * request type (filtered, received, not-exported)
 */
const routesQueryUrl = (type) => ({routeServerId, neighborId}) => (params) => {
  const q = new URLSearchParams(params);
  return `/api/v1/routeservers/${routeServerId}` +
    `/neighbors/${neighborId}/routes/${type}?${q.toString()}`;
}


// State: routes, isLoading
const initialState = {
  loading: false,
  requested: false,
  page: 0,
  pageSize: 0,
  totalPages: 0,
  totalResults: 0,
  routes: [],
  apiStatus: {},
};

/**
 * Decode routes state
 */
const paginationState = ({pagination}) => ({
  page: pagination.page,
  pageSize: pagination.page_size,
  totalPages: pagination.total_pages,
  totalResults: pagination.total_results,
})

const filtersState = ({filters_applied, filters_available}) => ({
  filtersApplied: filters_applied,
  filtersAvailable: filters_available,
});

const apiStatusState = ({api}) => ({apiStatus: api});

const routesPayloadState = (type) => (data) => {
  const key = ROUTES_PROPERTIES[type];
  return {
    routes: data[key],
  }
}

const routesState = (type) => (data) => {
  const state = {
    ...paginationState(data),
    ...filtersState(data),
    ...apiStatusState(data),
    ...routesPayloadState(type)(data),
  };
  return state;
}


const createFetchRoutesState = (type) => ({
  routeServerId,
  neighborId,
  page,
  query,
  enabled = true,
}) => {
  const [state, setState] = useState(initialState);
  const handleError = useErrorHandler();

  const url = useMemo(() =>
    routesQueryUrl(type)({routeServerId, neighborId})({
      page: page,
      q: query,
    }),
    [query, page, neighborId, routeServerId]);

  useEffect(() => {
    if (!enabled) {
      return;
    };

    setState((s) => ({...s, requested: true, loading: true}));
    axios.get(url).then(({data}) => {
        setState({
          ...routesState(type)(data), 
          loading: false,
          requested: true,
        })
      },
      (error) => handleError(error)
    );
  }, [url, handleError, enabled]);

  return state;
}

const useFetchReceivedState = createFetchRoutesState(ROUTES_RECEIVED);
const useFetchFilteredState = createFetchRoutesState(ROUTES_FILTERED);
const useFetchNotExportedState = createFetchRoutesState(ROUTES_NOT_EXPORTED);

/**
 * Create routes provider makes a new routes provider
 * for a given context.
 */
const createRoutesProvider = (Context, useFetchRoutesState) => ({
  routeServerId,
  neighborId,
  children,
  page = 0,
  query = "",
  enabled = true,
}) => {
  const state = useFetchRoutesState({
    routeServerId,
    neighborId,
    page,
    query,
    enabled,
  });
  return (
    <Context.Provider value={state}>
      <ApiStatusProvider api={state.apiStatus}>
        {children}
      </ApiStatusProvider>
    </Context.Provider>
  );
}


/**
 * RoutesReceivedProvider loads all routes recieved for a neighbor
 */
export const RoutesReceivedProvider = createRoutesProvider(
  RoutesReceivedContext,
  useFetchReceivedState,
);


/**
 * RoutesFilteredProvider loads all routes filtered for a neighbor
 */
export const RoutesFilteredProvider = createRoutesProvider(
  RoutesFilteredContext,
  useFetchFilteredState,
);

/**
 * RoutesNotExportedProvider loads all routes not exported for 
 * a neighbor.
 */
export const RoutesNotExportedProvider = createRoutesProvider(
  RoutesNotExportedContext,
  useFetchNotExportedState,
);

/**
 * useRoutesLoading checks if any routes are loading
 */
export const useRoutesLoading = () => {
  const received = useRoutesReceived();
  const filtered = useRoutesFiltered();
  const noexport = useRoutesNotExported();

  return (received.requested && received.loading) ||
    (filtered.requested && filtered.loading) ||
    (noexport.requested && noexport.loading);
}


/**
 * RouteDetails Context
 */
const RouteDetailsContext = createContext();

export const useRouteDetails = () => useContext(RouteDetailsContext);

export const useSetRouteDetails = () => useRouteDetails()[1];

export const RouteDetailsProvider = ({children}) => {
  const state = useState();
  return (
    <RouteDetailsContext.Provider value={state}>
      {children}
    </RouteDetailsContext.Provider>
  );
}

