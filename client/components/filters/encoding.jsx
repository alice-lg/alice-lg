
import {
  FILTER_GROUP_SOURCES,
  FILTER_GROUP_ASNS,
  FILTER_GROUP_COMMUNITIES,
  FILTER_GROUP_EXT_COMMUNITIES,
  FILTER_GROUP_LARGE_COMMUNITIES,
} from './groups'


function _makeFilter(value) {
  return {
    name: "",
    value: value,
    cardinality: 1,
  }
}

export function decodeFiltersSources(params) {
  if (!params.sources) {
    return []; // No params available
  }
  const sources = params.sources.split(",");
  return sources.map((sid) => _makeFilter(sid));
}

export function decodeFiltersAsns(params) {
  if (!params.asns) {
    return []; // No params available
  }
  const asns = params.asns.split(",");
  return asns.map((asn) => _makeFilter(parseInt(asn, 10)));
}

function _decodeCommunity(community) {
  const parts = community.split(":");
  return parts.map((p) => parseInt(p, 10));
}

function _decodeExtCommunity(community) {
  const parts = community.split(":");
  return [parts[0]].concat(parts.slice(1).map((p) => parseInt(p, 10)));
}

export function decodeFiltersCommunities(params) {
  if (!params.communities) {
    return []; // No params available
  }
  const communities = params.communities.split(",");
  return communities.map((c) => _makeFilter(_decodeCommunity(c)));
}

export function decodeFiltersExtCommunities(params) {
  if (!params.ext_communities) {
    return []; // No params available
  }
  const communities = params.ext_communities.split(",");
  return communities.map((c) => _makeFilter(_decodeExtCommunity(c)));
}

export function decodeFiltersLargeCommunities(params) {
  if (!params.large_communities) {
    return []; // No params available
  }
  const communities = params.large_communities.split(",");
  return communities.map((c) => _makeFilter(_decodeCommunity(c)));
}


export function encodeGroupInt(group) {
  if (!group.filters.length) {
    return "";
  }
  const values = group.filters.map((f) => f.value).join(",");
  return `&${group.key}=${values}`;
}

export function encodeGroupCommunities(group) {
  if (!group.filters.length) {
    return "";
  }
  const values = group.filters.map((f) => f.value.join(":")).join(",");
  return `&${group.key}=${values}`;
}

export function filtersUrlEncode(filters) {
  let encoded = "";

  encoded += encodeGroupInt(filters[FILTER_GROUP_SOURCES]);
  encoded += encodeGroupInt(filters[FILTER_GROUP_ASNS]);
  encoded += encodeGroupCommunities(filters[FILTER_GROUP_COMMUNITIES]);
  encoded += encodeGroupCommunities(filters[FILTER_GROUP_EXT_COMMUNITIES]);
  encoded += encodeGroupCommunities(filters[FILTER_GROUP_LARGE_COMMUNITIES]);

  return encoded;
}

