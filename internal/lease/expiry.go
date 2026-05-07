package lease

import "time"

// ExpiryLevel categorises how urgently a lease needs attention.
type ExpiryLevel int

const (
	ExpiryLevelOK      ExpiryLevel = iota // more than warnThreshold remaining
	ExpiryLevelWarning                    // within warnThreshold
	ExpiryLevelCritical                   // within criticalThreshold
	ExpiryLevelExpired                    // already past expiry
)

// String returns a human-readable label for the level.
func (l ExpiryLevel) String() string {
	switch l {
	case ExpiryLevelWarning:
		return "warning"
	case ExpiryLevelCritical:
		return "critical"
	case ExpiryLevelExpired:
		return "expired"
	default:
		return "ok"
	}
}

// ClassifyExpiry returns the ExpiryLevel for a lease expiry time given the
// configured warning and critical thresholds.
//
// warnThreshold    – alert when TTL falls below this value (e.g. 24h)
// criticalThreshold – alert at a higher urgency below this value (e.g. 1h)
func ClassifyExpiry(expiry time.Time, warnThreshold, criticalThreshold time.Duration) ExpiryLevel {
	now := time.Now()
	if !now.Before(expiry) {
		return ExpiryLevelExpired
	}
	ttl := expiry.Sub(now)
	switch {
	case ttl <= criticalThreshold:
		return ExpiryLevelCritical
	case ttl <= warnThreshold:
		return ExpiryLevelWarning
	default:
		return ExpiryLevelOK
	}
}
