package store

import (
	"sync"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/sources"
)

// SourcesRegistry holds the live set of source configurations
// that the API and housekeeping read, so a configuration reload
// is visible to every consumer and not only to the prefix search
// stores. It is reconciled ahead of the stores on a reload and is
// safe for concurrent use.
type SourcesRegistry struct {
	mu      sync.RWMutex
	sources []*config.SourceConfig
	byID    map[string]*config.SourceConfig
}

// NewSourcesRegistry seeds the registry from the initial
// configuration.
func NewSourcesRegistry(cfg *config.Config) *SourcesRegistry {
	byID := make(map[string]*config.SourceConfig, len(cfg.Sources))
	for _, src := range cfg.Sources {
		byID[src.ID] = src
	}
	return &SourcesRegistry{
		sources: cfg.Sources,
		byID:    byID,
	}
}

// Reconcile aligns the registry with the provided configuration,
// keeping an unchanged source's pointer so its memoized client
// survives; a source that only moved takes its new position.
func (r *SourcesRegistry) Reconcile(cfg *config.Config) {
	r.mu.Lock()
	defer r.mu.Unlock()

	next := make([]*config.SourceConfig, 0, len(cfg.Sources))
	byID := make(map[string]*config.SourceConfig, len(cfg.Sources))
	for _, desired := range cfg.Sources {
		src := desired
		if held, ok := r.byID[desired.ID]; ok && held.Equal(desired) {
			src = held // keep the memoized instance
			if held.Order != desired.Order {
				src = held.WithOrder(desired.Order)
			}
		}
		next = append(next, src)
		byID[src.ID] = src
	}
	r.sources = next
	r.byID = byID
}

// Sources returns a snapshot of the current source set.
func (r *SourcesRegistry) Sources() []*config.SourceConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*config.SourceConfig, len(r.sources))
	copy(out, r.sources)
	return out
}

// SourceByID returns the source configuration for an ID, or nil.
func (r *SourcesRegistry) SourceByID(id string) *config.SourceConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byID[id]
}

// SourceInstanceByID returns the source instance for an ID, or
// nil if the source is not registered.
func (r *SourcesRegistry) SourceInstanceByID(id string) sources.Source {
	r.mu.RLock()
	src := r.byID[id]
	r.mu.RUnlock()
	if src == nil {
		return nil
	}
	return src.GetInstance()
}
