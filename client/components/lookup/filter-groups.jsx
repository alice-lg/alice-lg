
export const FILTER_KEY_SOURCES = "sources"
export const FILTER_KEY_ASNS = "asns"
export const FILTER_KEY_COMMUNITIES = "communities"
export const FILTER_KEY_EXT_COMMUNITIES = "ext_communities"
export const FILTER_KEY_LARGE_COMMUNITIES = "large_communities"

export const FILTER_GROUP_SOURCES = 0
export const FILTER_GROUP_ASNS = 1
export const FILTER_GROUP_COMMUNITIES = 2
export const FILTER_GROUP_EXT_COMMUNITIES = 3
export const FILTER_GROUP_LARGE_COMMUNITIES = 4


export function filtersEqual(a, b) {
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

