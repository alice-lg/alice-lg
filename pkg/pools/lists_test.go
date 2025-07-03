package pools

import (
	"fmt"
	"log"
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

	log.Println("r1", r1)
	log.Println("r2", r2)

	if fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b) {
		t.Errorf("lists should not be same pointer: %p %p", a, b)
	}
	if fmt.Sprintf("%p", r1) != fmt.Sprintf("%p", r2) {
		t.Errorf("lists should be same pointer: %p %p", r1, r2)
	}

	t.Logf("Ptr: %p %p => %p %p", a, b, r1, r2)
}

func TestAcquireStringList(t *testing.T) {
	q := []string{"foo", "bar", "bgp"}
	w := []string{"foo", "bar", "bgp"}
	e := []string{"foo", "bpf"}

	p2 := NewStringListPool()
	x1 := p2.Acquire(q)
	p2.Acquire(e)
	x2 := p2.Acquire(w)
	t.Logf("Ptr: %p %p => %p %p", q, w, x1, x2)
}
