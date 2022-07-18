
import { useCallback }
  from 'react'

import { useQuery }
  from 'app/components/query';

import SearchInput
  from 'app/components/search/Input';

/**
 * SearchQueryInput is a SearchInput, updating the query.
 */
const SearchQueryInput = ({
  queryKey = "q",
  queryDefault = "",
  debounce=300,
  ...props
}) => {
  const [query, setQuery] = useQuery({[queryKey]: queryDefault});
  const updateQuery = useCallback(
    (v) => setQuery({[queryKey]: v}),
    [setQuery, queryKey]);
  return (
    <SearchInput
      value={query[queryKey]}
      debounce={300}
      onChange={updateQuery}
      {...props}
    />
  );
}

export default SearchQueryInput;
