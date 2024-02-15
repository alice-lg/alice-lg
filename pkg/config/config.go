// Package config provides runtime configuration
// for the Alice Looking Glass.
//
// This configuration is read from a config file.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"

	"github.com/alice-lg/alice-lg/pkg/api"
	"github.com/alice-lg/alice-lg/pkg/decoders"
	"github.com/alice-lg/alice-lg/pkg/pools"
	"github.com/alice-lg/alice-lg/pkg/sources"
	"github.com/alice-lg/alice-lg/pkg/sources/birdwatcher"
	"github.com/alice-lg/alice-lg/pkg/sources/gobgp"
	"github.com/alice-lg/alice-lg/pkg/sources/openbgpd"
)

var (
	// ErrSourceTypeUnknown will be used if the type could
	// not be identified from the section.
	ErrSourceTypeUnknown = errors.New("source type unknown")

	// ErrPostgresUnconfigured will occur when the
	// postgres database URL is required, but missing.
	ErrPostgresUnconfigured = errors.New(
		"the selected postgres backend requires configuration")
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

const (
	// SourceBackendBirdwatcher is used to indicate that
	// the source is using a birdwatcher interface.
	SourceBackendBirdwatcher = "birdwatcher"

	// SourceBackendGoBGP is used when the source is consuming
	// a GoBGP daemon via grpc API.
	SourceBackendGoBGP = "gobgp"

	// SourceBackendOpenBGPDStateServer is used when the openbgpd
	// is exported using the openbgpd-state-server.
	SourceBackendOpenBGPDStateServer = "openbgpd-state-server"

	// SourceBackendOpenBGPDBgplgd is used when the openbgpd
	// state is exported through the bgplgd.
	SourceBackendOpenBGPDBgplgd = "openbgpd-bgplgd"
)

const (
	// DefaultHTTPTimeout is the time in seconds after which the
	// server will timeout.
	DefaultHTTPTimeout = 120

	// DefaultPrefixLookupCommunityFilterCutoff is the number of
	// routes after which the community filter will not be
	// available.
	DefaultPrefixLookupCommunityFilterCutoff = 100000

	// DefaultRoutesStoreQueryLimit is the default limit for
	// prefixes returned from the store.
	DefaultRoutesStoreQueryLimit = 200000
)

// A ServerConfig holds the runtime configuration
// for the backend.
type ServerConfig struct {
	Listen                            string `ini:"listen_http"`
	HTTPTimeout                       int    `ini:"http_timeout"`
	EnablePrefixLookup                bool   `ini:"enable_prefix_lookup"`
	PrefixLookupCommunityFilterCutoff int    `ini:"prefix_lookup_community_filter_cutoff"`
	NeighborsStoreRefreshInterval     int    `ini:"neighbors_store_refresh_interval"`
	NeighborsStoreRefreshParallelism  int    `ini:"neighbors_store_refresh_parallelism"`
	RoutesStoreRefreshInterval        int    `ini:"routes_store_refresh_interval"`
	RoutesStoreRefreshParallelism     int    `ini:"routes_store_refresh_parallelism"`
	RoutesStoreQueryLimit             uint   `ini:"routes_store_query_limit"`
	StoreBackend                      string `ini:"store_backend"`
	DefaultAsn                        int    `ini:"asn"`
	EnableNeighborsStatusRefresh      bool   `ini:"enable_neighbors_status_refresh"`
	StreamParserThrottle              int    `ini:"stream_parser_throttle"`
}

// PostgresConfig is the configuration for the database
// connection when the postgres backend is used.
type PostgresConfig struct {
	URL      string `ini:"url"`
	MaxConns int32  `ini:"max_connections"`
	MinConns int32  `ini:"min_connections"`
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
	Reasons api.BGPCommunityMap
}

// NoexportsConfig holds no-export reasons
// associated with BGP communities and behaviour
// tweaks.
type NoexportsConfig struct {
	Reasons      api.BGPCommunityMap
	LoadOnDemand bool `ini:"load_on_demand"`
}

// RejectCandidatesConfig holds reasons for rejection
// candidates (e.g. routes that will be dropped if
// a hard filtering would be applied.)
type RejectCandidatesConfig struct {
	Communities api.BGPCommunityMap
}

// RpkiConfig defines BGP communities describing the RPKI
// validation state.
type RpkiConfig struct {
	// Define communities
	Enabled    bool       `ini:"enabled"`
	Valid      [][]string `ini:"valid"`
	Unknown    [][]string `ini:"unknown"`
	NotChecked [][]string `ini:"not_checked"`
	Invalid    [][]string `ini:"invalid"`
}

// UIConfig holds runtime settings for the web client
type UIConfig struct {
	RoutesColumns      map[string]string
	RoutesColumnsOrder []string

	NeighborsColumns      map[string]string
	NeighborsColumnsOrder []string

	LookupColumns      map[string]string
	LookupColumnsOrder []string

	RoutesRejections       RejectionsConfig
	RoutesNoexports        NoexportsConfig
	RoutesRejectCandidates RejectCandidatesConfig

	BGPCommunities          api.BGPCommunityMap
	BGPBlackholeCommunities api.BGPCommunitiesSet
	Rpki                    RpkiConfig

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
	Backend     string
	Birdwatcher birdwatcher.Config
	GoBGP       gobgp.Config
	OpenBGPD    openbgpd.Config

	// Source instance
	instance sources.Source
}

// Config is the application configuration
type Config struct {
	Server       ServerConfig
	Postgres     *PostgresConfig
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
	return sourceConfig.GetInstance()
}

func isSourceBase(section *ini.Section) bool {
	return len(strings.Split(section.Name(), ".")) == 2
}

// Get backend configuration type
func sourceBackendTypeFromConfig(section *ini.Section) (string, error) {
	name := section.Name()
	if strings.HasSuffix(name, "birdwatcher") {
		return SourceBackendBirdwatcher, nil
	} else if strings.HasSuffix(name, "gobgp") {
		return SourceBackendGoBGP, nil
	} else if strings.HasSuffix(name, "openbgpd-bgplgd") {
		return SourceBackendOpenBGPDBgplgd, nil
	} else if strings.HasSuffix(name, "openbgpd-state-server") {
		return SourceBackendOpenBGPDStateServer, nil
	}

	return "", ErrSourceTypeUnknown
}

// sourceTypeFromBackendType will return the backend source type
// for a given backend type
func sourceTypeFromBackendType(t string) string {
	switch t {
	case SourceBackendBirdwatcher:
		return SourceTypeBird
	case SourceBackendGoBGP:
		return SourceTypeGoBGP
	case SourceBackendOpenBGPDStateServer:
		return SourceTypeOpenBGPD
	case SourceBackendOpenBGPDBgplgd:
		return SourceTypeOpenBGPD
	default:
		return ""
	}
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

// Get UI config: Get Neighbors Columns Defaults
func getNeighborsColumnsDefaults() (map[string]string, []string, error) {
	columns := map[string]string{
		"address":         "Neighbor",
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

// Get UI config: Get Neighbors Columns
// basically the same as with the routes columns.
func getNeighborsColumns(config *ini.File) (
	map[string]string,
	[]string,
	error,
) {
	columns := make(map[string]string)
	order := []string{}

	section := config.Section("neighbors_columns")
	keys := section.Keys()

	if len(keys) == 0 {
		return getNeighborsColumnsDefaults()
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
		"network":              "Network",
		"gateway":              "Gateway",
		"neighbor.asn":         "ASN",
		"neighbor.description": "Neighbor",
		"bgp.as_path":          "AS Path",
		"routeserver.name":     "RS",
	}

	order := []string{
		"network",
		"gateway",
		"bgp.as_path",
		"neighbor.asn",
		"neighbor.description",
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

// Get UI config: BGP Communities
func getBGPCommunityMap(config *ini.File) api.BGPCommunityMap {
	// Load defaults
	communities := api.MakeWellKnownBGPCommunities()
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
		make(api.BGPCommunityMap),
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
	if err := baseConfig.MapTo(&noexportsConfig); err != nil {
		return noexportsConfig, err
	}

	reasons := parseAndMergeCommunities(
		make(api.BGPCommunityMap),
		reasonsConfig.Body())

	noexportsConfig.Reasons = reasons

	return noexportsConfig, nil
}

// Get UI config: blackhole communities
func getBlackholeCommunities(config *ini.File) (api.BGPCommunitiesSet, error) {
	section := config.Section("blackhole_communities")
	defaultBlackholes := api.BGPCommunitiesSet{
		Standard: []api.BGPCommunityRange{
			{[]interface{}{65535, 65535}, []interface{}{666, 666}},
		},
	}
	if section == nil {
		return defaultBlackholes, nil
	}
	set, err := parseRangeCommunitiesSet(section.Body())
	if err != nil {
		return defaultBlackholes, err
	}
	set.Standard = append(set.Standard, defaultBlackholes.Standard...)
	return *set, nil
}

// Get UI config: Reject candidates
func getRejectCandidatesConfig(config *ini.File) (RejectCandidatesConfig, error) {
	candidateCommunities := config.Section(
		"rejection_candidates").Key("communities").String()

	if candidateCommunities == "" {
		return RejectCandidatesConfig{}, nil
	}

	communities := api.BGPCommunityMap{}
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
	lines := strings.Split(section.Body(), "\n")

	for _, line := range lines {
		l := strings.TrimSpace(line)
		if !strings.Contains(l, "=") {
			continue
		}
		parts := strings.SplitN(l, "=", 2)
		if len(parts) != 2 {
			return rpki, fmt.Errorf("invalid rpki config line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		comm := strings.Split(value, ":")

		if key == "enabled" {
			rpki.Enabled = value == "true"
		} else if key == "valid" {
			rpki.Valid = append(rpki.Valid, comm)
		} else if key == "not_checked" {
			rpki.NotChecked = append(rpki.NotChecked, comm)
		} else if key == "invalid" {
			rpki.Invalid = append(rpki.Invalid, comm)
		} else if key == "unknown" {
			rpki.Unknown = append(rpki.Unknown, comm)
		} else {
			return rpki, fmt.Errorf("invalid rpki config line: %s", line)
		}

	}

	if err := section.MapTo(&rpki); err != nil {
		return rpki, err
	}
	hasDefaultASN := true
	asn, err := getDefaultASN(config)
	if err != nil {
		hasDefaultASN = false
	}

	// Fill in defaults or postprocess config value
	if len(rpki.Valid) == 0 && !hasDefaultASN && rpki.Enabled {
		return rpki, fmt.Errorf(
			"rpki.valid must be set if no server.asn is configured")
	}
	if len(rpki.Valid) == 0 && rpki.Enabled {
		log.Printf("Using default rpki.valid: %s:1000:1\n", asn)
		rpki.Valid = [][]string{{asn, "1000", "1"}}
	}

	if len(rpki.Unknown) == 0 && !hasDefaultASN && rpki.Enabled {
		return rpki, fmt.Errorf(
			"rpki.unknown must be set if no server.asn is configured")
	}
	if len(rpki.Unknown) == 0 && rpki.Enabled {
		log.Printf("Using default rpki.unknown: %s:1000:2\n", asn)
		rpki.Unknown = [][]string{{asn, "1000", "2"}}
	}

	if len(rpki.NotChecked) == 0 && !hasDefaultASN && rpki.Enabled {
		return rpki, fmt.Errorf(
			"rpki.not_checked must be set if no server.asn is set")
	}
	if len(rpki.NotChecked) == 0 {
		log.Printf("Using default rpki.not_checked: %s:1000:3\n", asn)
		rpki.NotChecked = [][]string{{asn, "1000", "3"}}
	}

	// As the euro-ix document states, this can be a range.
	for i, com := range rpki.Invalid {
		if len(com) != 3 {
			return rpki, fmt.Errorf("Invalid rpki.invalid config: %v", com)
		}
		tokens := strings.Split(com[2], "-")
		rpki.Invalid[i] = append([]string{com[0], com[1]}, tokens...)
	}
	if len(rpki.Invalid) == 0 && !hasDefaultASN && rpki.Enabled {
		return rpki, fmt.Errorf(
			"rpki.invalid must be set if no server.asn is configured")
	}
	if len(rpki.Invalid) == 0 && rpki.Enabled {
		log.Printf("Using default rpki.invalid: %s:1000:4-*\n", asn)
		rpki.Invalid = [][]string{{asn, "1000", "4", "*"}}
	}

	return rpki, nil
}

// Helper: Get own ASN from ini
// This is now easy, since we enforce an ASN in
// the [server] section.
func getDefaultASN(config *ini.File) (string, error) {
	server := config.Section("server")
	asn := server.Key("asn").MustString("")

	if asn == "" {
		return "", fmt.Errorf("could not get default ASN from config")
	}

	return asn, nil
}

// Get UI config: Theme settings
func getThemeConfig(config *ini.File) ThemeConfig {
	baseConfig := config.Section("theme")

	themeConfig := ThemeConfig{}
	_ = baseConfig.MapTo(&themeConfig)

	if themeConfig.BasePath == "" {
		themeConfig.BasePath = "/theme"
	}

	return themeConfig
}

// Get UI config: Pagination settings
func getPaginationConfig(config *ini.File) PaginationConfig {
	baseConfig := config.Section("pagination")

	paginationConfig := PaginationConfig{}
	_ = baseConfig.MapTo(&paginationConfig)

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

	// Get neighbors table columns
	neighborsColumns,
		neighborsColumnsOrder,
		err := getNeighborsColumns(config)
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

	// Blackhole communities
	blackholeCommunities, err := getBlackholeCommunities(config)
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

		NeighborsColumns:      neighborsColumns,
		NeighborsColumnsOrder: neighborsColumnsOrder,

		LookupColumns:      lookupColumns,
		LookupColumnsOrder: lookupColumnsOrder,

		RoutesRejections:       rejections,
		RoutesNoexports:        noexports,
		RoutesRejectCandidates: rejectCandidates,

		BGPBlackholeCommunities: blackholeCommunities,
		BGPCommunities:          getBGPCommunityMap(config),
		Rpki:                    rpki,

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
			return nil, fmt.Errorf("%s has ambiguous backends", section.Name())
		}

		// Configure backend
		backendConfig := sourceConfigSections[0]
		backendType, err := sourceBackendTypeFromConfig(backendConfig)
		if err != nil {
			return nil, fmt.Errorf("%s has an unsupported backend", section.Name())
		}

		sourceType := sourceTypeFromBackendType(backendType)

		// Make config
		sourceName := section.Key("name").MustString("Unknown Source")
		sourceGroup := section.Key("group").MustString("")
		sourceBlackholes := decoders.TrimmedCSVStringList(
			section.Key("blackholes").MustString(""))

		srcCfg := &SourceConfig{
			ID:         sourceID,
			Order:      order,
			Name:       sourceName,
			Group:      sourceGroup,
			Blackholes: sourceBlackholes,
			Backend:    backendType,
			Type:       sourceType,
		}

		// Register route server ID with pool
		pools.RouteServers.Acquire(sourceID)

		// Set backend
		switch backendType {
		case SourceBackendBirdwatcher:
			sourceType := backendConfig.Key("type").MustString("")
			mainTable := backendConfig.Key("main_table").MustString("master")
			peerTablePrefix := backendConfig.Key("peer_table_prefix").MustString("T")
			pipeProtocolPrefix := backendConfig.Key("pipe_protocol_prefix").MustString("M")

			if sourceType != "single_table" &&
				sourceType != "multi_table" {
				log.Fatal("Configuration error (birdwatcher source) unknown birdwatcher type:", sourceType)
			}

			c := birdwatcher.Config{
				ID:   srcCfg.ID,
				Name: srcCfg.Name,

				Timezone:        "UTC",
				ServerTime:      "2006-01-02T15:04:05.999999999Z07:00",
				ServerTimeShort: "2006-01-02",
				ServerTimeExt:   "Mon, 02 Jan 2006 15:04:05 -0700",

				Type:               sourceType,
				MainTable:          mainTable,
				PeerTablePrefix:    peerTablePrefix,
				PipeProtocolPrefix: pipeProtocolPrefix,
			}

			if err := backendConfig.MapTo(&c); err != nil {
				return nil, err
			}
			srcCfg.Birdwatcher = c

			log.Println("Adding birdwatcher source",
				c.Name, "of type", sourceType,
				"with peer_table_prefix", peerTablePrefix,
				"and pipe_protocol_prefix", pipeProtocolPrefix)
			if c.AltPipeProtocolSuffix != "" {
				log.Println(
					"Alt pipe protocol prefix:", c.AltPipeProtocolPrefix,
					"suffix:", c.AltPipeProtocolSuffix,
				)
			}

		case SourceBackendGoBGP:
			c := gobgp.Config{
				ID:   srcCfg.ID,
				Name: srcCfg.Name,
			}

			if err := backendConfig.MapTo(&c); err != nil {
				return nil, err
			}
			// Update defaults:
			//  - processing_timeout
			if c.ProcessingTimeout == 0 {
				c.ProcessingTimeout = 300
			}

			srcCfg.GoBGP = c

		case SourceBackendOpenBGPDStateServer:
			// Get cache TTL and reject communities from the config
			cacheTTL := time.Second * time.Duration(backendConfig.Key("cache_ttl").MustInt(300))
			routesCacheSize := backendConfig.Key("routes_cache_size").MustInt(1024)
			rc, err := getRoutesRejections(config)
			if err != nil {
				return nil, err
			}
			rejectComms := rc.Reasons.Communities()

			c := openbgpd.Config{
				ID:                srcCfg.ID,
				Name:              srcCfg.Name,
				CacheTTL:          cacheTTL,
				RoutesCacheSize:   routesCacheSize,
				RejectCommunities: rejectComms,
			}
			if err := backendConfig.MapTo(&c); err != nil {
				return nil, err
			}
			srcCfg.OpenBGPD = c

		case SourceBackendOpenBGPDBgplgd:
			// Get cache TTL from the config
			cacheTTL := time.Second * time.Duration(backendConfig.Key("cache_ttl").MustInt(300))
			routesCacheSize := backendConfig.Key("routes_cache_size").MustInt(1024)
			rc, err := getRoutesRejections(config)
			if err != nil {
				return nil, err
			}
			rejectComms := rc.Reasons.Communities()

			c := openbgpd.Config{
				ID:                srcCfg.ID,
				Name:              srcCfg.Name,
				CacheTTL:          cacheTTL,
				RoutesCacheSize:   routesCacheSize,
				RejectCommunities: rejectComms,
			}
			if err := backendConfig.MapTo(&c); err != nil {
				return nil, err
			}
			srcCfg.OpenBGPD = c
		}

		// Add to list of sources
		sources = append(sources, srcCfg)
		order++
	}

	return sources, nil
}

// preprocessConfig parses the variables in the config
// and applies it to the rest of the config.
func preprocessConfig(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	config := make([][]byte, 0, len(lines))

	expMap := ExpandMap{}
	for _, line := range lines {
		l := strings.TrimSpace(string(line))
		if strings.HasPrefix(l, "$") {
			expMap.AddExpr(l[1:])
			continue
		}
		config = append(config, line)
	}

	// Now apply to config
	configLines := []string{}
	for _, line := range config {
		l := string(line)
		exp, err := expMap.Expand(l)
		if err != nil {
			log.Fatal("Error expanding expression in config:", l, err)
		}
		for _, e := range exp {
			configLines = append(configLines, e)
		}
	}
	return []byte(strings.Join(configLines, "\n"))
}

// LoadConfig reads a configuration from a file.
func LoadConfig(file string) (*Config, error) {
	// Try to get config file, fallback to alternatives
	file, err := getConfigFile(file)
	if err != nil {
		return nil, err
	}

	// Read the config file and preprocess it
	configData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	configData = preprocessConfig(configData)

	// Load configuration, but handle bgp communities section
	// with our own parser
	parsedConfig, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{
			"bgp_communities",
			"blackhole_communities",
			"rejection_reasons",
			"noexport_reasons",
			"rpki",
		},
	}, configData)
	if err != nil {
		return nil, err
	}

	// Map sections
	server := ServerConfig{
		HTTPTimeout:                       DefaultHTTPTimeout,
		PrefixLookupCommunityFilterCutoff: DefaultPrefixLookupCommunityFilterCutoff,
		StoreBackend:                      "memory",
		RoutesStoreRefreshParallelism:     1,
		NeighborsStoreRefreshParallelism:  1,
		RoutesStoreQueryLimit:             DefaultRoutesStoreQueryLimit,
	}
	if err := parsedConfig.Section("server").MapTo(&server); err != nil {
		return nil, err
	}

	// Database config
	psql := &PostgresConfig{
		MinConns: 2,
		MaxConns: 128,
	}
	parsedConfig.Section("postgres").MapTo(&psql)
	if server.StoreBackend == "postgres" {
		if psql.URL == "" {
			psql.URL = "postgres:///?sslmode=prefer"
		}
	}

	housekeeping := HousekeepingConfig{}
	if err := parsedConfig.Section("housekeeping").MapTo(&housekeeping); err != nil {
		return nil, err
	}

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

	// Update stream parser throttle on all birdwatcher sources
	for _, src := range sources {
		if src.Backend == SourceBackendBirdwatcher {
			src.Birdwatcher.StreamParserThrottle = server.StreamParserThrottle
		}
	}

	config := &Config{
		Server:       server,
		Postgres:     psql,
		Housekeeping: housekeeping,
		UI:           ui,
		Sources:      sources,
		File:         file,
	}

	return config, nil
}

// GetInstance gets a source instance from config
func (cfg *SourceConfig) GetInstance() sources.Source {
	if cfg.instance != nil {
		return cfg.instance
	}

	var instance sources.Source
	switch cfg.Backend {
	case SourceBackendBirdwatcher:
		instance = birdwatcher.NewBirdwatcher(cfg.Birdwatcher)
	case SourceBackendGoBGP:
		instance = gobgp.NewGoBGP(cfg.GoBGP)
	case SourceBackendOpenBGPDStateServer:
		instance = openbgpd.NewStateServerSource(&cfg.OpenBGPD)
	case SourceBackendOpenBGPDBgplgd:
		instance = openbgpd.NewBgplgdSource(&cfg.OpenBGPD)
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
