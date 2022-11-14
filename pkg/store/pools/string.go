package pools

import "sync"

// String is a pool for strings.
// This will most likely be a pool for IP addresses.
type String struct {
	values map[string]*string

	counter map[string]uint
	top     uint

	sync.Mutex
}

// NewString creates a new string pool
func NewString() *String {
	return &String{
		values:  map[string]*string{},
		counter: map[string]uint{},
	}
}

// Acquire a pointer to a string value
func (p *String) Acquire(s string) *string {
	p.Lock()
	defer p.Unlock()
	// Deduplicate value
	ptr, ok := p.values[s]
	if !ok {
		p.values[s] = &s
		ptr = &s
	}
	p.counter[s] = p.top
	return ptr
}

// GarbageCollect releases all values, which have not been seen
// again.
func (p *String) GarbageCollect() uint {
	p.Lock()
	defer p.Unlock()
	var released uint = 0
	for k, cnt := range p.counter {
		if cnt < p.top {
			delete(p.counter, k)
			delete(p.values, k)
			released++
		}
	}
	p.top++ // Next generation
	return released
}
