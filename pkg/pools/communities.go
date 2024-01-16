package pools

import (
	"math"
	"reflect"
	"sync"
	"unsafe"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// CommunitiesPool is for deduplicating a single BGP community.
// This works with large and standard communities. For extended
// communities, use the ExtCommunityPool.
type CommunitiesPool struct {
	root *Node[int, api.Community]
	sync.RWMutex
}

// NewCommunitiesPool creates a new pool for a single BGP community
func NewCommunitiesPool() *CommunitiesPool {
	return &CommunitiesPool{
		root: NewNode[int, api.Community](api.Community{}),
	}
}

// Acquire a single bgp community
func (p *CommunitiesPool) Acquire(c api.Community) api.Community {
	p.Lock()
	defer p.Unlock()
	if len(c) == 0 {
		return p.root.value
	}
	return p.root.traverse(c, c)
}

// Read a single bgp community
func (p *CommunitiesPool) Read(c api.Community) api.Community {
	p.RLock()
	defer p.RUnlock()
	if len(c) == 0 {
		return p.root.value // root
	}
	v := p.root.read(c)
	if v == nil {
		return nil
	}
	return v
}

// CommunitiesSetPool is for deduplicating a list of BGP communities
// (Large and default. The ext communities representation right now
// makes problems and need to be fixed. TODO.)
type CommunitiesSetPool struct {
	root *Node[unsafe.Pointer, []api.Community]
	sync.Mutex
}

// NewCommunitiesSetPool creates a new pool for lists
// of BGP communities.
func NewCommunitiesSetPool() *CommunitiesSetPool {
	return &CommunitiesSetPool{
		root: NewNode[unsafe.Pointer, []api.Community]([]api.Community{}),
	}
}

// Acquire a list of bgp communities
func (p *CommunitiesSetPool) Acquire(communities []api.Community) []api.Community {
	p.Lock()
	defer p.Unlock()
	// Make identification list by using the pointer address
	// of the deduplicated community as ID
	ids := make([]unsafe.Pointer, len(communities))
	set := make([]api.Community, len(communities))
	for i, comm := range communities {
		commPtr := Communities.Acquire(comm)
		ids[i] = reflect.ValueOf(commPtr).UnsafePointer()
		set[i] = commPtr
	}
	if len(ids) == 0 {
		return p.root.value
	}
	return p.root.traverse(set, ids)
}

// ExtCommunitiesSetPool is for deduplicating a list of ext. BGP communities
type ExtCommunitiesSetPool struct {
	root *Node[unsafe.Pointer, []api.ExtCommunity]
	sync.Mutex
}

// NewExtCommunitiesSetPool creates a new pool for lists
// of BGP communities.
func NewExtCommunitiesSetPool() *ExtCommunitiesSetPool {
	return &ExtCommunitiesSetPool{
		root: NewNode[unsafe.Pointer, []api.ExtCommunity]([]api.ExtCommunity{}),
	}
}

func extPrefixToInt(s string) int {
	v := 0
	for i, c := range s {
		v += int(math.Pow(1000.0, float64(i))) * int(c)
	}
	return v
}

// AcquireExt a list of ext bgp communities
func (p *ExtCommunitiesSetPool) Acquire(communities []api.ExtCommunity) []api.ExtCommunity {
	p.Lock()
	defer p.Unlock()

	// Make identification list
	ids := make([]unsafe.Pointer, len(communities))
	for i, comm := range communities {
		r := extPrefixToInt(comm[0].(string))
		icomm := []int{r, comm[1].(int), comm[2].(int)}

		// get community identifier
		commPtr := ExtCommunities.Acquire(icomm)
		ids[i] = reflect.ValueOf(commPtr).UnsafePointer()
	}
	if len(ids) == 0 {
		return p.root.value
	}
	return p.root.traverse(communities, ids)
}
