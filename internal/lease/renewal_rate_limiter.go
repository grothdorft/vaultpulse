package lease

import (
	"sync"
	"time"
)

// RateLimiterPolicy defines the configuration for the token-bucket rate limiter.
type RateLimiterPolicy struct {
	// TokensPerInterval is the number of renewal tokens replenished each interval.
	TokensPerInterval int
	// Interval is the duration between token replenishments.
	Interval time.Duration
	// BurstCapacity is the maximum number of tokens that can accumulate.
	BurstCapacity int
}

// DefaultRateLimiterPolicy returns a sensible default rate limiter policy.
func DefaultRateLimiterPolicy() RateLimiterPolicy {
	return RateLimiterPolicy{
		TokensPerInterval: 10,
		Interval:          time.Minute,
		BurstCapacity:     20,
	}
}

// RenewalRateLimiter is a token-bucket rate limiter for lease renewals.
type RenewalRateLimiter struct {
	mu     sync.Mutex
	policy RateLimiterPolicy
	tokens int
	lastAt time.Time
	nowFn  func() time.Time
}

// NewRenewalRateLimiter creates a new RenewalRateLimiter with the given policy.
func NewRenewalRateLimiter(p RateLimiterPolicy) *RenewalRateLimiter {
	return &RenewalRateLimiter{
		policy: p,
		tokens: p.BurstCapacity,
		lastAt: time.Now(),
		nowFn:  time.Now,
	}
}

// Allow returns true if a renewal token is available and consumes one.
func (r *RenewalRateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.replenish()
	if r.tokens <= 0 {
		return false
	}
	r.tokens--
	return true
}

// Available returns the current number of available tokens.
func (r *RenewalRateLimiter) Available() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.replenish()
	return r.tokens
}

// Policy returns the current rate limiter policy.
func (r *RenewalRateLimiter) Policy() RateLimiterPolicy {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.policy
}

// replenish adds tokens based on elapsed time. Must be called with mu held.
func (r *RenewalRateLimiter) replenish() {
	now := r.nowFn()
	elapsed := now.Sub(r.lastAt)
	if elapsed < r.policy.Interval {
		return
	}
	intervals := int(elapsed / r.policy.Interval)
	added := intervals * r.policy.TokensPerInterval
	r.tokens += added
	if r.tokens > r.policy.BurstCapacity {
		r.tokens = r.policy.BurstCapacity
	}
	r.lastAt = r.lastAt.Add(time.Duration(intervals) * r.policy.Interval)
}
