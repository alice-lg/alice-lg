
import axios from 'axios';

import { useState
       , useEffect
       , useMemo
       }
  from 'react';

import { RoutesReceivedContext
       , RoutesFilteredContext
       , RoutesNotExportedContext
       , paginationState
       , filtersState
       , apiStatusState 
       }
  from 'app/context/routes';
import { ApiStatusProvider }
  from 'app/context/api-status';
import { useErrorHandler }
  from 'app/context/errors';
import { encodeQuery }
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

  apiStatus: {},
};

const initialState = {
  requested: false,
  loading: false,

  received: initialRoutesState,
  filtered: initialRoutesState,

  apiStatus: {},
}


const useRoutesSearchUrl = ({
  query,
  filters,
  pageReceived,
  pageFiltered,
}) => useMemo(() => {
  const qry = encodeQuery({
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
  };
  console.log(result);
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

export const RoutesSearchProvider = ({
  children,
  filters,
  query,
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
      ...s, received: {...s.received, loading: true},
    }));
    axios.get(searchUrl).then(({data}) => {
        setState(decodeSearchResult(data));
      },
      (error) => handleError(error));
  }, [searchUrl, setState, handleError]);

  console.log(state);

  // RoutesContexts
  return (
    <RoutesFilteredContext.Provider value={state.filtered}>
    <RoutesReceivedContext.Provider value={state.received}>
    <RoutesNotExportedContext.Provider value={initialRoutesState}>
      <ApiStatusProvider api={state.apiStatus}>
        {children}
      </ApiStatusProvider>
    </RoutesNotExportedContext.Provider>
    </RoutesReceivedContext.Provider>
    </RoutesFilteredContext.Provider>
  );
}


