package lease

import (
	"sync"
	"time"
)

// AuditEventKind classifies what happened during a renewal attempt.
type AuditEventKind string

const (
	AuditKindSuccess    AuditEventKind = "success"
	AuditKindFailure    AuditEventKind = "failure"
	AuditKindSkipped    AuditEventKind = "skipped"
	AuditKindThrottled  AuditEventKind = "throttled"
	AuditKindCircuitOpen AuditEventKind = "circuit_open"
)

// AuditEvent records a single renewal lifecycle event for a lease.
type AuditEvent struct {
	LeaseID   string         `json:"lease_id"`
	Kind      AuditEventKind `json:"kind"`
	Message   string         `json:"message,omitempty"`
	OccurredAt time.Time     `json:"occurred_at"`
}

// RenewalAuditLog is a bounded, thread-safe log of renewal audit events.
type RenewalAuditLog struct {
	mu     sync.RWMutex
	events []AuditEvent
	cap    int
}

// NewRenewalAuditLog creates an audit log that retains at most cap events.
func NewRenewalAuditLog(cap int) *RenewalAuditLog {
	if cap <= 0 {
		cap = 256
	}
	return &RenewalAuditLog{
		events: make([]AuditEvent, 0, cap),
		cap:    cap,
	}
}

// Record appends an audit event, evicting the oldest when the log is full.
func (a *RenewalAuditLog) Record(leaseID string, kind AuditEventKind, message string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.events) >= a.cap {
		a.events = a.events[1:]
	}
	a.events = append(a.events, AuditEvent{
		LeaseID:    leaseID,
		Kind:       kind,
		Message:    message,
		OccurredAt: time.Now().UTC(),
	})
}

// All returns a snapshot of all stored audit events, oldest first.
func (a *RenewalAuditLog) All() []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]AuditEvent, len(a.events))
	copy(out, a.events)
	return out
}

// ForLease returns all events recorded for the given lease ID.
func (a *RenewalAuditLog) ForLease(leaseID string) []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var out []AuditEvent
	for _, e := range a.events {
		if e.LeaseID == leaseID {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the current number of stored events.
func (a *RenewalAuditLog) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.events)
}
