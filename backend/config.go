package main

import (
	"fmt"
	"strings"

	"github.com/ecix/alice-lg/backend/sources/birdwatcher"

	"github.com/go-ini/ini"
	_ "github.com/imdario/mergo"
)

const SOURCE_UNKNOWN = 0
const SOURCE_BIRDWATCHER = 1

type ServerConfig struct {
	Listen string `ini:"listen_http"`
}

type SourceConfig struct {
	Name string
	Type int

	// Source configurations
	Birdwatcher birdwatcher.Config
}

type Config struct {
	Server  ServerConfig
	Sources []SourceConfig
}

// Get sources keys form ini
func getSourcesKeys(config *ini.File) []string {
	sources := []string{}
	sections := config.SectionStrings()
	for _, section := range sections {
		if strings.HasPrefix(section, "source") {
			sources = append(sources, section)
		}
	}
	return sources
}

func isSourceBase(section *ini.Section) bool {
	return len(strings.Split(section.Name(), ".")) == 2
}

// Get backend configuration type
func getBackendType(section *ini.Section) int {
	name := section.Name()
	if strings.HasSuffix(name, "birdwatcher") {
		return SOURCE_BIRDWATCHER
	}

	return SOURCE_UNKNOWN
}

func getSources(config *ini.File) ([]SourceConfig, error) {
	sources := []SourceConfig{}

	sourceSections := config.ChildSections("source")
	for _, section := range sourceSections {
		if !isSourceBase(section) {
			continue
		}

		// Try to get child configs and determine
		// Source type
		sourceConfigSections := section.ChildSections()
		if len(sourceConfigSections) == 0 {
			// This source has no configured backend
			return sources, fmt.Errorf("%s has no backend configuration", section.Name())
		}

		if len(sourceConfigSections) > 1 {
			// The source is ambiguous
			return sources, fmt.Errorf("%s has ambigous backends", section.Name())
		}

		// Configure backend
		backendConfig := sourceConfigSections[0]
		backendType := getBackendType(backendConfig)

		if backendType == SOURCE_UNKNOWN {
			return sources, fmt.Errorf("%s has an unsupported backend", section.Name())
		}

		// Make config
		config := SourceConfig{
			Type: backendType,
		}

		// Set backend
		switch backendType {
		case SOURCE_BIRDWATCHER:
			backendConfig.MapTo(&config.Birdwatcher)
		}

		// Add to list of sources
		sources = append(sources, config)
	}

	return sources, nil
}

// Try to load configfiles as specified in the files
// list. For example:
//
//    ./etc/alicelg/alice.conf
//    /etc/alicelg/alice.conf
//    ./etc/alicelg/alice.local.conf
//
func loadConfigs(base, global, local string) (*Config, error) {
	parsedConfig, err := ini.LooseLoad(base, global, local)
	if err != nil {
		return nil, err
	}

	// Map sections
	server := ServerConfig{}
	parsedConfig.Section("server").MapTo(&server)

	// Get all sources
	sources, err := getSources(parsedConfig)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server:  server,
		Sources: sources,
	}

	fmt.Println(config)
	return config, nil
}

func configOptions(filename string) []string {
	return []string{
		strings.Join([]string{"/", filename}, ""),
		strings.Join([]string{"./", filename}, ""),
		strings.Replace(filename, ".conf", ".local.conf", 1),
	}
}
