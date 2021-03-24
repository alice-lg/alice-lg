package openbgpd

// Config is a OpenBGPD source config
type Config struct {
	ID   string
	Name string

	API string `ini:"api"`
}
