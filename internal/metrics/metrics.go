package metrics

import "sync/atomic"

// Counters holds runtime statistics for VaultPulse.
type Counters struct {
	AlertsTotal    atomic.Int64
	AlertErrors    atomic.Int64
	LeasesTracked  atomic.Int64
	LeasesExpired  atomic.Int64
	ChecksRun      atomic.Int64
}

// Global is the process-wide metrics instance.
var Global = &Counters{}

// Snapshot returns a point-in-time copy of all counters.
type Snapshot struct {
	AlertsTotal   int64 `json:"alerts_total"`
	AlertErrors   int64 `json:"alert_errors"`
	LeasesTracked int64 `json:"leases_tracked"`
	LeasesExpired int64 `json:"leases_expired"`
	ChecksRun     int64 `json:"checks_run"`
}

// Snapshot returns a consistent read of all counters.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		AlertsTotal:   c.AlertsTotal.Load(),
		AlertErrors:   c.AlertErrors.Load(),
		LeasesTracked: c.LeasesTracked.Load(),
		LeasesExpired: c.LeasesExpired.Load(),
		ChecksRun:     c.ChecksRun.Load(),
	}
}

// Reset zeroes all counters (useful in tests).
func (c *Counters) Reset() {
	c.AlertsTotal.Store(0)
	c.AlertErrors.Store(0)
	c.LeasesTracked.Store(0)
	c.LeasesExpired.Store(0)
	c.ChecksRun.Store(0)
}
