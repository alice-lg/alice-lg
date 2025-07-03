package pools

import (
	"sync"
)

// A IntListPool can be used to deduplicate
// lists of integers. Like an AS path or BGP communities.
//
// A Tree datastructure is used.
type IntListPool struct {
	root *Node[int, []int]
	sync.Mutex
}

// NewIntListPool creates a new int list pool
func NewIntListPool() *IntListPool {
	return &IntListPool{
		root: NewNode[int, []int]([]int{}),
	}
}

// Acquire int list from pool
func (p *IntListPool) Acquire(list []int) []int {
	p.Lock()
	defer p.Unlock()

	if len(list) == 0 {
		return p.root.value // root
	}
	return p.root.traverse(list, list)
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

// Acquire the string list pointer from the pool.
func (p *StringListPool) Acquire(list []string) []string {
	if len(list) == 0 {
		return p.root.value
	}

	// Make identifier list
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

	return p.root.traverse(list, id)
}
