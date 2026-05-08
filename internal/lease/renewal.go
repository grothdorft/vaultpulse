package lease

import (
	"time"
)

// RenewalAdvice describes whether a lease should be renewed and why.
type RenewalAdvice struct {
	LeaseID     string
	ShouldRenew bool
	Reason      string
	Urgency     ExpiryStatus
}

// RenewalPolicy holds thresholds used to decide when to recommend renewal.
type RenewalPolicy struct {
	// RenewBeforeTTL is how far before expiry we recommend renewal.
	RenewBeforeTTL time.Duration
	// MinRemainingAfterRenewal is the minimum TTL we expect after renewal;
	// if the lease is already beyond this, skip.
	MinRemainingAfterRenewal time.Duration
}

// DefaultRenewalPolicy returns a sensible default policy.
func DefaultRenewalPolicy() RenewalPolicy {
	return RenewalPolicy{
		RenewBeforeTTL:           10 * time.Minute,
		MinRemainingAfterRenewal: 30 * time.Second,
	}
}

// AdviseRenewal evaluates a tracker entry against the given policy and
// returns a RenewalAdvice indicating whether the lease should be renewed.
func AdviseRenewal(entry *Entry, policy RenewalPolicy, now time.Time) RenewalAdvice {
	if entry == nil {
		return RenewalAdvice{
			ShouldRenew: false,
			Reason:      "nil entry",
		}
	}

	if entry.Expired {
		return RenewalAdvice{
			LeaseID:     entry.LeaseID,
			ShouldRenew: false,
			Reason:      "lease already expired",
			Urgency:     StatusExpired,
		}
	}

	remaining := entry.ExpiresAt.Sub(now)
	if remaining <= policy.MinRemainingAfterRenewal {
		return RenewalAdvice{
			LeaseID:     entry.LeaseID,
			ShouldRenew: false,
			Reason:      "insufficient time remaining to justify renewal",
			Urgency:     StatusExpired,
		}
	}

	status := ClassifyExpiry(entry.ExpiresAt, now)

	if remaining <= policy.RenewBeforeTTL {
		return RenewalAdvice{
			LeaseID:     entry.LeaseID,
			ShouldRenew: true,
			Reason:      "lease expiring within renewal window",
			Urgency:     status,
		}
	}

	return RenewalAdvice{
		LeaseID:     entry.LeaseID,
		ShouldRenew: false,
		Reason:      "lease has sufficient TTL remaining",
		Urgency:     status,
	}
}
