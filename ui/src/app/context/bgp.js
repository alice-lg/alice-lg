
import { useMemo
       , useCallback
       }
  from 'react';

import { useConfig }
  from 'app/context/config';

/**
 * Create string representation from community
 */
// const toRepr = (community) => community.join(":");

/*
 * Expand variables in string:
 *    "Test AS$0 rejects $2"
 * will expand with [23, 42, 123] to
 *    "Test AS23 rejects 123"
 */
const expandVars = (str, vars) => {
  if (!str) {
    return str; // We don't have to do anything.
  }
  let res = str;
  vars.map((v, i) => {
    res = res.replace(`$${i}`, v); 
    return 0;
  });
  return res;
}

/*
 * Resolve a community description from the above described
 * tree structure.
 */
const resolveCommunity = (base, community) => {
  let lookup = base;
  for (const part of community) {
    if (typeof(lookup) !== "object") {
      return null;
    }
    let res = lookup[part];
    if (!res) {
      // Try the wildcard
      if (lookup["*"]) {
        res = lookup["*"]
      } else {
        return null; // We did everything we could
      }
    }
    lookup = res;
  }
  return lookup;
}

/*
 * Resolve all communities
 */
const resolveCommunities = (base, communities) => {
  let results = [];
  if (!communities) {
    return results;
  }
  for (const c of communities) {
    const description = resolveCommunity(base, c);
    if (description != null) {
      results.push([c, description]);
    }
  }
  return results;
}

export const useResolvedCommunities = (base, communities) =>
  useMemo(() =>
    resolveCommunities(base, communities),
    [base, communities]);


export const useResolvedCommunity = (base, community) =>
  useMemo(() =>
    resolveCommunity(base, community),
    [base, community]);

/**
 * Reject candidate:
 *  - check if prefix is a reject candidate
 */
export const useRejectCandidate = (route) => {
  const { reject_candidates } = useConfig();
  const rejectCommunities = reject_candidates?.communities;
  const largeCommunities = route?.bgp?.large_communities;
  const resolved = useResolvedCommunities(
    rejectCommunities, largeCommunities);
  return (resolved.length > 0);
};

const getReadableCommunity = (communities, community) => {
  const label = resolveCommunity(communities, community);
  return expandVars(label, community);
}

export const useReadableCommunities = () => {
  const { bgp_communities } = useConfig();
  return useCallback((community) => getReadableCommunity(
    bgp_communities,
    community,
  ), [bgp_communities]);
}

export const useReadableCommunity = (community) => {
  const getLabel = useReadableCommunities();
  return useMemo(() => getLabel(community), [
    community, getLabel,
  ]);
}

/**
 * Get blackhole communities from config
 */
export const useBlackholeCommunities = () => {
  let config = useConfig();
  return config.bgp_blackhole_communities;
}

/**
 * Match community ranges. When doing so, make sure
 * you compare the right communities. e.g. comparing
 * a large and an extended community will lead to unexpected
 * results.
 */
export const matchCommunityRange = (community, range) => {
  if (community.length !== range.length) {
    return false; 
  }
  
  for (let i in community) {
    let c = community[i];
    let rs = range[i][0];
    let re = range[i][1];
    if ((c < rs) || (c > re)) {
      return false;
    }
  }
  return true;
}

