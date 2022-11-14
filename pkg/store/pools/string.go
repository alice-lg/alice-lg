package pools

import "sync"

// String is a pool for strings.
type String struct {
	values map[string]*string
	counter map[string]uint
	top uint

	sync.Mutex
}

// NewString creates a new string pool
func NewString() *String {
	return &String{
		values: map[string]*string{},
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
	// Increment counter and top value
	cnt, ok := p.counter[s]
	if !ok {
		cnt = p.top + 1
	} else {
		cnt = cnt + 1
	}
	p.counter[s] = cnt
	if p.top < cnt {
		p.top = cnt
	}
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
	return released
}
