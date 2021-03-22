package backend

import (
	"fmt"
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

// BgpCommunities is a tree representation of BGP communities
type BgpCommunities map[string]interface{}

// MakeWellKnownBgpCommunities returns a BgpCommunities
// map with well known communities.
func MakeWellKnownBgpCommunities() BgpCommunities {
	c := BgpCommunities{
		"65535": BgpCommunities{
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
func (c BgpCommunities) Lookup(community string) (string, error) {
	path := strings.Split(community, ":")
	var lookup interface{} // This is all much too dynamic...
	lookup = c

	for _, key := range path {
		key = strings.TrimSpace(key)

		clookup, ok := lookup.(BgpCommunities)
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
func (c BgpCommunities) Set(community string, label string) {
	path := strings.Split(community, ":")
	var lookup interface{} // Again, this is all much too dynamic...
	lookup = c

	for _, key := range path[:len(path)-1] {
		key = strings.TrimSpace(key)
		clookup, ok := lookup.(BgpCommunities)
		if !ok {
			break
		}

		res, ok := clookup[key]
		if !ok {
			// The key does not exist, create it!
			clookup[key] = BgpCommunities{}
			res = clookup[key]
		}

		lookup = res
	}

	slookup := lookup.(BgpCommunities)
	slookup[path[len(path)-1]] = label
}
