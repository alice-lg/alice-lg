package openbgpd

import (
	"fmt"
	"strings"
	"time"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Config is a OpenBGPD source config
type Config struct {
	ID   string
	Name string

	CacheTTL        time.Duration
	RoutesCacheSize int

	API string `ini:"api"`

	RejectCommunities api.Communities
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
