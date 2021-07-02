package backend

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"

	"github.com/alice-lg/alice-lg/pkg/sources"
	"github.com/alice-lg/alice-lg/pkg/sources/birdwatcher"
	"github.com/alice-lg/alice-lg/pkg/sources/gobgp"
	"github.com/alice-lg/alice-lg/pkg/sources/openbgpd"
)

var (
	// ErrSourceTypeUnknown will be used if the type could
	// not be identified from the section.
	ErrSourceTypeUnknown = errors.New("source type unknown")
)

const (
	// SourceTypeBird is used for either bird 1x and 2x
	// based route servers with a birdwatcher backend.
	SourceTypeBird = "bird"

	// SourceTypeGoBGP indicates a GoBGP based source.
	SourceTypeGoBGP = "gobgp"

	// SourceTypeOpenBGPD is used for an OpenBGPD source.
	SourceTypeOpenBGPD = "openbgpd"
)

// A ServerConfig holds the runtime configuration
// for the backend.
type ServerConfig struct {
	Listen                         string `ini:"listen_http"`
	EnablePrefixLookup             bool   `ini:"enable_prefix_lookup"`
	NeighboursStoreRefreshInterval int    `ini:"neighbours_store_refresh_interval"`
	RoutesStoreRefreshInterval     int    `ini:"routes_store_refresh_interval"`
	Asn                            int    `ini:"asn"`
	EnableNeighborsStatusRefresh   bool   `ini:"enable_neighbors_status_refresh"`
}

// HousekeepingConfig describes the housekeeping interval
// and flags.
type HousekeepingConfig struct {
	Interval           int  `ini:"interval"`
	ForceReleaseMemory bool `ini:"force_release_memory"`
}

// RejectionsConfig holds rejection reasons
// associated with BGP communities
type RejectionsConfig struct {
	Reasons BgpCommunities
}

// NoexportsConfig holds no-export reasons
// associated with BGP communities and behaviour
// tweaks.
type NoexportsConfig struct {
	Reasons      BgpCommunities
	LoadOnDemand bool `ini:"load_on_demand"`
}

// RejectCandidatesConfig holds reasons for rejection
// candidates (e.g. routes that will be dropped if
// a hard filtering would be applied.)
type RejectCandidatesConfig struct {
	Communities BgpCommunities
}

// RpkiConfig defines BGP communities describing the RPKI
// validation state.
type RpkiConfig struct {
	// Define communities
	Enabled    bool     `ini:"enabled"`
	Valid      []string `ini:"valid"`
	Unknown    []string `ini:"unknown"`
	NotChecked []string `ini:"not_checked"`
	Invalid    []string `ini:"invalid"`
}

// UIConfig holds runtime settings for the web client
type UIConfig struct {
	RoutesColumns      map[string]string
	RoutesColumnsOrder []string

	NeighboursColumns      map[string]string
	NeighboursColumnsOrder []string

	LookupColumns      map[string]string
	LookupColumnsOrder []string

	RoutesRejections       RejectionsConfig
	RoutesNoexports        NoexportsConfig
	RoutesRejectCandidates RejectCandidatesConfig

	BgpCommunities BgpCommunities
	Rpki           RpkiConfig

	Theme ThemeConfig

	Pagination PaginationConfig
}

// ThemeConfig describes a theme configuration
type ThemeConfig struct {
	Path     string `ini:"path"`
	BasePath string `ini:"url_base"` // Optional, default: /theme
}

// PaginationConfig holds settings for route pagination
type PaginationConfig struct {
	RoutesFilteredPageSize    int `ini:"routes_filtered_page_size"`
	RoutesAcceptedPageSize    int `ini:"routes_accepted_page_size"`
	RoutesNotExportedPageSize int `ini:"routes_not_exported_page_size"`
}

// A SourceConfig is a generic source configuration
type SourceConfig struct {
	ID    string
	Order int
	Name  string
	Group string

	// Blackhole IPs
	Blackholes []string

	// Source configurations
	Type        string
	Birdwatcher birdwatcher.Config
	GoBGP       gobgp.Config
	OpenBGPd    openbgpd.Config

	// Source instance
	instance sources.Source
}

// Config is the application configuration
type Config struct {
	Server       ServerConfig
	Housekeeping HousekeepingConfig
	UI           UIConfig
	Sources      []*SourceConfig
	File         string
}

// SourceByID returns a source from the config by id
func (cfg *Config) SourceByID(id string) *SourceConfig {
	for _, sourceConfig := range cfg.Sources {
		if sourceConfig.ID == id {
			return sourceConfig
		}
	}
	return nil
}

// SourceInstanceByID returns an instance by id
func (cfg *Config) SourceInstanceByID(id string) sources.Source {
	sourceConfig := cfg.SourceByID(id)
	if sourceConfig == nil {
		return nil // Nothing to do here.
	}

	// Get instance from config
	return sourceConfig.getInstance()
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
func getBackendType(section *ini.Section) (string, error) {
	name := section.Name()
	if strings.HasSuffix(name, "birdwatcher") {
		return SourceTypeBird, nil
	} else if strings.HasSuffix(name, "gobgp") {
		return SourceTypeGoBGP, nil
	} else if strings.HasSuffix(name, "openbgpd") {
		return SourceTypeOpenBGPD, nil
	}

	return "", ErrSourceTypeUnknown
}

// Get UI config: Routes Columns Default
func getRoutesColumnsDefaults() (map[string]string, []string, error) {
	columns := map[string]string{
		"network":     "Network",
		"bgp.as_path": "AS Path",
		"gateway":     "Gateway",
		"interface":   "Interface",
	}
	order := []string{"network", "bgp.as_path", "gateway", "interface"}
	return columns, order, nil
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
		return getRoutesColumnsDefaults()
	}

	for _, key := range keys {
		columns[key.Name()] = section.Key(key.Name()).MustString("")
		order = append(order, key.Name())
	}

	return columns, order, nil
}

// Get UI config: Get Neighbours Columns Defaults
func getNeighboursColumnsDefaults() (map[string]string, []string, error) {
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

	return columns, order, nil
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
		return getNeighboursColumnsDefaults()
	}

	for _, key := range keys {
		columns[key.Name()] = section.Key(key.Name()).MustString("")
		order = append(order, key.Name())
	}

	return columns, order, nil
}

// Get UI config: Get Prefix search / Routes lookup columns
// As these differ slightly from our routes in the response
// (e.g. the neighbor and source rs is referenced as a nested object)
// we provide an additional configuration for this
func getLookupColumnsDefaults() (map[string]string, []string, error) {
	columns := map[string]string{
		"network":               "Network",
		"gateway":               "Gateway",
		"neighbour.asn":         "ASN",
		"neighbour.description": "Neighbor",
		"bgp.as_path":           "AS Path",
		"routeserver.name":      "RS",
	}

	order := []string{
		"network",
		"gateway",
		"bgp.as_path",
		"neighbour.asn",
		"neighbour.description",
		"routeserver.name",
	}

	return columns, order, nil
}

func getLookupColumns(config *ini.File) (
	map[string]string,
	[]string,
	error,
) {
	columns := make(map[string]string)
	order := []string{}

	section := config.Section("lookup_columns")
	keys := section.Keys()

	if len(keys) == 0 {
		return getLookupColumnsDefaults()
	}

	for _, key := range keys {
		columns[key.Name()] = section.Key(key.Name()).MustString("")
		order = append(order, key.Name())
	}

	return columns, order, nil
}

// Helper parse communities from a section body
func parseAndMergeCommunities(
	communities BgpCommunities, body string,
) BgpCommunities {

	// Parse and merge communities
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			log.Println("Skipping malformed BGP community:", line)
			continue
		}

		community := strings.TrimSpace(kv[0])
		label := strings.TrimSpace(kv[1])
		communities.Set(community, label)
	}

	return communities
}

// Get UI config: Bgp Communities
func getBgpCommunities(config *ini.File) BgpCommunities {
	// Load defaults
	communities := MakeWellKnownBgpCommunities()
	communitiesConfig := config.Section("bgp_communities")
	if communitiesConfig == nil {
		return communities // nothing else to do here, go with the default
	}

	return parseAndMergeCommunities(communities, communitiesConfig.Body())
}

// Get UI config: Get rejections
func getRoutesRejections(config *ini.File) (RejectionsConfig, error) {
	reasonsConfig := config.Section("rejection_reasons")
	if reasonsConfig == nil {
		return RejectionsConfig{}, nil
	}

	reasons := parseAndMergeCommunities(
		make(BgpCommunities),
		reasonsConfig.Body())

	rejectionsConfig := RejectionsConfig{
		Reasons: reasons,
	}

	return rejectionsConfig, nil
}

// Get UI config: Get no export config
func getRoutesNoexports(config *ini.File) (NoexportsConfig, error) {
	baseConfig := config.Section("noexport")
	reasonsConfig := config.Section("noexport_reasons")

	// Map base configuration
	noexportsConfig := NoexportsConfig{}
	baseConfig.MapTo(&noexportsConfig)

	reasons := parseAndMergeCommunities(
		make(BgpCommunities),
		reasonsConfig.Body())

	noexportsConfig.Reasons = reasons

	return noexportsConfig, nil
}

// Get UI config: Reject candidates
func getRejectCandidatesConfig(config *ini.File) (RejectCandidatesConfig, error) {
	candidateCommunities := config.Section(
		"rejection_candidates").Key("communities").String()

	if candidateCommunities == "" {
		return RejectCandidatesConfig{}, nil
	}

	communities := BgpCommunities{}
	for i, c := range strings.Split(candidateCommunities, ",") {
		communities.Set(c, fmt.Sprintf("reject-candidate-%d", i+1))
	}

	conf := RejectCandidatesConfig{
		Communities: communities,
	}

	return conf, nil
}

// Get UI config: RPKI configuration
func getRpkiConfig(config *ini.File) (RpkiConfig, error) {
	var rpki RpkiConfig
	// Defaults taken from:
	//   https://www.euro-ix.net/en/forixps/large-bgp-communities/
	section := config.Section("rpki")
	section.MapTo(&rpki)

	fallbackAsn, err := getOwnASN(config)
	if err != nil {
		log.Println(
			"Own ASN is not configured.",
			"This might lead to unexpected behaviour with BGP large communities",
		)
	}
	ownAsn := fmt.Sprintf("%d", fallbackAsn)

	// Fill in defaults or postprocess config value
	if len(rpki.Valid) == 0 {
		rpki.Valid = []string{ownAsn, "1000", "1"}
	} else {
		rpki.Valid = strings.SplitN(rpki.Valid[0], ":", 3)
	}

	if len(rpki.Unknown) == 0 {
		rpki.Unknown = []string{ownAsn, "1000", "2"}
	} else {
		rpki.Unknown = strings.SplitN(rpki.Unknown[0], ":", 3)
	}

	if len(rpki.NotChecked) == 0 {
		rpki.NotChecked = []string{ownAsn, "1000", "3"}
	} else {
		rpki.NotChecked = strings.SplitN(rpki.NotChecked[0], ":", 3)
	}

	// As the euro-ix document states, this can be a range.
	if len(rpki.Invalid) == 0 {
		rpki.Invalid = []string{ownAsn, "1000", "4", "*"}
	} else {
		// Preprocess
		rpki.Invalid = strings.SplitN(rpki.Invalid[0], ":", 3)
		tokens := []string{}
		if len(rpki.Invalid) != 3 {
			// This is wrong, we should have three parts (RS):1000:[range]
			return rpki, fmt.Errorf(
				"unexpected rpki.Invalid configuration: %v", rpki.Invalid)
		}
		tokens = strings.Split(rpki.Invalid[2], "-")
		rpki.Invalid = append([]string{rpki.Invalid[0], rpki.Invalid[1]}, tokens...)
	}

	return rpki, nil
}

// Helper: Get own ASN from ini
// This is now easy, since we enforce an ASN in
// the [server] section.
func getOwnASN(config *ini.File) (int, error) {
	server := config.Section("server")
	asn := server.Key("asn").MustInt(-1)

	if asn == -1 {
		return 0, fmt.Errorf("Could not get own ASN from config")
	}

	return asn, nil
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

// Get UI config: Pagination settings
func getPaginationConfig(config *ini.File) PaginationConfig {
	baseConfig := config.Section("pagination")

	paginationConfig := PaginationConfig{}
	baseConfig.MapTo(&paginationConfig)

	return paginationConfig
}

// Get the UI configuration from the config file
func getUIConfig(config *ini.File) (UIConfig, error) {
	uiConfig := UIConfig{}

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

	// Lookup table columns
	lookupColumns, lookupColumnsOrder, err := getLookupColumns(config)
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

	// Get reject candidates
	rejectCandidates, _ := getRejectCandidatesConfig(config)

	// RPKI filter config
	rpki, err := getRpkiConfig(config)
	if err != nil {
		return uiConfig, err
	}

	// Theme configuration: Theming is optional, if no settings
	// are found, it will be ignored
	themeConfig := getThemeConfig(config)

	// Pagination
	paginationConfig := getPaginationConfig(config)

	// Make config
	uiConfig = UIConfig{
		RoutesColumns:      routesColumns,
		RoutesColumnsOrder: routesColumnsOrder,

		NeighboursColumns:      neighboursColumns,
		NeighboursColumnsOrder: neighboursColumnsOrder,

		LookupColumns:      lookupColumns,
		LookupColumnsOrder: lookupColumnsOrder,

		RoutesRejections:       rejections,
		RoutesNoexports:        noexports,
		RoutesRejectCandidates: rejectCandidates,

		BgpCommunities: getBgpCommunities(config),
		Rpki:           rpki,

		Theme: themeConfig,

		Pagination: paginationConfig,
	}

	return uiConfig, nil
}

func getSources(config *ini.File) ([]*SourceConfig, error) {
	sources := []*SourceConfig{}

	order := 0
	sourceSections := config.ChildSections("source")
	for _, section := range sourceSections {
		if !isSourceBase(section) {
			continue
		}

		// Derive source-id from name
		sourceID := section.Name()[len("source:"):]

		// Try to get child configs and determine
		// Source type
		sourceConfigSections := section.ChildSections()
		if len(sourceConfigSections) == 0 {
			// This source has no configured backend
			return nil, fmt.Errorf("%s has no backend configuration", section.Name())
		}

		if len(sourceConfigSections) > 1 {
			// The source is ambiguous
			return nil, fmt.Errorf("%s has ambigous backends", section.Name())
		}

		// Configure backend
		backendConfig := sourceConfigSections[0]
		backendType, err := getBackendType(backendConfig)
		if err != nil {
			return nil, fmt.Errorf("%s has an unsupported backend", section.Name())
		}

		// Make config
		sourceName := section.Key("name").MustString("Unknown Source")
		sourceGroup := section.Key("group").MustString("")
		sourceBlackholes := TrimmedStringList(
			section.Key("blackholes").MustString(""))

		config := &SourceConfig{
			ID:         sourceID,
			Order:      order,
			Name:       sourceName,
			Group:      sourceGroup,
			Blackholes: sourceBlackholes,
			Type:       backendType,
		}

		// Set backend
		switch backendType {
		case SourceTypeBird:
			sourceType := backendConfig.Key("type").MustString("")
			mainTable := backendConfig.Key("main_table").MustString("master")
			peerTablePrefix := backendConfig.Key("peer_table_prefix").MustString("T")
			pipeProtocolPrefix := backendConfig.Key("pipe_protocol_prefix").MustString("M")

			if sourceType != "single_table" &&
				sourceType != "multi_table" {
				log.Fatal("Configuration error (birdwatcher source) unknown birdwatcher type:", sourceType)
			}

			log.Println("Adding birdwatcher source of type", sourceType,
				"with peer_table_prefix", peerTablePrefix,
				"and pipe_protocol_prefix", pipeProtocolPrefix)

			c := birdwatcher.Config{
				ID:   config.ID,
				Name: config.Name,

				Timezone:        "UTC",
				ServerTime:      "2006-01-02T15:04:05.999999999Z07:00",
				ServerTimeShort: "2006-01-02",
				ServerTimeExt:   "Mon, 02 Jan 2006 15:04:05 -0700",

				Type:               sourceType,
				MainTable:          mainTable,
				PeerTablePrefix:    peerTablePrefix,
				PipeProtocolPrefix: pipeProtocolPrefix,
			}

			backendConfig.MapTo(&c)
			config.Birdwatcher = c

		case SourceTypeGoBGP:
			c := gobgp.Config{
				Id:   config.ID,
				Name: config.Name,
			}

			backendConfig.MapTo(&c)
			// Update defaults:
			//  - processing_timeout
			if c.ProcessingTimeout == 0 {
				c.ProcessingTimeout = 300
			}

			config.GoBGP = c

		case SourceTypeOpenBGPD:
			c := openbgpd.Config{
				ID:   config.ID,
				Name: config.Name,
			}
			backendConfig.MapTo(&c)
			config.OpenBGPd = c
		}

		// Add to list of sources
		sources = append(sources, config)
		order++
	}

	return sources, nil
}

// Try to load configfiles as specified in the files
// list. For example:
//
//    ./etc/alice-lg/alice.conf
//    /etc/alice-lg/alice.conf
//    ./etc/alice-lg/alice.local.conf
//
func loadConfig(file string) (*Config, error) {

	// Try to get config file, fallback to alternatives
	file, err := getConfigFile(file)
	if err != nil {
		return nil, err
	}

	// Load configuration, but handle bgp communities section
	// with our own parser
	parsedConfig, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{
			"bgp_communities",
			"rejection_reasons",
			"noexport_reasons",
		},
	}, file)
	if err != nil {
		return nil, err
	}

	// Map sections
	server := ServerConfig{}
	parsedConfig.Section("server").MapTo(&server)

	housekeeping := HousekeepingConfig{}
	parsedConfig.Section("housekeeping").MapTo(&housekeeping)

	// Get all sources
	sources, err := getSources(parsedConfig)
	if err != nil {
		return nil, err
	}

	// Get UI configurations
	ui, err := getUIConfig(parsedConfig)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server:       server,
		Housekeeping: housekeeping,
		UI:           ui,
		Sources:      sources,
		File:         file,
	}

	return config, nil
}

// Get source instance from config
func (cfg *SourceConfig) getInstance() sources.Source {
	if cfg.instance != nil {
		return cfg.instance
	}

	var instance sources.Source
	switch cfg.Type {
	case SourceTypeBird:
		instance = birdwatcher.NewBirdwatcher(cfg.Birdwatcher)
	case SourceTypeGoBGP:
		instance = gobgp.NewGoBGP(cfg.GoBGP)
	case SourceTypeOpenBGPD:
		instance = openbgpd.NewSource(&cfg.OpenBGPd)
	}

	cfg.instance = instance
	return instance
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
