package lease

import (
	"sync"
	"time"
)

// State represents the alert state of a lease.
type State int

const (
	StateUnknown State = iota
	StateHealthy
	StateAlerted
	StateExpired
)

// Entry tracks a single lease and its alert state.
type Entry struct {
	LeaseID   string
	ExpiresAt time.Time
	State     State
	AlertedAt *time.Time
}

// Tracker maintains in-memory state for monitored leases.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// NewTracker creates a new Tracker instance.
func NewTracker() *Tracker {
	return &Tracker{
		entries: make(map[string]*Entry),
	}
}

// Upsert inserts or updates a lease entry.
func (t *Tracker) Upsert(leaseID string, expiresAt time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if existing, ok := t.entries[leaseID]; ok {
		existing.ExpiresAt = expiresAt
		return
	}
	t.entries[leaseID] = &Entry{
		LeaseID:   leaseID,
		ExpiresAt: expiresAt,
		State:     StateHealthy,
	}
}

// MarkAlerted records that an alert was sent for the given lease.
func (t *Tracker) MarkAlerted(leaseID string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if e, ok := t.entries[leaseID]; ok {
		now := time.Now()
		e.State = StateAlerted
		e.AlertedAt = &now
	}
}

// MarkExpired marks a lease as expired and removes it from active tracking.
func (t *Tracker) MarkExpired(leaseID string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if e, ok := t.entries[leaseID]; ok {
		e.State = StateExpired
	}
}

// Get returns the entry for a lease ID, or nil if not found.
func (t *Tracker) Get(leaseID string) *Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.entries[leaseID]
}

// Remove deletes a lease entry from the tracker.
func (t *Tracker) Remove(leaseID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, leaseID)
}

// All returns a snapshot of all current entries.
func (t *Tracker) All() []*Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]*Entry, 0, len(t.entries))
	for _, e := range t.entries {
		copy := *e
		result = append(result, &copy)
	}
	return result
}
