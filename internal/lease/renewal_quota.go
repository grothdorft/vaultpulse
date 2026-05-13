package lease

import (
	"fmt"
	"sync"
	"time"
)

// QuotaPolicy defines per-namespace or per-mount renewal quotas.
type QuotaPolicy struct {
	// MaxRenewalsPerWindow is the maximum renewals allowed in the window.
	MaxRenewalsPerWindow int
	// Window is the duration of the quota window.
	Window time.Duration
}

// DefaultQuotaPolicy returns a sensible default quota policy.
func DefaultQuotaPolicy() QuotaPolicy {
	return QuotaPolicy{
		MaxRenewalsPerWindow: 50,
		Window:               time.Minute,
	}
}

// QuotaEntry tracks usage for a single namespace/mount key.
type QuotaEntry struct {
	Count     int
	WindowEnd time.Time
}

// RenewalQuota enforces per-key renewal quotas.
type RenewalQuota struct {
	mu     sync.Mutex
	policy QuotaPolicy
	usage  map[string]*QuotaEntry
}

// NewRenewalQuota creates a new RenewalQuota with the given policy.
func NewRenewalQuota(policy QuotaPolicy) *RenewalQuota {
	return &RenewalQuota{
		policy: policy,
		usage:  make(map[string]*QuotaEntry),
	}
}

// Allow returns true if the key is within quota and increments the counter.
// Returns an error describing the quota violation if denied.
func (q *RenewalQuota) Allow(key string) (bool, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	entry, ok := q.usage[key]
	if !ok || now.After(entry.WindowEnd) {
		q.usage[key] = &QuotaEntry{
			Count:     1,
			WindowEnd: now.Add(q.policy.Window),
		}
		return true, nil
	}

	if entry.Count >= q.policy.MaxRenewalsPerWindow {
		return false, fmt.Errorf("quota exceeded for %q: %d/%d renewals in window ending %s",
			key, entry.Count, q.policy.MaxRenewalsPerWindow, entry.WindowEnd.Format(time.RFC3339))
	}

	entry.Count++
	return true, nil
}

// Snapshot returns a copy of current usage keyed by namespace/mount.
func (q *RenewalQuota) Snapshot() map[string]QuotaEntry {
	q.mu.Lock()
	defer q.mu.Unlock()

	out := make(map[string]QuotaEntry, len(q.usage))
	for k, v := range q.usage {
		out[k] = *v
	}
	return out
}

// Policy returns the current quota policy.
func (q *RenewalQuota) Policy() QuotaPolicy {
	return q.policy
}
