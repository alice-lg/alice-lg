package openbgpd

import (
	"fmt"
	"strings"
)

// Config is a OpenBGPD source config
type Config struct {
	ID   string
	Name string

	API string `ini:"api"`
}

// APIURL creates an url from the config
func (cfg *Config) APIURL(path string, params ...interface{}) string {
	u := cfg.API
	if strings.HasSuffix(u, "/") {
		u = u[:len(u)-1]
	}
	u += fmt.Sprintf(path, params...)
	return u
}
