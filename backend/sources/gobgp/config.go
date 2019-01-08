package gobgp;

type Config struct {
	Id   			string
	Name 			string

	Host            string `ini:"api"`
}