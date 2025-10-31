package store

import (
	"context"
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// ErrSourceNotInitialized is returned if a source
// is known but not yet initialized
var ErrSourceNotInitialized = errors.New(
	"source is not initialized")

// Store State Constants
const (
	StateInit = iota
	StateReady
	StateBusy
	StateError
)

// State is an enum of the above States
type State int

// String  converts a state into a string
func (s State) String() string {
	switch s {
	case StateInit:
		return "INIT"
	case StateReady:
		return "READY"
	case StateBusy:
		return "BUSY"
	case StateError:
		return "ERROR"
	}
	return "INVALID"
}

// Status defines a status the store can be in.
type Status struct {
	RefreshInterval     time.Duration `json:"refresh_interval"`
	RefreshParallelism  int           `json:"-"`
	LastRefresh         time.Time     `json:"last_refresh"`
	LastRefreshDuration time.Duration `json:"-"`
	LastError           interface{}   `json:"-"`
	State               State         `json:"state"`
	Initialized         bool          `json:"initialized"`
	SourceID            string        `json:"source_id"`

	lastRefreshStart time.Time
}

// SourceStatusList is a sortable list of source status
type SourceStatusList []*Status

// Len implements the sort interface
func (l SourceStatusList) Len() int {
	return len(l)
}

// Less implements the sort interface
func (l SourceStatusList) Less(i, j int) bool {
	return l[i].lastRefreshStart.Before(l[j].lastRefreshStart)
}

// Swap implements the sort interface
func (l SourceStatusList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// SourcesStore provides methods for retrieving
// the current status of a source.
type SourcesStore struct {
	refreshInterval    time.Duration
	refreshParallelism int
	status             map[string]*Status
	sources            map[string]*config.SourceConfig
	sync.Mutex
}

// NewSourcesStore initializes a new source store
func NewSourcesStore(
	cfg *config.Config,
	refreshInterval time.Duration,
	refreshParallelism int,
) *SourcesStore {
	status := make(map[string]*Status)
	sources := make(map[string]*config.SourceConfig)

	// Add sources from config
	for _, src := range cfg.Sources {
		sourceID := src.ID
		sources[sourceID] = src
		status[sourceID] = &Status{
			RefreshInterval: refreshInterval,
			SourceID:        sourceID,
		}
	}

	return &SourcesStore{
		status:             status,
		sources:            sources,
		refreshInterval:    refreshInterval,
		refreshParallelism: refreshParallelism,
	}
}

// GetSourcesStatus will retrieve the status for all sources
// as a list.
func (s *SourcesStore) GetSourcesStatus() []*Status {
	s.Lock()
	defer s.Unlock()
	status := []*Status{}
	for _, s := range s.status {
		status = append(status, s)
	}
	return status
}

// GetStatus will retrieve the status of a source
func (s *SourcesStore) GetStatus(sourceID string) (*Status, error) {
	s.Lock()
	defer s.Unlock()
	return s.getStatus(sourceID)
}

// Internal getStatus
func (s *SourcesStore) getStatus(sourceID string) (*Status, error) {
	status, ok := s.status[sourceID]
	if !ok {
		return nil, sources.ErrSourceNotFound
	}
	return status, nil
}

// IsInitialized will retrieve the status of the source
// and check if a successful refresh happened at least
// once.
func (s *SourcesStore) IsInitialized(sourceID string) (bool, error) {
	s.Lock()
	defer s.Unlock()
	status, err := s.getStatus(sourceID)
	if err != nil {
		return false, err
	}
	return status.Initialized, nil
}

// NextRefresh calculates the next refresh time
func (s *SourcesStore) NextRefresh(
	ctx context.Context,
) time.Time {
	s.Lock()
	defer s.Unlock()

	t := time.Time{}

	for _, status := range s.status {
		nextRefresh := status.LastRefresh.Add(
			s.refreshInterval)
		if nextRefresh.After(t) {
			t = nextRefresh
		}
	}
	return t
}

// ShouldRefresh checks if the source needs a
// new refresh according to the provided refreshInterval.
func (s *SourcesStore) ShouldRefresh(
	sourceID string,
) bool {
	status, err := s.GetStatus(sourceID)
	if err != nil {
		log.Println("get status error:", err)
		return false
	}

	nextRefresh := status.LastRefresh.Add(s.refreshInterval)

	if status.State == StateBusy {
		return false // Source is busy
	}
	if status.State == StateError {
		// The refresh interval in the config is ok if the
		// success case. When an error occurs it is desirable
		// to retry sooner, without spamming the server.
		nextRefresh = status.LastRefresh.Add(10 * time.Second)
	}

	if time.Now().UTC().Before(nextRefresh) {
		return false // Too soon
	}

	return true // Go for it
}

// CachedAt retrieves the oldest refresh time
// from all sources. All data is then guaranteed to be older
// than the CachedAt date.
func (s *SourcesStore) CachedAt(ctx context.Context) time.Time {
	s.Lock()
	defer s.Unlock()

	t := time.Now().UTC()
	for _, status := range s.status {
		if status.LastRefresh.Before(t) {
			t = status.LastRefresh
		}
	}
	return t
}

// GetInstance retrieves a source instance by ID
func (s *SourcesStore) GetInstance(sourceID string) sources.Source {
	s.Lock()
	defer s.Unlock()
	return s.sources[sourceID].GetInstance()
}

// GetName retrieves a source name by ID
func (s *SourcesStore) GetName(sourceID string) string {
	s.Lock()
	defer s.Unlock()
	return s.sources[sourceID].Name
}

// Get retrieves the source
func (s *SourcesStore) Get(sourceID string) *config.SourceConfig {
	s.Lock()
	defer s.Unlock()
	return s.sources[sourceID]
}

// GetSourceIDs returns a list of registered source ids.
func (s *SourcesStore) GetSourceIDs() []string {
	s.Lock()
	defer s.Unlock()
	ids := make([]string, 0, len(s.sources))
	for id := range s.sources {
		ids = append(ids, id)
	}
	return ids
}

// GetSourceIDsForRefresh will retrieve a list of source IDs,
// which are currently not locked, sorted by least refreshed.
// The number of sources returned is limited through the
// refresh parallelism parameter.
func (s *SourcesStore) GetSourceIDsForRefresh() []string {
	s.Lock()
	defer s.Unlock()

	locked := 0
	sources := make(SourceStatusList, 0, len(s.status))
	for _, status := range s.status {
		sources = append(sources, status)
		if status.State == StateBusy {
			locked++
		}
	}

	// Sort by refresh start time ascending
	sort.Sort(sources)

	slots := s.refreshParallelism - locked
	if slots <= 0 {
		slots = 0
	}

	ids := make([]string, 0, slots)
	i := 0
	for _, status := range sources {
		if i >= slots {
			break
		}
		ids = append(ids, status.SourceID)
		i++
	}
	return ids
}

// LockSource indicates the start of a refresh
func (s *SourcesStore) LockSource(sourceID string) error {
	s.Lock()
	defer s.Unlock()
	status, err := s.getStatus(sourceID)
	if err != nil {
		return err
	}
	if status.State == StateBusy {
		return sources.ErrSourceBusy
	}
	status.State = StateBusy
	status.lastRefreshStart = time.Now()
	return nil
}

// RefreshSuccess indicates a successful update
// of the store's content.
func (s *SourcesStore) RefreshSuccess(sourceID string) error {
	s.Lock()
	defer s.Unlock()
	status, err := s.getStatus(sourceID)
	if err != nil {
		return err
	}
	status.State = StateReady
	status.LastRefresh = time.Now().UTC()
	status.LastRefreshDuration = time.Since(status.lastRefreshStart)
	status.LastError = nil
	status.Initialized = true // We now have data
	return nil
}

// RefreshError indicates that the refresh has failed
func (s *SourcesStore) RefreshError(
	sourceID string,
	sourceErr interface{},
) {
	s.Lock()
	defer s.Unlock()
	status, err := s.getStatus(sourceID)
	if err != nil {
		log.Println("error getting source status:", err)
		return
	}
	status.State = StateError
	status.LastRefresh = time.Now().UTC()
	status.LastRefreshDuration = time.Since(status.lastRefreshStart)
	status.LastError = sourceErr
}
