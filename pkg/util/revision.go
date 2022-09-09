package util

import (
	"sync"
	"time"
)

// Revision contains the current interval loop revision.
type Revision struct {
	mu  sync.Mutex
	rev string
}

// Next generates and retrieves a new revision.
func (r *Revision) Next() string {
	r.mu.Lock()
	r.rev = time.Now().Format("20060102.150405.000")
	rev := r.rev
	r.mu.Unlock()
	return rev
}

// Current retrieves the current revision.
func (r *Revision) Current() string {
	r.mu.Lock()
	rev := r.rev
	r.mu.Unlock()
	return rev
}
