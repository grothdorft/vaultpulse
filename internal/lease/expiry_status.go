package lease

import "time"

// ExpiryStatus represents the current state of a lease's expiry classification.
type ExpiryStatus string

const (
	StatusOK       ExpiryStatus = "ok"
	StatusWarning  ExpiryStatus = "warning"
	StatusCritical ExpiryStatus = "critical"
	StatusExpired  ExpiryStatus = "expired"
)

// StatusSummary holds a snapshot of a lease's expiry information.
type StatusSummary struct {
	LeaseID   string
	Status    ExpiryStatus
	ExpiresAt time.Time
	TTL       time.Duration
}

// BuildStatusSummary constructs a StatusSummary for a given lease entry
// using the provided warning and critical thresholds.
func BuildStatusSummary(entry *Entry, warnThreshold, critThreshold time.Duration) StatusSummary {
	if entry == nil {
		return StatusSummary{Status: StatusExpired}
	}

	now := time.Now()
	ttl := entry.ExpiresAt.Sub(now)

	classification := ClassifyExpiry(entry.ExpiresAt, warnThreshold, critThreshold)

	var status ExpiryStatus
	switch classification {
	case ExpiryOK:
		status = StatusOK
	case ExpiryWarning:
		status = StatusWarning
	case ExpiryCritical:
		status = StatusCritical
	default:
		status = StatusExpired
	}

	return StatusSummary{
		LeaseID:   entry.LeaseID,
		Status:    status,
		ExpiresAt: entry.ExpiresAt,
		TTL:       ttl,
	}
}
