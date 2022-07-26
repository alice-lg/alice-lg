
import { useCallback
       , forwardRef
       }
  from 'react'

import { useQuery
       , PARAM_LOAD_NOT_EXPORTED
       }
  from 'app/context/query';

import SearchInput
  from 'app/components/search/SearchInput';

/**
 * SearchQueryInput is a SearchInput, updating the query.
 */
const SearchQueryInput = forwardRef(({
  queryKey = "q",
  queryDefault = "",
  debounce=300,
  ...props
}, ref) => {
  const [query, setQuery] = useQuery({
    [queryKey]: queryDefault,
  });
  const updateQuery = useCallback(
    (v) => setQuery((q) => ({
      [queryKey]: v,
      [PARAM_LOAD_NOT_EXPORTED]: q[PARAM_LOAD_NOT_EXPORTED], // Keep state
    })),
    [setQuery, queryKey]);
  return (
    <SearchInput
      value={query[queryKey]}
      debounce={300}
      onChange={updateQuery}
      ref={ref}
      {...props}
    />
  );
});

export default SearchQueryInput;
