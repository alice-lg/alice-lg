package gobgp

import "github.com/alice-lg/alice-lg/pkg/api"

// Config is a GoBGP source config
type Config struct {
	ID              string
	Name            string
	HiddenNeighbors *api.NeighborFilter `ini:"-"`

	Host     string `ini:"host"`
	Insecure bool   `ini:"insecure"`
	// ProcessingTimeout is a timeout in seconds configured per gRPC call to a given GoBGP daemon
	ProcessingTimeout int    `ini:"processing_timeout"`
	TLSCert           string `ini:"tls_crt"`
	TLSCommonName     string `ini:"tls_common_name"`
}
