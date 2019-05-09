package gobgp;

type Config struct {
	Id   			string
	Name 			string

	Host            string `ini:"host"`
	Insecure 		bool `ini:"insecure"`
	TLSCert			string `ini:"tls_crt"`
	TLSCommonName	string `ini:"tls_common_name"`
}