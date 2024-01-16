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

	r1, gid1 := p.AcquireGid(a)
	p.Acquire(c)
	r2, gid2 := p.AcquireGid(b)

	log.Println("r1", r1, "gid1", gid1)
	log.Println("r2", r2, "gid2", gid2)

	if fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b) {
		t.Error("lists should not be same pointer", fmt.Sprintf("%p %p", a, b))
	}
	if fmt.Sprintf("%p", r1) != fmt.Sprintf("%p", r2) {
		t.Error("lists should be same pointer", fmt.Sprintf("%p %p", r1, r2))
	}
	if gid1 != gid2 {
		t.Error("gid should be same, got:", gid1, gid2)
	}

	t.Log(fmt.Sprintf("Ptr: %p %p => %p %p", a, b, r1, r2))

	_, gid3 := p.AcquireGid(c)
	if gid3 == gid1 {
		t.Error("gid should not be same, got:", gid3, gid1)
	}
	t.Log("gid3", gid3, "gid1", gid1)
}

func TestAcquireStringList(t *testing.T) {
	q := []string{"foo", "bar", "bgp"}
	w := []string{"foo", "bar", "bgp"}
	e := []string{"foo", "bpf"}

	p2 := NewStringListPool()
	x1, g1 := p2.AcquireGid(q)
	x2, g2 := p2.AcquireGid(w)
	x3, g3 := p2.AcquireGid(e)
	fmt.Printf("Ptr: %p %p => %p %d %p %d \n", q, w, x1, g1, x2, g2)
	fmt.Printf("Ptr: %p => %p %d\n", e, x3, g3)
}
