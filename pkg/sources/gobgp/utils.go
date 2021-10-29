package gobgp

import (
	// Standard imports
	"crypto/sha1"
	"fmt"
	"io"

	// External imports
	api "github.com/osrg/gobgp/api"
	// Internal imports
)

// PeerHash calculates a peer hash
func PeerHash(peer *api.Peer) string {
	return PeerHashWithASAndAddress(peer.State.PeerAs, peer.State.NeighborAddress)
}

// PeerHashWithASAndAddress creates a peer hash (sha1) from
// the ASN and the address.
func PeerHashWithASAndAddress(asn uint32, address string) string {
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%v", asn))
	io.WriteString(h, address)
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum[0:5])
}
