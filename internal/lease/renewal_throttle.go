package lease

import (
	"sync"
	"time"
)

// ThrottlePolicy defines rate-limiting rules for lease renewals.
type ThrottlePolicy struct {
	// MaxPerInterval is the maximum number of renewals allowed within the interval.
	MaxPerInterval int
	// Interval is the rolling window duration.
	Interval time.Duration
}

// DefaultThrottlePolicy returns a sensible default: 10 renewals per minute.
func DefaultThrottlePolicy() ThrottlePolicy {
	return ThrottlePolicy{
		MaxPerInterval: 10,
		Interval:       time.Minute,
	}
}

// RenewalThrottle tracks renewal attempts and enforces the ThrottlePolicy.
type RenewalThrottle struct {
	mu       sync.Mutex
	policy   ThrottlePolicy
	attempts []time.Time
	nowFn    func() time.Time
}

// NewRenewalThrottle creates a RenewalThrottle with the given policy.
func NewRenewalThrottle(policy ThrottlePolicy) *RenewalThrottle {
	return &RenewalThrottle{
		policy: policy,
		nowFn:  time.Now,
	}
}

// Allow returns true if a renewal attempt is permitted under the current policy,
// and records the attempt if so. It is safe for concurrent use.
func (t *RenewalThrottle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFn()
	cutoff := now.Add(-t.policy.Interval)

	// Evict attempts outside the rolling window.
	valid := t.attempts[:0]
	for _, ts := range t.attempts {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	t.attempts = valid

	if len(t.attempts) >= t.policy.MaxPerInterval {
		return false
	}

	t.attempts = append(t.attempts, now)
	return true
}

// Remaining returns the number of renewal slots still available in the current window.
func (t *RenewalThrottle) Remaining() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFn()
	cutoff := now.Add(-t.policy.Interval)

	count := 0
	for _, ts := range t.attempts {
		if ts.After(cutoff) {
			count++
		}
	}

	remaining := t.policy.MaxPerInterval - count
	if remaining < 0 {
		return 0
	}
	return remaining
}
