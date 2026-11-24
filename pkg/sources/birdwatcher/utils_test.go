package birdwatcher

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestIsProtocolUp(t *testing.T) {
	tests := map[string]bool{
		"up":   true,
		"uP":   true,
		"Up":   true,
		"UP":   true,
		"down": false,
	}

	for up, expected := range tests {
		if isProtocolUp(up) != expected {
			t.Error("f(", up, ") != ", expected)
		}
	}
}

func TestLockMapDropsReleasedKeys(t *testing.T) {
	m := NewLockMap()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("n%d", i)
		m.Lock(key)
		m.Unlock(key)
	}

	if len(m.locks) != 0 {
		t.Errorf("locks after release: got %d, want 0", len(m.locks))
	}
}

func TestLockMapExcludesPerKey(t *testing.T) {
	m := NewLockMap()
	var holders int32
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.Lock("n1")
				if atomic.AddInt32(&holders, 1) != 1 {
					t.Error("two callers hold the lock for n1")
				}
				atomic.AddInt32(&holders, -1)
				m.Unlock("n1")
			}
		}()
	}
	wg.Wait()

	if len(m.locks) != 0 {
		t.Errorf("locks after release: got %d, want 0", len(m.locks))
	}
}
