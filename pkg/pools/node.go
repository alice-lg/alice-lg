package pools

// Node is a generic tree node
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
