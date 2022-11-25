package pools

import (
	"math"
	"reflect"
	"sync"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// CommunitiesPool is for deduplicating a list of BGP communities
// (Large and default. The ext communities representation right now
// makes problems and need to be fixed. TODO.)
type CommunitiesPool struct {
	communitiesRoot *Node
	root            *Node
	sync.Mutex
}

// NewCommunitiesPool creates a new pool for lists
// of BGP communities.
func NewCommunitiesPool() *CommunitiesPool {
	return &CommunitiesPool{
		communitiesRoot: NewNode([]int{}),
		root:            NewNode([]api.Community{}),
	}
}

// Acquire a list of bgp communities
func (p *CommunitiesPool) Acquire(communities []api.Community) []api.Community {
	p.Lock()
	defer p.Unlock()
	// Make identification list by using the pointer address
	// of the deduplicated community as ID
	ids := make([]int, len(communities))
	for i, comm := range communities {
		commPtr := p.communitiesRoot.traverse(comm, comm)
		addr := reflect.ValueOf(commPtr).UnsafePointer()
		ids[i] = int(uintptr(addr))
	}
	if len(ids) == 0 {
		return p.root.ptr.([]api.Community)
	}
	return p.root.traverse(communities, ids).([]api.Community)
}

// NewExtCommunitiesPool creates a new pool for extended communities
func NewExtCommunitiesPool() *CommunitiesPool {
	return &CommunitiesPool{
		communitiesRoot: NewNode([]int{}),
		root:            NewNode([]api.ExtCommunity{}),
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
func (p *CommunitiesPool) AcquireExt(communities []api.ExtCommunity) []api.ExtCommunity {
	p.Lock()
	defer p.Unlock()

	// Make identification list
	ids := make([]int, len(communities))
	for i, comm := range communities {
		r := extPrefixToInt(comm[0].(string))
		icomm := []int{r, comm[1].(int), comm[2].(int)}

		// get community identifier
		commPtr := p.communitiesRoot.traverse(icomm, icomm)
		addr := reflect.ValueOf(commPtr).UnsafePointer()
		ids[i] = int(uintptr(addr))
	}
	if len(ids) == 0 {
		return p.root.ptr.([]api.ExtCommunity)
	}
	return p.root.traverse(communities, ids).([]api.ExtCommunity)
}
