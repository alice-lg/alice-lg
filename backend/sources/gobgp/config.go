package gobgp;

type Config struct {
	Id   			string
	Name 			string

	Host            string `ini:"host"`
	Insecure 		bool `ini:"insecure"`
}