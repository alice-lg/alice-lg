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
	AliceConfig, err := loadConfig(*configFilenameFlag)
	return err
}

// InitStores initializes the routes and neighbors cache
func InitStores() {
	AliceNeighboursStore = NewNeighboursStore(AliceConfig)
	AliceRoutesStore = NewRoutesStore(AliceConfig)
}
