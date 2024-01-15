package pools

import (
	"math"
	"reflect"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// CommunitiesPool is for deduplicating a single BGP community.
// This works with large and standard communities. For extended
// communities, use the ExtCommunityPool.
type CommunitiesPool struct {
	root *Node
	sync.RWMutex
}

// NewCommunitiesPool creates a new pool for a single BGP community
func NewCommunitiesPool() *CommunitiesPool {
	return &CommunitiesPool{
		root: NewNode(api.Community{}),
	}
}

// Acquire a single bgp community
func (p *CommunitiesPool) Acquire(c api.Community) api.Community {
	p.Lock()
	defer p.Unlock()
	if len(c) == 0 {
		return p.root.ptr.(api.Community) // root
	}
	return p.root.traverse(c, c).(api.Community)
}

// Read a single bgp community
func (p *CommunitiesPool) Read(c api.Community) api.Community {
	p.RLock()
	defer p.RUnlock()
	if len(c) == 0 {
		return p.root.ptr.(api.Community) // root
	}
	v := p.root.read(c, c)
	if v == nil {
		return nil
	}
	return v.(api.Community)
}

// CommunitiesSetPool is for deduplicating a list of BGP communities
// (Large and default. The ext communities representation right now
// makes problems and need to be fixed. TODO.)
type CommunitiesSetPool struct {
	root *Node
	sync.Mutex
}

// NewCommunitiesSetPool creates a new pool for lists
// of BGP communities.
func NewCommunitiesSetPool() *CommunitiesSetPool {
	return &CommunitiesSetPool{
		root: NewNode([]api.Community{}),
	}
}

// Acquire a list of bgp communities
func (p *CommunitiesSetPool) Acquire(communities []api.Community) []api.Community {
	p.Lock()
	defer p.Unlock()
	// Make identification list by using the pointer address
	// of the deduplicated community as ID
	ids := make([]int, len(communities))
	set := make([]api.Community, len(communities))
	for i, comm := range communities {
		commPtr := Communities.Acquire(comm)
		addr := reflect.ValueOf(commPtr).UnsafePointer()
		ids[i] = int(uintptr(addr))
		set[i] = commPtr
	}
	if len(ids) == 0 {
		return p.root.ptr.([]api.Community)
	}
	return p.root.traverse(set, ids).([]api.Community)
}

// NewExtCommunitiesSetPool creates a new pool for extended communities
func NewExtCommunitiesSetPool() *CommunitiesSetPool {
	return &CommunitiesSetPool{
		root: NewNode([]api.ExtCommunity{}),
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
func (p *CommunitiesSetPool) AcquireExt(communities []api.ExtCommunity) []api.ExtCommunity {
	p.Lock()
	defer p.Unlock()

	// Make identification list
	ids := make([]int, len(communities))
	for i, comm := range communities {
		r := extPrefixToInt(comm[0].(string))
		icomm := []int{r, comm[1].(int), comm[2].(int)}

		// get community identifier
		commPtr := ExtCommunities.Acquire(icomm)
		addr := reflect.ValueOf(commPtr).UnsafePointer()
		ids[i] = int(uintptr(addr))
	}
	if len(ids) == 0 {
		return p.root.ptr.([]api.ExtCommunity)
	}
	return p.root.traverse(communities, ids).([]api.ExtCommunity)
}
