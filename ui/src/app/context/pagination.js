import { useCallback }
  from 'react';

import { useQuery
       , useMakeQueryLocation
       , PARAM_PAGE_FILTERED
       , PARAM_PAGE_RECEIVED
       , PARAM_PAGE_NOT_EXPORTED
       }
  from 'app/context/query';


/**
 * usePageQuery retrievs the pagination
 * query paramters and decodes the value
 */
export const usePageQuery = () => {
  const [query] = useQuery({
    [PARAM_PAGE_FILTERED]: "0",
    [PARAM_PAGE_RECEIVED]: "0",
    [PARAM_PAGE_NOT_EXPORTED]: "0",
  });
  const filtered = parseInt(query[PARAM_PAGE_FILTERED], 10);
  const received = parseInt(query[PARAM_PAGE_RECEIVED], 10);
  const notExported = parseInt(query[PARAM_PAGE_NOT_EXPORTED], 10);

  return {
    received: received,
    filtered: filtered,
    notExported: notExported,
    [PARAM_PAGE_FILTERED]: filtered,
    [PARAM_PAGE_RECEIVED]: received, 
    [PARAM_PAGE_NOT_EXPORTED]: notExported, 
  };
}

/**
 * usePageLocation creates a location for a given
 * page with a query key
 */
export const useMakePageLocation = (pageKey, anchor) => {
  const [query] = useQuery();
  const makeLocation = useMakeQueryLocation();
  return useCallback((p) =>
    makeLocation({
      ...query,
      [pageKey]: p
    }, {
      hash: anchor,
    }), [
    makeLocation,
    pageKey,
    anchor,
    query,
  ]);
}

