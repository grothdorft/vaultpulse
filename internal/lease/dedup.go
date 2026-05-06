package lease

import "time"

// ShouldAlert determines whether an alert should be sent for a lease entry.
// It returns true if the lease is expiring soon and has not already been alerted.
// cooldown prevents repeated alerts for the same lease within a given window.
func ShouldAlert(e *Entry, threshold, cooldown time.Duration) bool {
	if e == nil {
		return false
	}

	now := time.Now()

	// Already expired — no alert needed.
	if now.After(e.ExpiresAt) {
		return false
	}

	// Not within the alert threshold window yet.
	if e.ExpiresAt.Sub(now) > threshold {
		return false
	}

	// Already alerted and still within cooldown period.
	if e.State == StateAlerted && e.AlertedAt != nil {
		if now.Sub(*e.AlertedAt) < cooldown {
			return false
		}
	}

	return true
}
