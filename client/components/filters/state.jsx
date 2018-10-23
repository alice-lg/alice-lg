
import {FILTER_GROUP_SOURCES,
        FILTER_GROUP_ASNS,
        FILTER_GROUP_COMMUNITIES,
        FILTER_GROUP_EXT_COMMUNITIES,
        FILTER_GROUP_LARGE_COMMUNITIES}
  from './groups'

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

