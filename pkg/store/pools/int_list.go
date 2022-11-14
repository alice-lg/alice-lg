package pools

import "sync"

// IntNode is a node in the tree.
type IntNode struct {
	children []*IntNode
	value    int
	ptr      any
}

// IntTree is a tree structure for deduplicating
// lists of integers.
type IntTree struct {
	children []*IntNode
}

// Acquire int list
func (t *IntTree) Acquire(list []int) []int {
	if len(list) == 0 {
		return nil
	}
	v := list[0]
	var node *IntNode
	// TODO
	for _, c := range t.children {
		if c.value == v {
			// TODO
		}
	}
}

// A IntList pool can be used to deduplicate
// lists of integers. Like an AS path.
//
// A Tree datastructure is used.
type IntList struct {
	values  *IntTree
	counter *IntTree

	sync.Mutex
}
