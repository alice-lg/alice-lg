
import { useMemo
       , useContext
       , createContext
       , useCallback
       }
  from 'react';

import { useRoutesReceived
       , useRoutesFiltered
       , useRoutesNotExported
       }
  from 'app/context/routes';
import { useQuery }
  from 'app/context/query';


/**
 * Filter Groups
 */
const FILTER_KEY_SOURCES = "sources";
const FILTER_KEY_ASNS = "asns";
const FILTER_KEY_COMMUNITIES = "communities";
const FILTER_KEY_EXT_COMMUNITIES = "ext_communities";
const FILTER_KEY_LARGE_COMMUNITIES = "large_communities";

export const FILTER_GROUP_SOURCES = 0;
export const FILTER_GROUP_ASNS = 1;
export const FILTER_GROUP_COMMUNITIES = 2;
export const FILTER_GROUP_EXT_COMMUNITIES = 3;
export const FILTER_GROUP_LARGE_COMMUNITIES = 4;

const FILTER_GROUP_KEYS = {
  [FILTER_GROUP_SOURCES]: FILTER_KEY_SOURCES,
  [FILTER_GROUP_ASNS]: FILTER_KEY_ASNS,
  [FILTER_GROUP_COMMUNITIES]: FILTER_KEY_COMMUNITIES,
  [FILTER_GROUP_EXT_COMMUNITIES]: FILTER_KEY_EXT_COMMUNITIES,
  [FILTER_GROUP_LARGE_COMMUNITIES]: FILTER_KEY_LARGE_COMMUNITIES,
};


/**
 * Make initialState
 */
export const initializeFilterState = () => ([
  {"key": "sources", "filters": []},
  {"key": "asns", "filters": []},
  {"key": "communities", "filters": []},
  {"key": "ext_communities", "filters": []},
  {"key": "large_communities", "filters": []},
]);



// Compare values
const cmpValue = (a, b) => a === b;

// Compare list values
const cmpList = (a, b) => 
  a.map((v, i) => v === b[i]).reduce(
    (part, match) => (match && part), true);


const FILTER_VALUE_CMP = {
  [FILTER_GROUP_SOURCES]: cmpValue,
  [FILTER_GROUP_ASNS]: cmpValue,
  [FILTER_GROUP_COMMUNITIES]: cmpList,
  [FILTER_GROUP_EXT_COMMUNITIES]: cmpList,
  [FILTER_GROUP_LARGE_COMMUNITIES]: cmpList,
}

/*
 * Filters set compare
 */
const cmpFilterValue = (set, filter) => {
  for (const f of set) {
    if(f.value === filter.value) {
      return f;
    }
  }
  return null;
}

const cmpFilterCommunity = (set, filter) => {
  for (const f of set) {
    let match = true;
    for (const i in f.value) {
      if (f.value[i] !== filter.value[i]) {
        match = false;
        break;
      }
    }
    if (match) {
      return f;
    }
  }
  return null;
}



/*
 * Merge list of filters
 */
const mergeFilterSet = (inSet, a, b) => {
  let result = a;
  for (const f of b) {
    const present = inSet(result, f);
    if (present) {
      // Update filter cardinality
      present.cardinality += f.cardinality;
      continue;
    }
    result.push(f);
  }
  return result;
}


/*
 * Merge filters
 */
const _mergeFilters = (a, b) => {
  let groups = initializeFilterState();
  if (!a || !b) {
    return groups;
  }
  let setCmp = [];
  setCmp[FILTER_GROUP_SOURCES] = cmpFilterValue;
  setCmp[FILTER_GROUP_ASNS] = cmpFilterValue;
  setCmp[FILTER_GROUP_COMMUNITIES] = cmpFilterCommunity;
  setCmp[FILTER_GROUP_EXT_COMMUNITIES] = cmpFilterCommunity;
  setCmp[FILTER_GROUP_LARGE_COMMUNITIES] = cmpFilterCommunity;
  for (const i in groups) {
    if (a[i]?.filters && b[i]?.filters) {
      groups[i].filters = mergeFilterSet(setCmp[i], a[i].filters, b[i].filters);
    }
    else if(a[i]?.filters) {
      groups[i].filters = a[i].filters;
    }
    else if(b[i]?.filters) {
      groups[i].filters = b[i].filters;
    }
  }
  return groups;
}

const mergeFilters = (a, ...other) => {
  let result = a;
  for (const filters of other) {
    result = _mergeFilters(result, filters);
  }
  return result;
}

/*
 * Does a single group have any filters?
 */
const groupHasFilters = (group) =>
  group.filters.length > 0;

/*
 * Do we have filters in general?
 */
const hasFilters = (groups) => {
  for (const g of groups) {
    if (groupHasFilters(g)) {
      return true;
    }
  }
  return false;
}


/*
 * Filter Query Decoding
 */
const decodeStringList = (value) => {
  if (value === "") {
    return [];
  }
  return value.split(",");
}

const decodeIntList = (value) =>
  decodeStringList(value).map((v) =>
    parseInt(v, 10));

const decodeCommunity = (community) =>
  community.split(":").map((p) =>
    parseInt(p, 10));

const decodeExtCommunity = (community) =>
  community.split(":");

const decodeCommunities = (value) =>
  decodeStringList(value).map((v) =>
    decodeCommunity(v));

const decodeExtCommunities = (value) =>
  decodeStringList(value).map((v) =>
    decodeExtCommunity(v));


const decodeQuery = (query) => {
  const sources = decodeStringList(query[FILTER_KEY_SOURCES]);
  const asns = decodeIntList(query[FILTER_KEY_ASNS]);
  const communities = decodeCommunities(query[FILTER_KEY_COMMUNITIES]);
  const extCommunities = decodeExtCommunities(query[FILTER_KEY_EXT_COMMUNITIES]);
  const largeCommunities = decodeCommunities(query[FILTER_KEY_LARGE_COMMUNITIES]);
  let filters = {};
  filters[FILTER_KEY_SOURCES] = sources;
  filters[FILTER_KEY_ASNS] = asns;
  filters[FILTER_KEY_COMMUNITIES] = communities;
  filters[FILTER_KEY_EXT_COMMUNITIES] = extCommunities;
  filters[FILTER_KEY_LARGE_COMMUNITIES] = largeCommunities;
  return filters;
}

/*
 * Filter Query Encoding
 */
const encodeList = (value) =>
  value.join(",");

const encodeCommunity = (community) =>
  community.join(":");

const encodeCommunities = (communities) =>
  encodeList(communities.map((c) => encodeCommunity(c)));

export const encodeFilters = (filters) => {
  let query = {};
  const sources = filters[FILTER_KEY_SOURCES];
  const asns = filters[FILTER_KEY_ASNS];
  const communities = filters[FILTER_KEY_COMMUNITIES];
  const extCommunities = filters[FILTER_KEY_EXT_COMMUNITIES];
  const largeCommunities = filters[FILTER_KEY_LARGE_COMMUNITIES];
  query[FILTER_KEY_SOURCES] = encodeList(sources);
  query[FILTER_KEY_ASNS] = encodeList(asns);
  query[FILTER_KEY_COMMUNITIES] = encodeCommunities(communities);
  query[FILTER_KEY_EXT_COMMUNITIES] = encodeCommunities(extCommunities);
  query[FILTER_KEY_LARGE_COMMUNITIES] = encodeCommunities(largeCommunities);
  return query;
}


/**
 * FiltersQuery Context
 */
export const useFiltersQuery = () => {
  const [, setQuery] = useQuery();
  const [query] = useQuery({
    [FILTER_KEY_SOURCES]: "",
    [FILTER_KEY_ASNS]: "",
    [FILTER_KEY_COMMUNITIES]: "",
    [FILTER_KEY_EXT_COMMUNITIES]: "",
    [FILTER_KEY_LARGE_COMMUNITIES]: "",
  });

  const filterQuery = useMemo(() => decodeQuery(query), [query]);
  const setFilterQuery = useCallback((key, value) => {
    const next = {...filterQuery, [key]: value};
    setQuery(encodeFilters(next));
  }, [filterQuery, setQuery]);

  return [filterQuery, setFilterQuery];
};


/*
 * FiltersContext
 */
const FiltersContext = createContext();

export const useFilters = () => useContext(FiltersContext);

const useRoutesFilters = (routes) => {
  return useMemo(() => {
    if (!routes.requested || routes.loading) {
      return { applied: [], available: [] };
    }
    return {
      applied: routes.filtersApplied,
      available: routes.filtersAvailable,
    };
  }, [routes]);
}

export const useTotalFilters = () => {
  const {filters} = useFilters();
  const {applied, available} = filters;
  return useMemo(() => 
    applied.reduce(
      (total, group) => total + group.filters.length,
      0,
    ) + available.reduce(
      (total, group) => total + group.filters.length,
      0
    ), [applied, available]);
}


const createGroupFilters = (group) => () => {
  const {filters, applyFilter, removeFilter} = useFilters();
  const applyGroupFilter = useCallback((filters) => {
    applyFilter(group, filters);
  }, [applyFilter]);
  const removeGroupFilter = useCallback((filters) => {
    removeFilter(group, filters);
  }, [removeFilter]);
  return useMemo(() => ({
    filters: {
      applied:   filters.applied[group].filters,
      available: filters.available[group].filters,
    },
    applyFilter: applyGroupFilter,
    removeFilter: removeGroupFilter,
  }), [filters, applyGroupFilter, removeGroupFilter]);
}

export const useSourceFilters = createGroupFilters(
  FILTER_GROUP_SOURCES);
export const useAsnFilters = createGroupFilters(
  FILTER_GROUP_ASNS);
export const useCommunitiesFilters = createGroupFilters(
  FILTER_GROUP_COMMUNITIES);
export const useExtCommunitiesFilters = createGroupFilters(
  FILTER_GROUP_EXT_COMMUNITIES);
export const useLargeCommunitiesFilters = createGroupFilters(
  FILTER_GROUP_LARGE_COMMUNITIES);

/**
 * RoutesFiltersProvider merged the filters from the
 * received, filtered and noexport responses
 */
export const RoutesFiltersProvider = ({children}) => {
  const received = useRoutesFilters(useRoutesReceived());
  const filtered = useRoutesFilters(useRoutesFiltered());
  const noexport = useRoutesFilters(useRoutesNotExported());

  const [queryFilters, setFilterQuery] = useFiltersQuery();

  const applyFilter = useCallback((group, value) => {
    const key = FILTER_GROUP_KEYS[group];
    const values = [...queryFilters[key], value];
    setFilterQuery(key, values);
  }, [queryFilters, setFilterQuery]);

  const removeFilter = useCallback((group, value) => {
    const cmp = FILTER_VALUE_CMP[group];
    const key = FILTER_GROUP_KEYS[group];
    const values = queryFilters[key].filter((f) => !cmp(f, value));
    setFilterQuery(key, values);
  }, [queryFilters, setFilterQuery]);

  const filters = useMemo(() => {
    const applied = mergeFilters(
      received.applied,
      filtered.applied,
      noexport.applied,
    );
    const available = mergeFilters(
      received.available,
      filtered.available,
      noexport.available,
    );
    return { applied, available };
  }, [received, filtered, noexport]);

  const context = {filters, applyFilter, removeFilter};
  return (
    <FiltersContext.Provider value={context}>
      {children}
    </FiltersContext.Provider>
  );
};

