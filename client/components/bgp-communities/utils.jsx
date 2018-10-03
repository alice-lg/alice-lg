
/*
 * Communities helper
 */

/* 
 * Check if a community exists in a given set of communities.
 * Communities are represented as a nested object:
 *     {
 *         1234: {
 *            23: "community-leaf",
 *            42: { 
 *              1: "large-community-leaf"
 *            }
 *     }
 */

export function lookupCommunity(communities, community) {
  let lookup = communities;
  for (let c of community) {
    if (typeof(lookup) !== "object") {
      return null;
    }
    let res = lookup[c];
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
 * Reject candidate helpers:
 * 
 *  - check if prefix is a reject candidate
 *  - make css classes
 */

export function isRejectCandidate(route, rejectCommunities) {
  // Check if any reject candidate community is set
  const communities = props.route.bgp.communities;
  const largeCommunities = props.route.bgp.large_communities;

}


