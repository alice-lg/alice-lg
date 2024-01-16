package pools

import (
	"fmt"
	"testing"
)

func TestAcquireIntList(t *testing.T) {
	a := []int{23, 42, 1337, 65535, 1}
	b := []int{23, 42, 1337, 65535, 1}
	c := []int{23, 42, 1338, 65535, 2}

	p := NewIntListPool()

	r1 := p.Acquire(a)
	p.Acquire(c)
	r2 := p.Acquire(b)

	if fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b) {
		t.Error("lists should not be same pointer", fmt.Sprintf("%p %p", a, b))
	}
	if fmt.Sprintf("%p", r1) != fmt.Sprintf("%p", r2) {
		t.Error("lists should be same pointer", fmt.Sprintf("%p %p", r1, r2))
	}
	t.Log(fmt.Sprintf("Ptr: %p %p => %p %p", a, b, r1, r2))
}

func TestAcquireStringList(t *testing.T) {
	q := []string{"foo", "bar", "bgp"}
	w := []string{"foo", "bar", "bgp"}
	e := []string{"foo", "bpf"}

	p2 := NewStringListPool()
	x1 := p2.Acquire(q)
	x2 := p2.Acquire(w)
	x3 := p2.Acquire(e)
	fmt.Printf("Ptr: %p %p => %p %p\n", q, w, x1, x2)
	fmt.Printf("Ptr: %p => %p\n", e, x3)
}
