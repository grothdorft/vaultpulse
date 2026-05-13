package lease

import (
	"sync"
	"time"
)

// WindowPolicy defines the configuration for a renewal window.
type WindowPolicy struct {
	// MinWindow is the minimum allowed renewal window duration.
	MinWindow time.Duration
	// MaxWindow is the maximum allowed renewal window duration.
	MaxWindow time.Duration
	// DefaultFraction is the fraction of the lease TTL used to compute the window.
	DefaultFraction float64
}

// DefaultWindowPolicy returns a sensible default WindowPolicy.
func DefaultWindowPolicy() WindowPolicy {
	return WindowPolicy{
		MinWindow:       30 * time.Second,
		MaxWindow:       24 * time.Hour,
		DefaultFraction: 0.25,
	}
}

// RenewalWindow tracks per-lease computed renewal windows.
type RenewalWindow struct {
	mu      sync.RWMutex
	policy  WindowPolicy
	windows map[string]time.Duration
}

// NewRenewalWindow creates a new RenewalWindow with the given policy.
func NewRenewalWindow(p WindowPolicy) *RenewalWindow {
	return &RenewalWindow{
		policy:  p,
		windows: make(map[string]time.Duration),
	}
}

// Compute derives the renewal window for a given TTL and stores it by leaseID.
func (rw *RenewalWindow) Compute(leaseID string, ttl time.Duration) time.Duration {
	w := time.Duration(float64(ttl) * rw.policy.DefaultFraction)
	if w < rw.policy.MinWindow {
		w = rw.policy.MinWindow
	}
	if w > rw.policy.MaxWindow {
		w = rw.policy.MaxWindow
	}
	rw.mu.Lock()
	rw.windows[leaseID] = w
	rw.mu.Unlock()
	return w
}

// Get returns the stored renewal window for a leaseID, and whether it exists.
func (rw *RenewalWindow) Get(leaseID string) (time.Duration, bool) {
	rw.mu.RLock()
	defer rw.mu.RUnlock()
	w, ok := rw.windows[leaseID]
	return w, ok
}

// Remove deletes the stored window for a leaseID.
func (rw *RenewalWindow) Remove(leaseID string) {
	rw.mu.Lock()
	delete(rw.windows, leaseID)
	rw.mu.Unlock()
}

// Snapshot returns a copy of all stored windows.
func (rw *RenewalWindow) Snapshot() map[string]time.Duration {
	rw.mu.RLock()
	defer rw.mu.RUnlock()
	out := make(map[string]time.Duration, len(rw.windows))
	for k, v := range rw.windows {
		out[k] = v
	}
	return out
}
