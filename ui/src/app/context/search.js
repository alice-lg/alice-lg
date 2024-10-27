
import axios from 'axios';

import { useState
       , useEffect
       , useMemo
       , useContext
       , createContext
       }
  from 'react';

import { RoutesReceivedContext
       , RoutesFilteredContext
       , RoutesNotExportedContext
       , paginationState
       , filtersState
       , apiStatusState 
       , useRoutesReceived
       , useRoutesFiltered
       }
  from 'app/context/routes';
import { ApiStatusProvider }
  from 'app/context/api-status';
import { useErrorHandler, isTimeoutError }
  from 'app/context/errors';
import { useQuery
       , PARAM_QUERY
       }
  from 'app/context/query';
import { encodeFilters }
  from 'app/context/filters';


const initialRoutesState = {
  requested: true,
  loading: false,

  page: 0,
  pageSize: 0,
  totalPages: 0,
  totalResults: 0,

  routes: [],
  filtersApplied: [],
  filtersAvailable: [],
  filtersNotAvailable: [],

  apiStatus: {},
};

const initialState = {
  requested: false,
  loading: false,

  error: null,

  received: initialRoutesState,
  filtered: initialRoutesState,

  apiStatus: {},
}

/**
 * useSearchParam retrieves the query parameter
 */
export const useSearchQuery = () => {
  const [query] = useQuery({
    [PARAM_QUERY]: "",
  });
  return query[PARAM_QUERY];
}

/**
 * useRoutesSearchUrl creates a memoized URL for
 * a query with filters.
 */
const useRoutesSearchUrl = ({
  query,
  filters,
  pageReceived,
  pageFiltered,
}) => useMemo(() => {
  const qry = new URLSearchParams({
    ...encodeFilters(filters),
    q: query,
    page_filtered: pageFiltered,
    page_imported: pageReceived,
  }).toString();
  const url = `/api/v1/lookup/prefix?${qry}`;
  return url;
}, [
  query,
  filters,
  pageReceived,
  pageFiltered,
]);


const decodeSearchResult = (result) => {
  const filtered = {
    requested: true,
    loading: false,
    ...paginationState(result.filtered),
    ...apiStatusState(result),
    routes: result?.filtered?.routes,
    filtersApplied: [],
    filtersAvailable: [],
    filtersNotAvailable: [],
  };
  const received = {
    requested: true,
    loading: false,
    ...paginationState(result.imported),
    ...apiStatusState(result),
    ...filtersState(result),
    routes: result?.imported?.routes,
  };
  const state = {
    received: received,
    filtered: filtered,
    ...apiStatusState(result),
  }
  return state;
}


/**
 * useSearchResult retrieves the url and returns the state
 */
const useSearchResults = ({
  query,
  filters,
  pageFiltered,
  pageReceived,
}) => {
  const handleError = useErrorHandler();
  const [state, setState] = useState(initialState);
  const searchUrl = useRoutesSearchUrl({
    query,
    filters,
    pageFiltered,
    pageReceived,
  });

  // Search routes on backend
  useEffect(() => {
    setState((s) => ({
      ...s, 
      received: {
        ...initialRoutesState,
        loading: true,
      },
      filtered: initialRoutesState,
    }));
    axios.get(searchUrl).then(({data}) => {
        setState(decodeSearchResult(data));
      },
      (error) => {
        setState((s) => ({...s,
          error: error,
          received: {
            ...s.received,
            error: error,
            loading: false,
          },
          filtered: {
            ...s.filtered,
            loading: false,
          },
        }));

        // We handle timeout errors ourself. All other errors
        // are handled by the global error handler.
        if(!isTimeoutError(error)) {
          handleError(error);
        }
      });
  }, [searchUrl, setState, handleError]);

  return state;
}


const SearchStatusContext = createContext();

export const useSearchStatus = () => useContext(SearchStatusContext);

export const SearchStatusProvider = ({children, api}) => {
  const received = useRoutesReceived();
  const filtered = useRoutesFiltered();
  
  const context = {
    totalReceived: received.totalResults,
    totalFiltered: filtered.totalResults,
    queryDurationMs: api.request_duration_ms,
  };

  return (
    <SearchStatusContext.Provider value={context}>
      <ApiStatusProvider api={api}>
        {children}
      </ApiStatusProvider>
    </SearchStatusContext.Provider>
  );
}


/**
 * RoutesSearchProvider provides routes received, filtered
 * and not exported.
 */
export const RoutesSearchProvider = ({
  children,
  filters,
  query,
  pageFiltered,
  pageReceived,
}) => {
  const result = useSearchResults({
    query,
    filters,
    pageFiltered,
    pageReceived,
  });

  // RoutesContexts
  return (
    <RoutesFilteredContext.Provider value={result.filtered}>
    <RoutesReceivedContext.Provider value={result.received}>
    <RoutesNotExportedContext.Provider value={initialRoutesState}>
      <SearchStatusProvider api={result.apiStatus}>
        {children}
      </SearchStatusProvider>
    </RoutesNotExportedContext.Provider>
    </RoutesReceivedContext.Provider>
    </RoutesFilteredContext.Provider>
  );
}


