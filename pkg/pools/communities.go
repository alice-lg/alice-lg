package pools

import (
	"reflect"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// Communities is a pool for deduplicating a list of BGP communities
// (Large and default. The ext communities representation right now
// makes problems and need to be fixed. TODO.)
type Communities struct {
	communities *IntList
	root        *Node
	sync.Mutex
}

// NewCommunities creates a new pool for lists
// of BGP communities.
func NewCommunities() *Communities {
	return &Communities{
		communities: NewIntList(),
		root: &Node{
			ptr: []api.Community{},
		},
	}
}

// Acquire a list of bgp communities
func (p *Communities) Acquire(communities []api.Community) []api.Community {
	ids := make([]int, len(communities))
	for i, comm := range communities {
		commPtr := p.communities.Acquire(comm)
		addr := reflect.ValueOf(commPtr).UnsafePointer()
		ids[i] = int(uintptr(addr))
	}
	p.Lock()
	defer p.Unlock()
	if len(ids) == 0 {
		return p.root.ptr.([]api.Community)
	}

	return p.root.traverse(communities, ids).([]api.Community)
}
