package backend

// Globals
var (
	AliceConfig          *Config
	AliceRoutesStore     *RoutesStore
	AliceNeighborsStore *NeighborsStore
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
	AliceNeighborsStore = NewNeighborsStore(AliceConfig)
	AliceRoutesStore = NewRoutesStore(AliceConfig)
}
