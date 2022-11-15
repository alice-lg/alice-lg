package pools

import (
	"reflect"
	"sync"
)

// A IntListList pool can be used to deduplicate
// lists of lists of integers. Like a list of large
// BGP communities.
//
// A Tree datastructure is used.
type IntListList struct {
	lists *IntTree
	ptrs  *IntTree
	sync.Mutex
}

// NewIntListList creates a new int list pool
func NewIntListList() *IntListList {
	return &IntListList{
		lists: NewIntTree(),
		ptrs:  NewIntTree(),
	}
}

// Acquire int list from pool
func (p *IntListList) Acquire(list [][]int) [][]int {
	p.Lock()
	defer p.Unlock()

	// Convert list of list to list of ptrs
	ptrList := make([]int, len(list))
	for i, v := range list {
		ptr := p.ptrs.Acquire(v)
		ptrV := int(uintptr(reflect.ValueOf(ptr).UnsafePointer()))
		ptrList[i] = ptrV
	}

	return []int{}
}
