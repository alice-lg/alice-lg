
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

export const initialFilterState = [
  {"key": "sources", "filters": []},
  {"key": "asns", "filters": []},
  {"key": "communities", "filters": []},
  {"key": "ext_communities", "filters": []},
  {"key": "large_communities", "filters": []},
];

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
  let groups = cloneFilters(initialFilterState);

  groups[FILTER_GROUP_SOURCES].filters =           decodeFiltersSources(params);
  groups[FILTER_GROUP_ASNS].filters =              decodeFiltersAsns(params);
  groups[FILTER_GROUP_COMMUNITIES].filters =       decodeFiltersCommunities(params);
  groups[FILTER_GROUP_EXT_COMMUNITIES].filters =   decodeFiltersExtCommunities(params);
  groups[FILTER_GROUP_LARGE_COMMUNITIES].filters = decodeFiltersLargeCommunities(params);

  return groups;
}
