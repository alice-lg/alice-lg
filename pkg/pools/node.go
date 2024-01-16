package pools

type Node[T comparable, V any] struct {
	children map[T]*Node[T, V] // map of children
	value    V
	final    bool
	gid      uint64
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
func (n *Node[T, V]) traverse(gid uint64, value V, tail []T) (V, uint64) {
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
			child.gid = gid
		}
		return child.value, child.gid
	}

	return child.traverse(gid, value, tail)
}

// read returns the object if it exists or nil if not.
func (n *Node[T, V]) read(tail []T) (V, uint64) {
	id := tail[0]
	tail = tail[1:]

	// Seek for identifier in children
	child, ok := n.children[id]
	if !ok {
		var zero V
		return zero, 0
	}

	// Set obj if required
	if len(tail) == 0 {
		return child.value, child.gid
	}

	return child.read(tail)
}
