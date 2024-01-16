package pools

import (
	"sync"
)

// A IntListPool can be used to deduplicate
// lists of integers. Like an AS path or BGP communities.
//
// A Tree datastructure is used.
type IntListPool struct {
	root    *Node[int, []int]
	counter uint64
	sync.Mutex
}

// NewIntListPool creates a new int list pool
func NewIntListPool() *IntListPool {
	return &IntListPool{
		root: NewNode[int, []int]([]int{}),
	}
}

// AcquireGid int list from pool and return with gid
func (p *IntListPool) AcquireGid(list []int) ([]int, uint64) {
	p.Lock()
	defer p.Unlock()

	if len(list) == 0 {
		return p.root.value, p.root.gid // root
	}
	v, c := p.root.traverse(p.counter+1, list, list)
	if c > p.counter {
		p.counter = c
	}
	return v, c
}

// Acquire int list from pool without gid
func (p *IntListPool) Acquire(list []int) []int {
	v, _ := p.AcquireGid(list)
	return v
}

// A StringListPool can be used for deduplicating lists
// of strings. (This is a variant of an int list, as string
// values are converted to int.
type StringListPool struct {
	root   *Node[int, []string]
	values map[string]int
	head   int
	sync.Mutex
}

// NewStringListPool creates a new string list.
func NewStringListPool() *StringListPool {
	return &StringListPool{
		head:   1,
		values: map[string]int{},
		root:   NewNode[int, []string]([]string{}),
	}
}

// AcquireGid aquires the string list pointer from the pool
// and also returns the gid.
func (p *StringListPool) AcquireGid(list []string) ([]string, uint64) {
	if len(list) == 0 {
		return p.root.value, p.root.gid
	}

	// Make idenfier list
	id := make([]int, len(list))
	for i, s := range list {
		// Resolve string value into int
		v, ok := p.values[s]
		if !ok {
			p.head++
			p.values[s] = p.head
			v = p.head
		}
		id[i] = v
	}

	return p.root.traverse(uint64(p.head), list, id)
}

// Acquire aquires the string list pointer from the pool
func (p *StringListPool) Acquire(list []string) []string {
	v, _ := p.AcquireGid(list)
	return v
}
