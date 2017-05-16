package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

type ServerConfig struct {
	Listen string `toml:"listen_http"`
}

type UiConfig struct {
	ShowLastReboot bool `toml:"rs_show_last_reboot"`
}

type SourceConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
	Ui     UiConfig

	Sources []SourceConfig
}

// Try to load configfiles as specified in the files
// list. For example:
//
//    ./etc/alicelg/alice.conf
//    /etc/alicelg/alice.conf
//    ./etc/alicelg/alice.local.conf
//
func LoadConfigs(configFiles []string) (*Config, error) {
	config := &Config{}
	hasConfig := false
	var confError error

	for _, filename := range configFiles {
		tmp := &Config{}
		_, err := toml.DecodeFile(filename, tmp)
		if err != nil {
			continue
		} else {
			log.Println("Using config file:", filename)
			hasConfig = true
			// Merge configs
			if err := mergo.Merge(config, tmp); err != nil {
				return nil, err
			}
		}
	}

	if !hasConfig {
		confError = fmt.Errorf("Could not load any config file")
	}

	return config, confError
}

func ConfigOptions(filename string) []string {
	return []string{
		strings.Join([]string{"/", filename}, ""),
		strings.Join([]string{"./", filename}, ""),
		strings.Replace(filename, ".conf", ".local.conf", 1),
	}
}
