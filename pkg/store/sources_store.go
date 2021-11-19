package store

import (
	"sync"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// Store State Constants
const (
	StateInit = iota
	StateReady
	StateUpdatingNeighbors
	StateUpdatingRoutes
	StateError
)

// Status defines a status the store can be in
type Status struct {
	LastNeighborsRefresh time.Time
	LastRoutesRefresh    time.Time
	LastError            error
	State                int
}

// stateToString converts a state into a string
func stateToString(state int) string {
	switch state {
	case StateInit:
		return "INIT"
	case StateReady:
		return "READY"
	case StateUpdating:
		return "UPDATING"
	case StateError:
		return "ERROR"
	}
	return "INVALID"
}

// SourcesStore provides methods for retrieving
// the current status of a source.
type SourcesStore struct {
	status  map[string]*Status
	sources map[string]*config.SourceConfig

	sync.Mutex
}

// NewSourcesStore initializes a new source store
func NewSourcesStore() *SourcesStore {
	return &SourcesStore{
		status:  make(map[string]*Status),
		sources: make(map[string]*config.SourceConfig),
	}
}

// AddSource registers the source
func (s *SourcesStore) AddSource(src *config.SourceConfig) {
	s.Lock()
	defer s.Unlock()
	sourceID := src.ID
	sources[sourceID] = src
	status[sourceID] = &Status{}
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
}

// ShouldRefreshNeighbors checks if the source needs a
// new refresh according to the provided refreshInterval.
func (s *SourcesStore) ShouldRefreshNeighbors(
	sourceID string,
	interval time.Duration,
) bool {
	status = s.GetStatus(sourceID)
	if status.State == StateUpdating {
		return false // Source is busy
	}
	nextRefresh := status.LastNeighborsRefresh.Add(interval)
	if time.Now().UTC().Before(nextRefresh) {
		return false // Too soon
	}
	return true // Go for it
}

// ShouldRefreshRoutes checks if the source needs a
// new refresh according to the provided refreshInterval.
func (s *SourcesStore) ShouldRefreshRoutes(
	sourceID string, interval time.Duration,
) bool {
	status = s.GetStatus(sourceID)
	if status.State == StateUpdating {
		return false // Source is busy
	}
	nextRefresh := status.LastRoutesRefresh.Add(interval)
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
	if status.State == StateUpdating {
		return sources.ErrSourceBusy
	}
	status.State = StateUpdating
	return nil
}

// NeighborsRefreshSuccess indicates a successfull refresh
func (s *SourcesStore) NeighborsRefreshSuccess(sourceID string) error {
	s.Lock()
	defer s.Unlock()
	status, err := getStatus(sourceID)
	if err != nil {
		return err
	}
	status.State = StateReady
	status.LastNeighborsRefresh = time.Now().UTC()
	status.LastError = nil
}

// NeighborsRefreshError indicates that the refresh has failed
func (s *SourcesStore) NeighborsRefreshError(sourceID string, err error) error {
	s.Lock()
	defer s.Unlock()
	status, err := getStatus(sourceID)
	if err != nil {
		return err
	}
	status.State = StateError
	status.LastNeighborsRefresh = time.Now().UTC()
	status.LastError = err.String()
}
