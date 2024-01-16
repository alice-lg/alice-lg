package pools

import (
	"math"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// CommunitiesPool is for deduplicating a single BGP community.
// This works with large and standard communities. For extended
// communities, use the ExtCommunityPool.
type CommunitiesPool struct {
	root    *Node[int, api.Community]
	counter uint64
	sync.RWMutex
}

// NewCommunitiesPool creates a new pool for a single BGP community
func NewCommunitiesPool() *CommunitiesPool {
	return &CommunitiesPool{
		root: NewNode[int, api.Community](api.Community{}),
	}
}

// AcquireGid acquires a single bgp community with gid
func (p *CommunitiesPool) AcquireGid(c api.Community) (api.Community, uint64) {
	p.Lock()
	defer p.Unlock()
	if len(c) == 0 {
		return p.root.value, p.root.gid
	}
	v, gid := p.root.traverse(p.counter+1, c, c)
	if gid > p.counter {
		p.counter = gid
	}
	return v, gid
}

// Acquire a single bgp community without gid
func (p *CommunitiesPool) Acquire(c api.Community) api.Community {
	v, _ := p.AcquireGid(c)
	return v
}

// Read a single bgp community
func (p *CommunitiesPool) Read(c api.Community) (api.Community, uint64) {
	p.RLock()
	defer p.RUnlock()
	if len(c) == 0 {
		return p.root.value, p.root.gid
	}
	return p.root.read(c)
}

// CommunitiesSetPool is for deduplicating a list of BGP communities
// (Large and default. The ext communities representation right now
// makes problems and need to be fixed. TODO.)
type CommunitiesSetPool struct {
	root    *Node[uint64, []api.Community]
	counter uint64
	sync.Mutex
}

// NewCommunitiesSetPool creates a new pool for lists
// of BGP communities.
func NewCommunitiesSetPool() *CommunitiesSetPool {
	return &CommunitiesSetPool{
		root: NewNode[uint64, []api.Community]([]api.Community{}),
	}
}

// AcquireGid acquires a list of bgp communities and returns a gid
func (p *CommunitiesSetPool) AcquireGid(
	communities []api.Community,
) ([]api.Community, uint64) {
	p.Lock()
	defer p.Unlock()
	// Make identification list by using the pointer address
	// of the deduplicated community as ID
	ids := make([]uint64, len(communities))
	set := make([]api.Community, len(communities))
	for i, comm := range communities {
		ptr, gid := Communities.AcquireGid(comm)
		ids[i] = gid
		set[i] = ptr
	}
	if len(ids) == 0 {
		return p.root.value, p.root.gid
	}
	v, id := p.root.traverse(p.counter+1, set, ids)
	if id > p.counter {
		p.counter = id
	}
	return v, id
}

// Acquire a list of bgp communities
func (p *CommunitiesSetPool) Acquire(
	communities []api.Community,
) []api.Community {
	v, _ := p.AcquireGid(communities)
	return v
}

// ExtCommunitiesSetPool is for deduplicating a list of ext. BGP communities
type ExtCommunitiesSetPool struct {
	root    *Node[uint64, []api.ExtCommunity]
	counter uint64
	sync.Mutex
}

// NewExtCommunitiesSetPool creates a new pool for lists
// of BGP communities.
func NewExtCommunitiesSetPool() *ExtCommunitiesSetPool {
	return &ExtCommunitiesSetPool{
		root: NewNode[uint64, []api.ExtCommunity]([]api.ExtCommunity{}),
	}
}

func extPrefixToInt(s string) int {
	v := 0
	for i, c := range s {
		v += int(math.Pow(1000.0, float64(i))) * int(c)
	}
	return v
}

// AcquireGid acquires a list of ext bgp communities
func (p *ExtCommunitiesSetPool) AcquireGid(
	communities []api.ExtCommunity,
) ([]api.ExtCommunity, uint64) {
	p.Lock()
	defer p.Unlock()

	// Make identification list
	ids := make([]uint64, len(communities))
	for i, comm := range communities {
		r := extPrefixToInt(comm[0].(string))
		icomm := []int{r, comm[1].(int), comm[2].(int)}

		// get community identifier
		_, gid := ExtCommunities.AcquireGid(icomm)
		ids[i] = gid
	}
	if len(ids) == 0 {
		return p.root.value, p.root.gid
	}
	v, id := p.root.traverse(p.counter+1, communities, ids)
	if id > p.counter {
		p.counter = id
	}
	return v, id
}

// Acquire a list of ext bgp communities
func (p *ExtCommunitiesSetPool) Acquire(
	communities []api.ExtCommunity,
) []api.ExtCommunity {
	v, _ := p.AcquireGid(communities)
	return v
}
