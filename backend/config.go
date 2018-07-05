package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alice-lg/alice-lg/backend/sources"
	"github.com/alice-lg/alice-lg/backend/sources/birdwatcher"

	"github.com/go-ini/ini"
	_ "github.com/imdario/mergo"
)

const SOURCE_UNKNOWN = 0
const SOURCE_BIRDWATCHER = 1

type ServerConfig struct {
	Listen             string `ini:"listen_http"`
	EnablePrefixLookup bool   `ini:"enable_prefix_lookup"`
}

type RejectionsConfig struct {
	Asn      int `ini:"asn"`
	RejectId int `ini:"reject_id"`

	Reasons map[int]string
}

type NoexportsConfig struct {
	Asn        int `ini:"asn"`
	NoexportId int `ini:"noexport_id"`

	Reasons map[int]string
}

type UiConfig struct {
	RoutesColumns      map[string]string
	RoutesColumnsOrder []string

	NeighboursColumns      map[string]string
	NeighboursColumnsOrder []string

	RoutesRejections RejectionsConfig
	RoutesNoexports  NoexportsConfig

	Theme ThemeConfig
}

type ThemeConfig struct {
	Path     string `ini:"path"`
	BasePath string `ini:"url_base"` // Optional, default: /theme
}

type SourceConfig struct {
	Id   int
	Name string
	Type int

	// Source configurations
	Birdwatcher birdwatcher.Config
}

type Config struct {
	Server  ServerConfig
	Ui      UiConfig
	Sources []SourceConfig
	File    string

	instances map[SourceConfig]sources.Source
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

// Get UI config: Routes Columns Default
func getRoutesColumnsDefault() (map[string]string, []string) {
	columns := map[string]string{
		"bgp.as_path": "AS Path",
		"gateway":     "Gateway",
		"interface":   "Interface",
	}

	order := []string{"Network", "bgp.as_path", "gateway", "interface"}

	return columns, order
}

// Get UI config: Routes Columns
// The columns displayed in the frontend.
// The columns are ordered as in the config file.
//
// In case the configuration is empty, fall back to
// the defaults as defined in getRoutesColumnsDefault()
//
func getRoutesColumns(config *ini.File) (map[string]string, []string, error) {
	columns := make(map[string]string)
	order := []string{}

	section := config.Section("routes_columns")
	keys := section.Keys()

	if len(keys) == 0 {
		defaultColumns, defaultOrder := getRoutesColumnsDefault()
		return defaultColumns, defaultOrder, nil
	}

	for _, key := range keys {
		columns[key.Name()] = section.Key(key.Name()).MustString("")
		order = append(order, key.Name())
	}

	return columns, order, nil
}

// Get UI config: Get Neighbours Columns Defaults
func getNeighboursColumnsDefaults() (map[string]string, []string) {
	columns := map[string]string{
		"address":         "Neighbour",
		"asn":             "ASN",
		"state":           "State",
		"Uptime":          "Uptime",
		"Description":     "Description",
		"routes_received": "Routes Recv.",
		"routes_filtered": "Routes Filtered",
	}

	order := []string{
		"address", "asn", "state",
		"Uptime", "Description", "routes_received", "routes_filtered",
	}

	return columns, order
}

// Get UI config: Get Neighbours Columns
// basically the same as with the routes columns.
func getNeighboursColumns(config *ini.File) (
	map[string]string,
	[]string,
	error,
) {
	columns := make(map[string]string)
	order := []string{}

	section := config.Section("neighbours_columns")
	keys := section.Keys()

	if len(keys) == 0 {
		defaultColumns, defaultOrder := getNeighboursColumnsDefaults()
		return defaultColumns, defaultOrder, nil
	}

	for _, key := range keys {
		columns[key.Name()] = section.Key(key.Name()).MustString("")
		order = append(order, key.Name())
	}

	return columns, order, nil
}

// Get UI config: Get rejections
func getRoutesRejections(config *ini.File) (RejectionsConfig, error) {
	reasons := make(map[int]string)
	baseConfig := config.Section("rejection")
	reasonsConfig := config.Section("rejection_reasons")

	// Map base configuration
	rejectionsConfig := RejectionsConfig{}
	baseConfig.MapTo(&rejectionsConfig)

	// Map reasons
	keys := reasonsConfig.Keys()
	for _, key := range keys {
		id, err := strconv.Atoi(key.Name())
		if err != nil {
			return rejectionsConfig, err
		}
		reasons[id] = reasonsConfig.Key(key.Name()).MustString("")
	}

	rejectionsConfig.Reasons = reasons

	return rejectionsConfig, nil
}

// Get UI config: Get no export config
func getRoutesNoexports(config *ini.File) (NoexportsConfig, error) {
	reasons := make(map[int]string)
	baseConfig := config.Section("noexport")
	reasonsConfig := config.Section("noexport_reasons")

	// Map base configuration
	noexportsConfig := NoexportsConfig{}
	baseConfig.MapTo(&noexportsConfig)

	// Map reasons for missing export
	keys := reasonsConfig.Keys()
	for _, key := range keys {
		id, err := strconv.Atoi(key.Name())
		if err != nil {
			return noexportsConfig, err
		}
		reasons[id] = reasonsConfig.Key(key.Name()).MustString("")
	}

	noexportsConfig.Reasons = reasons

	return noexportsConfig, nil
}

// Get UI config: Theme settings
func getThemeConfig(config *ini.File) ThemeConfig {
	baseConfig := config.Section("theme")

	themeConfig := ThemeConfig{}
	baseConfig.MapTo(&themeConfig)

	if themeConfig.BasePath == "" {
		themeConfig.BasePath = "/theme"
	}

	return themeConfig
}

// Get the UI configuration from the config file
func getUiConfig(config *ini.File) (UiConfig, error) {
	uiConfig := UiConfig{}

	// Get route columns
	routesColumns, routesColumnsOrder, err := getRoutesColumns(config)
	if err != nil {
		return uiConfig, err
	}

	// Get neighbours table columns
	neighboursColumns,
		neighboursColumnsOrder,
		err := getNeighboursColumns(config)
	if err != nil {
		return uiConfig, err
	}

	// Get rejections and reasons
	rejections, err := getRoutesRejections(config)
	if err != nil {
		return uiConfig, err
	}

	noexports, err := getRoutesNoexports(config)
	if err != nil {
		return uiConfig, err
	}

	// Theme configuration: Theming is optional, if no settings
	// are found, it will be ignored
	themeConfig := getThemeConfig(config)

	// Make config
	uiConfig = UiConfig{
		RoutesColumns:      routesColumns,
		RoutesColumnsOrder: routesColumnsOrder,

		NeighboursColumns:      neighboursColumns,
		NeighboursColumnsOrder: neighboursColumnsOrder,

		RoutesRejections: rejections,
		RoutesNoexports:  noexports,

		Theme: themeConfig,
	}

	return uiConfig, nil
}

func getSources(config *ini.File) ([]SourceConfig, error) {
	sources := []SourceConfig{}

	sourceSections := config.ChildSections("source")
	sourceId := 0
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
			Id:   sourceId,
			Name: section.Key("name").MustString("Unknown Source"),
			Type: backendType,
		}

		// Set backend
		switch backendType {
		case SOURCE_BIRDWATCHER:
			c := birdwatcher.Config{
				Id:   config.Id,
				Name: config.Name,

				Timezone:        "UTC",
				ServerTime:      "2006-01-02T15:04:05.999999999Z07:00",
				ServerTimeShort: "2006-01-02",
				ServerTimeExt:   "Mon, 02 Jan 2006 15:04:05 -0700",
			}
			backendConfig.MapTo(&c)
			config.Birdwatcher = c
		}

		// Add to list of sources
		sources = append(sources, config)

		sourceId += 1
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
func loadConfig(file string) (*Config, error) {

	// Try to get config file, fallback to alternatives
	file, err := getConfigFile(file)
	if err != nil {
		return nil, err
	}

	parsedConfig, err := ini.LooseLoad(file)
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

	// Get UI configurations
	ui, err := getUiConfig(parsedConfig)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server:  server,
		Ui:      ui,
		Sources: sources,
		File:    file,
	}

	return config, nil
}

// Get source instance from config
func (source SourceConfig) getInstance() sources.Source {
	switch source.Type {
	case SOURCE_BIRDWATCHER:
		return birdwatcher.NewBirdwatcher(source.Birdwatcher)
	}

	return nil
}

// Get configuration file with fallbacks
func getConfigFile(filename string) (string, error) {
	// Check if requested file is present
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Fall back to local filename
		filename = ".." + filename
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = strings.Replace(filename, ".conf", ".local.conf", 1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "not_found", fmt.Errorf("could not find any configuration file")
	}

	return filename, nil
}
