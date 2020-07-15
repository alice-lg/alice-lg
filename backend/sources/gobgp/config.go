package gobgp

type Config struct {
	Id   string
	Name string

	Host     string `ini:"host"`
	Insecure bool   `ini:"insecure"`
	// ProcessingTimeout is a timeout configured per gRPC call to a given GoBGP backend
	ProcessingTimeout int    `ini:"processing_timeout"`
	TLSCert           string `ini:"tls_crt"`
	TLSCommonName     string `ini:"tls_common_name"`
}
