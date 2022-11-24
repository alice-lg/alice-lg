package api

import (
	"fmt"
	"strconv"
	"strings"
)

/*
Implement BGP Communities Lookup Base

We initialize the dictionary with well known communities and
store the representation as a string with : as delimiter.

From: https://www.iana.org/assignments/bgp-well-known-communities/bgp-well-known-communities.xhtml

    0x00000000-0x0000FFFF   Reserved    [RFC1997]
    0x00010000-0xFFFEFFFF   Reserved for Private Use    [RFC1997]

    0xFFFF0000  GRACEFUL_SHUTDOWN   [RFC8326]
    0xFFFF0001  ACCEPT_OWN          [RFC7611]
    0xFFFF0002  ROUTE_FILTER_TRANSLATED_v4  [draft-l3vpn-legacy-rtc]
    0xFFFF0003  ROUTE_FILTER_v4     [draft-l3vpn-legacy-rtc]
    0xFFFF0004  ROUTE_FILTER_TRANSLATED_v6  [draft-l3vpn-legacy-rtc]
    0xFFFF0005  ROUTE_FILTER_v6     [draft-l3vpn-legacy-rtc]
    0xFFFF0006  LLGR_STALE          [draft-uttaro-idr-bgp-persistence]
    0xFFFF0007  NO_LLGR             [draft-uttaro-idr-bgp-persistence]
    0xFFFF0008  accept-own-nexthop  [draft-agrewal-idr-accept-own-nexthop]

    0xFFFF0009-0xFFFF0299   Unassigned

    0xFFFF029A  BLACKHOLE           [RFC7999]

    0xFFFF029B-0xFFFFFF00   Unassigned

    0xFFFFFF01  NO_EXPORT           [RFC1997]
    0xFFFFFF02  NO_ADVERTISE        [RFC1997]
    0xFFFFFF03  NO_EXPORT_SUBCONFED [RFC1997]
    0xFFFFFF04  NOPEER              [RFC3765]
    0xFFFFFF05-0xFFFFFFFF   Unassigned
*/

// BGPCommunityMap is a tree representation of BGP communities
// where the leaf is a description or reason.
type BGPCommunityMap map[string]interface{}

// MakeWellKnownBGPCommunities returns a BGPCommunityMap
// map with well known communities.
func MakeWellKnownBGPCommunities() BGPCommunityMap {
	c := BGPCommunityMap{
		"65535": BGPCommunityMap{
			"0": "graceful shutdown",
			"1": "accept own",
			"2": "route filter translated v4",
			"3": "route filter v4",
			"4": "route filter translated v6",
			"5": "route filter v6",
			"6": "llgr stale",
			"7": "no llgr",
			"8": "accept-own-nexthop",

			"666": "blackhole",

			"1048321": "no export",
			"1048322": "no advertise",
			"1048323": "no export subconfed",
			"1048324": "nopeer",
		},
	}

	return c
}

// Lookup searches for a label in the communities map
func (c BGPCommunityMap) Lookup(community string) (string, error) {
	path := strings.Split(community, ":")
	var lookup interface{} // This is all much too dynamic...
	lookup = c

	for _, key := range path {
		key = strings.TrimSpace(key)

		clookup, ok := lookup.(BGPCommunityMap)
		if !ok {
			// This happens if path.len > depth
			return "", fmt.Errorf("community not found @ %v", key)
		}

		res, ok := clookup[key]
		if !ok {
			// Try to fall back to wildcard key
			res, ok = clookup["*"]
			if !ok {
				break // we did everything we could.
			}
		}

		lookup = res
	}

	label, ok := lookup.(string)
	if !ok {
		return "", fmt.Errorf("community not found: %v", community)
	}

	return label, nil
}

// Set assignes a label to a community
func (c BGPCommunityMap) Set(community string, label string) {
	path := strings.Split(community, ":")
	var lookup interface{} // Again, this is all much too dynamic...
	lookup = c

	for _, key := range path[:len(path)-1] {
		key = strings.TrimSpace(key)
		clookup, ok := lookup.(BGPCommunityMap)
		if !ok {
			break
		}

		res, ok := clookup[key]
		if !ok {
			// The key does not exist, create it!
			clookup[key] = BGPCommunityMap{}
			res = clookup[key]
		}

		lookup = res
	}

	slookup := lookup.(BGPCommunityMap)
	slookup[path[len(path)-1]] = label
}

// Communities enumerates all bgp communities into
// a set of api.Communities.
// CAVEAT: Wildcards are substituted by 0 and ** ARE NOT ** expanded.
func (c BGPCommunityMap) Communities() Communities {
	communities := Communities{}
	// We could do this recursive, or assume that
	// the max depth is 3.
	for uVal, c1 := range c {
		u, err := strconv.Atoi(uVal)
		if err != nil {
			u = 0
		}
		for vVal, c2 := range c1.(BGPCommunityMap) {
			v, err := strconv.Atoi(vVal)
			if err != nil {
				v = 0
			}

			com2, ok := c2.(BGPCommunityMap)
			if !ok {
				// we only have labels here
				communities = append(
					communities, Community{u, v})
				continue
			}

			for wVal := range com2 {
				w, err := strconv.Atoi(wVal)
				if err != nil {
					w = 0
				}

				communities = append(
					communities, Community{u, v, w})
				continue
			}
		}
	}
	return communities
}

// BGPCommunity types: Standard, Extended and Large
const (
	BGPCommunityTypeStd = iota
	BGPCommunityTypeExt
	BGPCommunityTypeLarge
)

// BGPCommunityRange is a list of tuples with the start and end
// of the range defining a community.
type BGPCommunityRange []interface{}

// Type classifies the BGP Ranged BGP Community into: std, large, ext
func (c BGPCommunityRange) Type() int {
	if len(c) == 2 {
		return BGPCommunityTypeStd
	}
	if _, ok := c[0].([]string); ok {
		return BGPCommunityTypeExt
	}
	return BGPCommunityTypeLarge
}

// A BGPCommunitiesSet is a set of communities, large and extended.
// The communities are described as ranges.
type BGPCommunitiesSet struct {
	Standard []BGPCommunityRange `json:"standard"`
	Extended []BGPCommunityRange `json:"extended"`
	Large    []BGPCommunityRange `json:"large"`
}
