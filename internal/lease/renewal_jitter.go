package lease

import (
	"math/rand"
	"time"
)

// JitterPolicy defines parameters for adding randomised jitter to renewal delays.
type JitterPolicy struct {
	// MaxFraction is the maximum fraction of the base delay to add as jitter.
	// For example, 0.25 means up to 25% of the base delay is added randomly.
	MaxFraction float64
}

// DefaultJitterPolicy returns a JitterPolicy suitable for most deployments.
func DefaultJitterPolicy() JitterPolicy {
	return JitterPolicy{
		MaxFraction: 0.20,
	}
}

// Apply returns the base delay with a random jitter applied according to the policy.
// The jitter is always non-negative, so the returned duration is >= base.
func (p JitterPolicy) Apply(base time.Duration, rng *rand.Rand) time.Duration {
	if base <= 0 || p.MaxFraction <= 0 {
		return base
	}
	maxJitter := float64(base) * p.MaxFraction
	jitter := time.Duration(rng.Float64() * maxJitter)
	return base + jitter
}

// ApplyGlobal is a convenience wrapper that uses the global random source.
func (p JitterPolicy) ApplyGlobal(base time.Duration) time.Duration {
	//nolint:gosec // non-cryptographic jitter is intentional
	return p.Apply(base, rand.New(rand.NewSource(time.Now().UnixNano())))
}
