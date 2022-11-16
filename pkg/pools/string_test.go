package pools

import (
	"fmt"
	"testing"
)

func TestAcquireString(t *testing.T) {
	p := NewStringPool()
	s1 := p.Acquire("hello")
	s2 := p.Acquire("hello")
	s3 := p.Acquire("world")
	s1 = p.Acquire("hello")

	if s1 != s2 {
		t.Error("expected s1 == s2")
	}
	t.Log(fmt.Sprintf("s1, s2: %x %x", s1, s2))

	if s2 == s3 {
		t.Error("expected s2 !== s3")
	}
	t.Log(fmt.Sprintf("s1, s2: %x %x", s1, s2))
}

func TestGarbageCollectString(t *testing.T) {
	p := NewStringPool()

	// Gen 1
	p.Acquire("hello")
	p.Acquire("world")

	r := p.GarbageCollect()
	if r > 0 {
		t.Error("first run should not collect anything.")
	}

	p.Acquire("hello")
	p.Acquire("foo")
	r = p.GarbageCollect()
	if r != 1 {
		t.Error("expected 1 released value")
	}

	for k := range p.values {
		if k == "world" {
			t.Error("did not expect to find world here")
		}
	}
	t.Log(p.values)
	t.Log(p.counter)
}
