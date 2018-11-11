
import _ from 'underscore'

import {FILTER_GROUP_SOURCES,
        FILTER_GROUP_ASNS,
        FILTER_GROUP_COMMUNITIES,
        FILTER_GROUP_EXT_COMMUNITIES,
        FILTER_GROUP_LARGE_COMMUNITIES}
  from './groups'

import {decodeFiltersSources,
        decodeFiltersAsns,
        decodeFiltersCommunities,
        decodeFiltersExtCommunities,
        decodeFiltersLargeCommunities}
  from 'components/filters/encoding'

export const initializeFilterState = () => ([
  {"key": "sources", "filters": []},
  {"key": "asns", "filters": []},
  {"key": "communities", "filters": []},
  {"key": "ext_communities", "filters": []},
  {"key": "large_communities", "filters": []},
]);

export function cloneFilters(filters) {
  const nextFilters = [
    Object.assign({}, filters[FILTER_GROUP_SOURCES]),
    Object.assign({}, filters[FILTER_GROUP_ASNS]),
    Object.assign({}, filters[FILTER_GROUP_COMMUNITIES]),
    Object.assign({}, filters[FILTER_GROUP_EXT_COMMUNITIES]),
    Object.assign({}, filters[FILTER_GROUP_LARGE_COMMUNITIES]),
  ];

  nextFilters[FILTER_GROUP_SOURCES].filters =
    [...nextFilters[FILTER_GROUP_SOURCES].filters];

  nextFilters[FILTER_GROUP_ASNS].filters =
    [...nextFilters[FILTER_GROUP_ASNS].filters];

  nextFilters[FILTER_GROUP_COMMUNITIES].filters =
    [...nextFilters[FILTER_GROUP_COMMUNITIES].filters];

  nextFilters[FILTER_GROUP_EXT_COMMUNITIES].filters =
    [...nextFilters[FILTER_GROUP_EXT_COMMUNITIES].filters];

  nextFilters[FILTER_GROUP_LARGE_COMMUNITIES].filters =
    [...nextFilters[FILTER_GROUP_LARGE_COMMUNITIES].filters];

  return nextFilters;
}

/*
 * Decode filters applied from params
 */
export function decodeFiltersApplied(params) {
  const groups = initializeFilterState();

  groups[FILTER_GROUP_SOURCES].filters =           decodeFiltersSources(params);
  groups[FILTER_GROUP_ASNS].filters =              decodeFiltersAsns(params);
  groups[FILTER_GROUP_COMMUNITIES].filters =       decodeFiltersCommunities(params);
  groups[FILTER_GROUP_EXT_COMMUNITIES].filters =   decodeFiltersExtCommunities(params);
  groups[FILTER_GROUP_LARGE_COMMUNITIES].filters = decodeFiltersLargeCommunities(params);

  return groups;
}

/*
 * Merge filters
 */
function _mergeFilters(a, b) {
  let groups = initializeFilterState();
  let setCmp = [];
  setCmp[FILTER_GROUP_SOURCES] = cmpFilterValue;
  setCmp[FILTER_GROUP_ASNS] = cmpFilterValue;
  setCmp[FILTER_GROUP_COMMUNITIES] = cmpFilterCommunity;
  setCmp[FILTER_GROUP_EXT_COMMUNITIES] = cmpFilterCommunity;
  setCmp[FILTER_GROUP_LARGE_COMMUNITIES] = cmpFilterCommunity;

  for (const i in groups) {
    groups[i].filters = mergeFilterSet(setCmp[i], a[i].filters, b[i].filters);
  }

  return groups;
}

export function mergeFilters(a, ...other) {
  let result = cloneFilters(a);
  for (const filters of other) {
    result = _mergeFilters(result, cloneFilters(filters));
  }
  return result;
}

/*
 * Merge list of filters
 */
function mergeFilterSet(inSet, a, b) {
  let result = a;
  for (const f of b) {
    const present = inSet(result, f);
    if (present) {
      // Update filter cardinality
      // present.cardinality = Math.max(f.cardinality, present.cardinality);
      present.cardinality += f.cardinality;
      continue;
    }
    result.push(f);
  }
  return result;
}

/*
 * Does a single group have any filters?
 */
export function groupHasFilters(group) {
  return group.filters.length > 0;
}

/*
 * Filters set compare
 */
function cmpFilterValue(set, filter) {
  for (const f of set) {
    if(f.value == filter.value) {
      return f;
    }
  }
  return null;
}

function cmpFilterCommunity(set, filter) {
  for (const f of set) {
    let match = true;
    for (const i in f.value) {
      if (f.value[i] != filter.value[i]) {
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
 * Do we have filters in general?
 */
export function hasFilters(groups) {
  for (const g of groups) {
    if (groupHasFilters(g)) {
      return true;
    }
  }
  return false;
}


