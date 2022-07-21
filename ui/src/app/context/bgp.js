
import { useMemo }
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


export const useReadableCommunity = (community) => {
  const { bgp_communities } = useConfig();
  return useMemo(() => {
    const label = resolveCommunity(bgp_communities, community);
    return expandVars(label, community);
  }, [bgp_communities, community]);
}


