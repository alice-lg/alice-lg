
import { useMemo
       , useContext
       , createContext
       }
  from 'react';

import { useRoutesReceived
       , useRoutesFiltered
       , useRoutesNotExported
       }
  from 'app/context/routes';


/**
 * Filter Groups
 */
const FILTER_KEY_SOURCES = "sources";
const FILTER_KEY_ASNS = "asns";
const FILTER_KEY_COMMUNITIES = "communities";
const FILTER_KEY_EXT_COMMUNITIES = "ext_communities";
const FILTER_KEY_LARGE_COMMUNITIES = "large_communities";

const FILTER_GROUP_SOURCES = 0;
const FILTER_GROUP_ASNS = 1;
const FILTER_GROUP_COMMUNITIES = 2;
const FILTER_GROUP_EXT_COMMUNITIES = 3;
const FILTER_GROUP_LARGE_COMMUNITIES = 4;


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


/**
 * Check filters for equality
 */
const filtersEqual = (a, b) => {
  return (a[FILTER_GROUP_SOURCES].filters.length ===
          b[FILTER_GROUP_SOURCES].filters.length) &&

         (a[FILTER_GROUP_ASNS].filters.length ===
          b[FILTER_GROUP_ASNS].filters.length) &&

         (a[FILTER_GROUP_COMMUNITIES].filters.length ===
          b[FILTER_GROUP_COMMUNITIES].filters.length) &&

         (a[FILTER_GROUP_EXT_COMMUNITIES].filters.length ===
          b[FILTER_GROUP_EXT_COMMUNITIES].filters.length) &&

         (a[FILTER_GROUP_LARGE_COMMUNITIES].filters.length ===
          b[FILTER_GROUP_LARGE_COMMUNITIES].filters.length);
}


/**
 * Deep copy of filters
 */
const cloneFilters = (filters) => {
  return filters;

  /*
  const nextFilters = [
    {...filters[FILTER_GROUP_SOURCES]},
    {...filters[FILTER_GROUP_ASNS]},
    {...filters[FILTER_GROUP_COMMUNITIES]},
    {...filters[FILTER_GROUP_EXT_COMMUNITIES]},
    {...filters[FILTER_GROUP_LARGE_COMMUNITIES]},
  ];
  nextFilters[FILTER_GROUP_SOURCES].filters =
    [...[FILTER_GROUP_SOURCES].filters];
  nextFilters[FILTER_GROUP_ASNS].filters =
    [...nextFilters[FILTER_GROUP_ASNS].filters];
  nextFilters[FILTER_GROUP_COMMUNITIES].filters =
    [...nextFilters[FILTER_GROUP_COMMUNITIES].filters];
  nextFilters[FILTER_GROUP_EXT_COMMUNITIES].filters =
    [...nextFilters[FILTER_GROUP_EXT_COMMUNITIES].filters];
  nextFilters[FILTER_GROUP_LARGE_COMMUNITIES].filters =
    [...nextFilters[FILTER_GROUP_LARGE_COMMUNITIES].filters];
  return nextFilters;
  */
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
 * Filter Query Encoding
 */

const makeFilter = (value) => {
  return {
    name: "",
    value: value,
    cardinality: 1,
  }
}

const decodeFiltersSources = ({sources}) => {
  if (!sources) {
    return []; // No params available
  }
  return sources.split(",").map(
    (sid) => makeFilter(sid));
}

const decodeFiltersAsns = ({asns}) => {
  if (!asns) {
    return []; // No params available
  }
  return asns.split(",").map(
    (asn) => makeFilter(parseInt(asn, 10)));
}

const decodeCommunity = (community) => {
  const parts = community.split(":");
  return parts.map((p) => parseInt(p, 10));
}

const decodeExtCommunity = (community) =>
  community.split(":");

const decodeFiltersCommunities = ({communities}) => {
  if (!communities) {
    return []; // No params available
  }
  return communities.split(",").map(
    (c) => makeFilter(decodeCommunity(c)));
}

const decodeFiltersExtCommunities = ({ext_communities}) => {
  if (!ext_communities) {
    return []; // No params available
  }
  const communities = ext_communities.split(",");
  return communities.map(
    (c) => makeFilter(decodeExtCommunity(c)));
}

const decodeFiltersLargeCommunities = ({large_communities}) => {
  if (!large_communities) {
    return []; // No params available
  }
  const communities = large_communities.split(",");
  return communities.map(
    (c) => makeFilter(decodeCommunity(c)));
}

const encodeGroupInt = (group) => {
  if (!group.filters.length) {
    return "";
  }
  const values = group.filters.map((f) => f.value).join(",");
  return `&${group.key}=${values}`;
}

const encodeGroupCommunities = (group) => {
  if (!group.filters.length) {
    return "";
  }
  const values = group.filters.map((f) => f.value.join(":")).join(",");
  return `&${group.key}=${values}`;
}

/**
 * Encode filters as URL params
 */
const filtersUrlEncode = (filters) => {
  let encoded = "";
  encoded += encodeGroupInt(filters[FILTER_GROUP_SOURCES]);
  encoded += encodeGroupInt(filters[FILTER_GROUP_ASNS]);
  encoded += encodeGroupCommunities(filters[FILTER_GROUP_COMMUNITIES]);
  encoded += encodeGroupCommunities(filters[FILTER_GROUP_EXT_COMMUNITIES]);
  encoded += encodeGroupCommunities(filters[FILTER_GROUP_LARGE_COMMUNITIES]);
  return encoded;
}

/*
 * Decode filters applied from params
 */
const decodeFiltersApplied = (params) => {
  const groups = initializeFilterState();

  groups[FILTER_GROUP_SOURCES].filters =           decodeFiltersSources(params);
  groups[FILTER_GROUP_ASNS].filters =              decodeFiltersAsns(params);
  groups[FILTER_GROUP_COMMUNITIES].filters =       decodeFiltersCommunities(params);
  groups[FILTER_GROUP_EXT_COMMUNITIES].filters =   decodeFiltersExtCommunities(params);
  groups[FILTER_GROUP_LARGE_COMMUNITIES].filters = decodeFiltersLargeCommunities(params);

  return groups;
}


/*
 * FiltersContext
 */
const initialContext = {
  applied: [],
  available: [],
};

const FiltersContext = createContext(initialContext);

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

/**
 * RoutesFiltersProvider merged the filters from the
 * received, filtered and noexport responses
 */
export const RoutesFiltersProvider = ({children}) => {
  const received = useRoutesFilters(useRoutesReceived());
  const filtered = useRoutesFilters(useRoutesFiltered());
  const noexport = useRoutesFilters(useRoutesNotExported());

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

  return (
    <FiltersContext.Provider value={filters}>
      {children}
    </FiltersContext.Provider>
  );
};

