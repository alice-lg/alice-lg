
/*
 * Communities helper
 */

/* 
 * Communities are represented as a nested object:
 *     {
 *         1234: {
 *            23: "community-leaf",
 *            42: { 
 *              1: "large-community-leaf"
 *            }
 *     }
 */

/*
 * Resolve a community description from the above described
 * tree structure.
 */
export function resolveCommunity(base, community) {
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
export function resolveCommunities(base, communities) {
  let results = [];
  for (const c of communities) {
    const description = resolveCommunity(base, c);
    if (description != null) {
      results.push([c, description]);
    }
  }
  return results;
}


/*
 * Reject candidate helpers:
 * 
 *  - check if prefix is a reject candidate
 *  - make css classes
 */

export function isRejectCandidate(rejectCommunities, route) {
  // Check if any reject candidate community is set
  const communities = route.bgp.communities;
  const largeCommunities = route.bgp.large_communities;

  const resolved = resolveCommunities(
    rejectCommunities, largeCommunities
  );

  return (resolved.length > 0);
}


