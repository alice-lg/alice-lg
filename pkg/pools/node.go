package pools

import (
	"unsafe"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// NOTE: Yes, generics could be used here.
// This also looks like a pretty good use case for them.
// However, the performance penalty is too high.
//
// With:    Refreshed routes of rs1.foo (BIRD2 IPv4) in 26.905169348s
// Without: Refreshed routes of rs1.foo (BIRD2 IPv4) in 46.760695927s
//
// So yeah. Copy and paste time!

type Node[T comparable, V any] struct {
	children map[T]*Node[T, V] // map of children
	value    V
	final    bool
}

// NewNode creates a new tree node
func NewNode[T comparable, V any](value V) *Node[T, V] {
	return &Node[T, V]{
		children: map[T]*Node[T, V]{},
		value:    value,
		final:    false,
	}
}

// traverse inserts a new node into the three if required
// or returns the object if it already exists.
func (n *Node[T, V]) traverse(value V, tail []T) V {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		var zero V
		child = NewNode[T, V](zero)
		n.children[id] = child
	}

	// Set obj if required
	if len(tail) == 0 {
		if !child.final {
			child.value = value
			child.final = true
		}
		return child.value
	}

	return child.traverse(value, tail)
}

// read returns the object if it exists or nil if not.
func (n *Node[T, V]) read(tail []T) V {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		var zero V
		return zero
	}

	// Set obj if required
	if len(tail) == 0 {
		return child.value
	}

	return child.read(tail)
}

// IntListNode is a node with an integer as key and
// a list of integers as value.
type IntListNode struct {
	children map[int]*IntListNode
	value    []int
}

// NewIntListNode creates a new int tree node
func NewIntListNode(value []int) *IntListNode {
	return &IntListNode{
		children: map[int]*IntListNode{},
		value:    value,
	}
}

// IntList traverse inserts a new node into the three if required
// or returns the object if it already exists.
func (n *IntListNode) traverse(value []int, tail []int) []int {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		child = NewIntListNode(nil)
		n.children[id] = child
	}

	// Set obj if required
	if len(tail) == 0 {
		if child.value == nil {
			child.value = value
		}
		return child.value
	}

	return child.traverse(value, tail)
}

// read returns the object if it exists or nil if not.
func (n *IntListNode) read(tail []int) []int {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		return nil
	}

	// Set obj if required
	if len(tail) == 0 {
		return child.value
	}

	return child.read(tail)
}

// StringListNode is a node with an integer as key and
// a list of integers as value.
type StringListNode struct {
	children map[int]*StringListNode
	value    []string
}

// NewStringListNode creates a new int tree node
func NewStringListNode(value []string) *StringListNode {
	return &StringListNode{
		children: map[int]*StringListNode{},
		value:    value,
	}
}

// StringList traverse inserts a new node into the three if required
// or returns the object if it already exists.
func (n *StringListNode) traverse(value []string, tail []int) []string {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		child = NewStringListNode(nil)
		n.children[id] = child
	}

	// Set obj if required
	if len(tail) == 0 {
		if child.value == nil {
			child.value = value
		}
		return child.value
	}

	return child.traverse(value, tail)
}

// read returns the object if it exists or nil if not.
func (n *StringListNode) read(tail []int) []string {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		return nil
	}

	// Set obj if required
	if len(tail) == 0 {
		return child.value
	}

	return child.read(tail)
}

// CommunityNode is a node with an integer as key
type CommunityNode struct {
	children map[int]*CommunityNode
	value    api.Community
}

// NewCommunityNode creates a new int tree node
func NewCommunityNode(value []int) *CommunityNode {
	return &CommunityNode{
		children: map[int]*CommunityNode{},
		value:    value,
	}
}

// CommunityNode: traverse inserts a new node into the three if required
// or returns the object if it already exists.
func (n *CommunityNode) traverse(value api.Community, tail []int) api.Community {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		child = NewCommunityNode(nil)
		n.children[id] = child
	}

	// Set obj if required
	if len(tail) == 0 {
		if child.value == nil {
			child.value = value
		}
		return child.value
	}

	return child.traverse(value, tail)
}

// read returns the object if it exists or nil if not.
func (n *CommunityNode) read(tail []int) api.Community {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		return nil
	}

	// Set obj if required
	if len(tail) == 0 {
		return child.value
	}

	return child.read(tail)
}

// CommunityListNode is a node with an integer as key and
// a list of integers as value.
type CommunityListNode struct {
	children map[unsafe.Pointer]*CommunityListNode
	value    []api.Community
}

// NewCommunityListNode creates a new int tree node
func NewCommunityListNode(value []api.Community) *CommunityListNode {
	return &CommunityListNode{
		children: map[unsafe.Pointer]*CommunityListNode{},
		value:    value,
	}
}

// CommunityList traverse inserts a new node into the three if required
// or returns the object if it already exists.
func (n *CommunityListNode) traverse(
	value []api.Community,
	tail []unsafe.Pointer,
) []api.Community {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		child = NewCommunityListNode(nil)
		n.children[id] = child
	}

	// Set obj if required
	if len(tail) == 0 {
		if child.value == nil {
			child.value = value
		}
		return child.value
	}

	return child.traverse(value, tail)
}

// read returns the object if it exists or nil if not.
func (n *CommunityListNode) read(tail []unsafe.Pointer) []api.Community {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		return nil
	}

	// Set obj if required
	if len(tail) == 0 {
		return child.value
	}

	return child.read(tail)
}
