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

func PeerHash(peer *api.Peer) string {
	h := sha1.New()
	io.WriteString(h, string(peer.State.PeerAs))
	io.WriteString(h, peer.State.NeighborAddress)
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum[0:5])
}
