package pools

import (
	"sync"
)

// IntNode is a node in the tree.
type IntNode struct {
	children []*IntNode
	value    int
	counter  int
	ptr      []int
}

// Internally acquire list by traversing the tree and
// creating nodes if required.
func (n *IntNode) traverse(list, tail []int) []int {
	head := tail[0]
	tail = tail[1:]

	// Seek for value in children
	var child *IntNode
	for _, c := range n.children {
		if c.value == head {
			child = c
		}
	}
	if child == nil {
		// Insert child
		child = &IntNode{
			children: []*IntNode{},
			value:    head,
			ptr:      nil,
		}
		n.children = append(n.children, child)
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

// IntTree is a tree structure for deduplicating
// lists of integers.
type IntTree struct {
	root *IntNode
}

// NewIntTree initializes an empty int tree
func NewIntTree() *IntTree {
	return &IntTree{
		root: &IntNode{
			ptr: []int{},
		},
	}
}

// Acquire int list element
func (t *IntTree) Acquire(list []int) []int {
	if len(list) == 0 {
		return t.root.ptr // root
	}
	return t.root.traverse(list, list)
}

// A IntList pool can be used to deduplicate
// lists of integers. Like an AS path.
//
// A Tree datastructure is used.
type IntList struct {
	values *IntTree
	sync.Mutex
}

// NewIntList creates a new int list pool
func NewIntList() *IntList {
	return &IntList{
		values: NewIntTree(),
	}
}

// Acquire int list from pool
func (p *IntList) Acquire(list []int) []int {
	p.Lock()
	defer p.Unlock()
	return p.values.Acquire(list)
}
