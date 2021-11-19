package store

import (
	"log"
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// Store State Constants
const (
	StateInit = iota
	StateReady
	StateBusy
	StateError
)

// State is an enum of the above States
type State int

// String()  converts a state into a string
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

// Status defines a status the store can be in
type Status struct {
	RefreshInterval     time.Duration
	LastRefresh         time.Time
	LastRefreshDuration time.Duration
	LastError           interface{}
	State               State

	lastRefreshStart time.Time
}

// SourcesStore provides methods for retrieving
// the current status of a source.
type SourcesStore struct {
	refreshInterval time.Duration
	status          map[string]*Status
	sources         map[string]*config.SourceConfig
	sync.Mutex
}

// NewSourcesStore initializes a new source store
func NewSourcesStore(
	cfg *config.Config,
	refreshInterval time.Duration,
) *SourcesStore {
	status := make(map[string]*Status)
	sources := make(map[string]*config.SourceConfig)

	// Add sources from config
	for _, src := range cfg.Sources {
		sourceID := src.ID
		sources[sourceID] = src
		status[sourceID] = &Status{
			RefreshInterval: refreshInterval,
		}
	}

	return &SourcesStore{
		status:          status,
		sources:         sources,
		refreshInterval: refreshInterval,
	}
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

// NextRefresh calculates the next refresh time
func (s *SourcesStore) NextRefresh(sourceID string) time.Time {
	status, err := s.GetStatus(sourceID)
	if err != nil {
		log.Println("get status error:", err)
		return false
	}
	if status.State == StateBusy {
		return false // Source is busy
	}
	nextRefresh := status.LastRefresh.Add(
		s.refreshInterval)
	return nextRefresh
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
	if status.State == StateBusy {
		return false // Source is busy
	}
	nextRefresh := status.LastRefresh.Add(
		s.refreshInterval)
	if time.Now().UTC().Before(nextRefresh) {
		return false // Too soon
	}
	return true // Go for it
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

// RefreshSuccess indicates a successfull update
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
	status.LastRefreshDuration = time.Now().Sub(
		status.lastRefreshStart)
	status.LastError = nil
	return nil
}

// RefreshError indicates that the refresh has failed
func (s *SourcesStore) RefreshError(
	sourceID string,
	err interface{},
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
	status.LastRefreshDuration = time.Now().Sub(
		status.lastRefreshStart)
	status.LastError = err
	return
}
