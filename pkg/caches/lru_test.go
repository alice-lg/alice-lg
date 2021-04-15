package caches

import (
	"testing"
	"time"
)

func TestLRUMap(t *testing.T) {

	accessedAt := LRUMap{}

	accessedAt["foo"] = time.Now()
	accessedAt["bar"] = time.Now().Add(-2 * time.Minute)
	accessedAt["bam"] = time.Now().Add(-1 * time.Minute)

	lru := accessedAt.LRU()
	if lru != "bar" {
		t.Error("Expected bar to be LRU, got:", lru)
	}
}
