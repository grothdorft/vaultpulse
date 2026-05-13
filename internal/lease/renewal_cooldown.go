package lease

import (
	"sync"
	"time"
)

// CooldownPolicy defines how long to wait between renewal attempts for a lease.
type CooldownPolicy struct {
	// MinInterval is the minimum time that must elapse between two renewals
	// of the same lease ID.
	MinInterval time.Duration
}

// DefaultCooldownPolicy returns a sensible production default.
func DefaultCooldownPolicy() CooldownPolicy {
	return CooldownPolicy{
		MinInterval: 30 * time.Second,
	}
}

// RenewalCooldown tracks the last renewal time per lease and enforces a
// minimum interval between successive renewals.
type RenewalCooldown struct {
	mu     sync.Mutex
	policy CooldownPolicy
	last   map[string]time.Time
	now    func() time.Time
}

// NewRenewalCooldown creates a RenewalCooldown using the given policy.
func NewRenewalCooldown(policy CooldownPolicy) *RenewalCooldown {
	return &RenewalCooldown{
		policy: policy,
		last:   make(map[string]time.Time),
		now:    time.Now,
	}
}

// Allow returns true when the lease is not in its cooldown window, meaning a
// renewal attempt is permitted. Calling Allow records the current time as the
// last renewal instant for the lease.
func (c *RenewalCooldown) Allow(leaseID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	if t, ok := c.last[leaseID]; ok {
		if now.Sub(t) < c.policy.MinInterval {
			return false
		}
	}
	c.last[leaseID] = now
	return true
}

// Reset removes the cooldown record for a specific lease, allowing an
// immediate renewal on the next Allow call.
func (c *RenewalCooldown) Reset(leaseID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.last, leaseID)
}

// Remaining returns how much cooldown time is left for a lease.
// Returns 0 if the lease is not in cooldown.
func (c *RenewalCooldown) Remaining(leaseID string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	t, ok := c.last[leaseID]
	if !ok {
		return 0
	}
	elapsed := c.now().Sub(t)
	if elapsed >= c.policy.MinInterval {
		return 0
	}
	return c.policy.MinInterval - elapsed
}
