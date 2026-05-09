package lease

import (
	"math"
	"time"
)

// BackoffPolicy defines parameters for exponential backoff on renewal failures.
type BackoffPolicy struct {
	// BaseDelay is the initial delay after the first failure.
	BaseDelay time.Duration
	// MaxDelay caps the computed delay regardless of failure count.
	MaxDelay time.Duration
	// Multiplier is the exponential growth factor (e.g. 2.0).
	Multiplier float64
}

// DefaultBackoffPolicy returns a sensible default backoff policy.
func DefaultBackoffPolicy() BackoffPolicy {
	return BackoffPolicy{
		BaseDelay:  30 * time.Second,
		MaxDelay:   10 * time.Minute,
		Multiplier: 2.0,
	}
}

// NextDelay computes the backoff duration for the given consecutive failure
// count (1-based). A failureCount of 0 or less always returns zero.
func (p BackoffPolicy) NextDelay(failureCount int) time.Duration {
	if failureCount <= 0 {
		return 0
	}
	// delay = base * multiplier^(failureCount-1)
	exponent := float64(failureCount - 1)
	delay := float64(p.BaseDelay) * math.Pow(p.Multiplier, exponent)
	if delay > float64(p.MaxDelay) {
		return p.MaxDelay
	}
	return time.Duration(delay)
}

// ShouldRetry returns true when the elapsed time since the last attempt
// exceeds the computed backoff delay for the given failure count.
func (p BackoffPolicy) ShouldRetry(failureCount int, lastAttempt time.Time) bool {
	if failureCount <= 0 {
		return true
	}
	delay := p.NextDelay(failureCount)
	return time.Since(lastAttempt) >= delay
}
