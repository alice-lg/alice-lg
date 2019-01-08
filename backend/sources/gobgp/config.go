package gobgp;

type Config struct {
	Id   			string
	Name 			string

	Host            string `ini:"api"`
	Insecure 		bool `ini:"insecure"`
}