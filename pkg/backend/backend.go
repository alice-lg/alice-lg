package backend

// Globals
var (
	AliceConfig          *Config
	AliceRoutesStore     *RoutesStore
	AliceNeighboursStore *NeighboursStore
)

// InitConfig loads the configuration into the global
// AliceConfig
func InitConfig(filename string) error {
	var err error
	AliceConfig, err = loadConfig(filename)
	return err
}

// InitStores initializes the routes and neighbors cache
func InitStores() {
	AliceNeighboursStore = NewNeighboursStore(AliceConfig)
	AliceRoutesStore = NewRoutesStore(AliceConfig)
}
