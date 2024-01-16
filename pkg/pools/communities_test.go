package pools

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestAcquireGidCommunity(t *testing.T) {
	c1 := api.Community{2342, 5, 1}
	c2 := api.Community{2342, 5, 1}
	c3 := api.Community{2342, 5}

	p := NewCommunitiesPool()

	pc1, gid1 := p.AcquireGid(c1)
	pc2, gid2 := p.AcquireGid(c2)
	pc3, gid3 := p.AcquireGid(c3)

	if fmt.Sprintf("%p", c1) == fmt.Sprintf("%p", c2) {
		t.Error("expected c1 !== c2")
	}

	if fmt.Sprintf("%p", pc1) != fmt.Sprintf("%p", pc2) {
		t.Error("expected pc1 == pc2")
	}

	if gid1 != gid2 {
		t.Error("expected gid1 == gid2")
	}

	if gid1 == gid3 {
		t.Error("expected gid1 != gid3")
	}

	fmt.Printf("c1:  %p, c2:  %p, c3:  %p\n", c1, c2, c3)
	fmt.Printf("pc1: %p, pc2: %p, pc3: %p\n", pc1, pc2, pc3)

	log.Println(c3, pc3)
	log.Println(gid1, gid2, gid3)
}

func TestCommunityRead(t *testing.T) {
	c1 := api.Community{1111, 5, 1}
	c2 := api.Community{1111, 5, 1}
	c3 := api.Community{1111, 5}

	p := NewCommunitiesPool()

	pc1, gid1 := p.AcquireGid(c1)
	pc2, gid2 := p.Read(c2)
	pc3, gid3 := p.Read(c3)

	fmt.Printf("pc1: %p, pc2: %p, pc3: %p\n", pc1, pc2, pc3)
	fmt.Printf("gid1: %d, gid2: %d, gid3: %d\n", gid1, gid2, gid3)

	if fmt.Sprintf("%p", pc1) != fmt.Sprintf("%p", pc2) {
		t.Error("expected pc1 == pc2")
	}
	if gid1 != gid2 {
		t.Error("expected gid1 == gid2")
	}

	if pc3 != nil {
		t.Error("expected pc3 == nil, got", pc3)
	}
	if gid3 != 0 {
		t.Error("expected gid3 == 0, got", gid3)
	}
}

func TestAcquireGidCommunitiesSets(t *testing.T) {
	c1 := []api.Community{
		{2342, 5, 1},
		{2342, 5, 2},
		{2342, 51, 1},
	}
	c2 := []api.Community{
		{2342, 5, 1},
		{2342, 5, 2},
		{2342, 51, 1},
	}
	c3 := []api.Community{
		{2341, 6, 1},
		{2341, 6, 2},
		{2341, 1, 1},
	}

	p := NewCommunitiesSetPool()

	pc1, gid1 := p.AcquireGid(c1)
	pc2, gid2 := p.AcquireGid(c2)
	pc3, gid3 := p.AcquireGid(c3)

	if fmt.Sprintf("%p", c1) == fmt.Sprintf("%p", c2) {
		t.Error("expected c1 !== c2")
	}

	if fmt.Sprintf("%p", pc1) != fmt.Sprintf("%p", pc2) {
		t.Error("expected pc1 == pc2")
	}
	if gid1 != gid2 {
		t.Error("expected gid1 == gid2")
	}

	fmt.Printf("c1:  %p, c2:  %p, c3:  %p\n", c1, c2, c3)
	fmt.Printf("pc1: %p, pc2: %p, pc3: %p\n", pc1, pc2, pc3)

	t.Logf("gid1: %d, gid2: %d, gid3: %d\n", gid1, gid2, gid3)
}

func TestSetCommunityIdentity(t *testing.T) {
	set := []api.Community{
		{2341, 6, 1},
		{2341, 6, 2},
		{2341, 1, 1},
	}

	pset, gid1 := CommunitiesSets.AcquireGid(set)
	pval, gid2 := Communities.AcquireGid(api.Community{2341, 6, 2})

	fmt.Printf("set:  %p, pset[1]:  %p, pval:  %p\n", set, pset[1], pval)
	fmt.Printf("gid1: %d, gid2: %d\n", gid1, gid2)

	p1 := reflect.ValueOf(pset[1]).UnsafePointer()
	p2 := reflect.ValueOf(pval).UnsafePointer()

	if p1 != p2 {
		t.Error("expected pset[1] == pval")
	}
}

func TestAcquireGidExtCommunitiesSets(t *testing.T) {
	c1 := []api.ExtCommunity{
		{"ro", 5, 1},
		{"ro", 5, 2},
		{"rt", 51, 1},
	}
	c2 := []api.ExtCommunity{
		{"ro", 5, 1},
		{"ro", 5, 2},
		{"rt", 51, 1},
	}
	c3 := []api.ExtCommunity{
		{"ro", 6, 1},
		{"rt", 6, 2},
		{"xyz", 1, 1},
	}

	p := NewExtCommunitiesSetPool()

	pc1, gid1 := p.AcquireGid(c1)
	pc2, gid2 := p.AcquireGid(c2)
	pc3, gid3 := p.AcquireGid(c3)

	if fmt.Sprintf("%p", c1) == fmt.Sprintf("%p", c2) {
		t.Error("expected c1 !== c2")
	}

	if fmt.Sprintf("%p", pc1) != fmt.Sprintf("%p", pc2) {
		t.Error("expected pc1 == pc2")
	}
	if gid1 != gid2 {
		t.Error("expected gid1 == gid2")
	}

	fmt.Printf("c1:  %p, c2:  %p, c3:  %p\n", c1, c2, c3)
	fmt.Printf("pc1: %p, pc2: %p, pc3: %p\n", pc1, pc2, pc3)
	fmt.Printf("gid1: %d, gid2: %d, gid3: %d\n", gid1, gid2, gid3)
}
