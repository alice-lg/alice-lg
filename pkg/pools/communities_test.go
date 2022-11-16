package pools

import (
	"fmt"
	"testing"

	"github.com/alice-lg/alice-lg/pkg/api"
)

func TestAcquireCommunities(t *testing.T) {
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

	p := NewCommunities()

	pc1 := p.Acquire(c1)
	pc2 := p.Acquire(c2)
	pc3 := p.Acquire(c3)

	if fmt.Sprintf("%p", c1) == fmt.Sprintf("%p", c2) {
		t.Error("expected c1 !== c2")
	}

	if fmt.Sprintf("%p", pc1) != fmt.Sprintf("%p", pc2) {
		t.Error("expected pc1 == pc2")
	}

	fmt.Printf("c1:  %p, c2:  %p, c3:  %p\n", c1, c2, c3)
	fmt.Printf("pc1: %p, pc2: %p, pc3: %p\n", pc1, pc2, pc3)

}
