package pools

import (
	"sync"
)

// Node is a node in the tree.
type Node struct {
	children map[int]*Node
	counter  int
	ptr      interface{}
}

// NewNode creates a new tree node
func NewNode(ptr interface{}) *Node {
	return &Node{
		children: map[int]*Node{},
		ptr:      ptr,
	}
}

// Internally acquire list by traversing the tree and
// creating nodes if required.
func (n *Node) traverse(list interface{}, tail []int) interface{} {
	value := tail[0]
	tail = tail[1:]

	// Seek for value in children
	child, ok := n.children[value]
	if !ok {
		child = NewNode(nil)
		n.children[value] = child
	}

	// Set list ptr if required
	if len(tail) == 0 {
		if child.ptr == nil {
			child.ptr = list
		}
		return child.ptr
	}

	return child.traverse(list, tail)
}

// A IntListPool can be used to deduplicate
// lists of integers. Like an AS path or BGP communities.
//
// A Tree datastructure is used.
type IntListPool struct {
	root *Node
	sync.Mutex
}

// NewIntListPool creates a new int list pool
func NewIntListPool() *IntListPool {
	return &IntListPool{
		root: NewNode([]int{}),
	}
}

// Acquire int list from pool
func (p *IntListPool) Acquire(list []int) []int {
	p.Lock()
	defer p.Unlock()

	if len(list) == 0 {
		return p.root.ptr.([]int) // root
	}
	return p.root.traverse(list, list).([]int)
}

// A StringListPool can be used for deduplicating lists
// of strings. (This is a variant of an int list, as string
// values are converted to int.
type StringListPool struct {
	root   *Node
	values map[string]int
	head   int
	sync.Mutex
}

// NewStringListPool creates a new string list.
func NewStringListPool() *StringListPool {
	return &StringListPool{
		head:   1,
		values: map[string]int{},
		root:   NewNode([]string{}),
	}
}

// Acquire the string list pointer from the pool
func (p *StringListPool) Acquire(list []string) []string {
	if len(list) == 0 {
		return p.root.ptr.([]string) // root
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

	return p.root.traverse(list, id).([]string)
}
