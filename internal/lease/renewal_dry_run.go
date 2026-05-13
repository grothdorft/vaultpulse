package lease

import (
	"sync"
	"time"
)

// DryRunRecord captures what would have happened during a renewal in dry-run mode.
type DryRunRecord struct {
	LeaseID   string
	At        time.Time
	Advice    RenewalAdvice
	Reason    string
}

// DryRunLog stores simulated renewal actions without executing them.
type DryRunLog struct {
	mu      sync.Mutex
	records []DryRunRecord
	cap     int
}

// NewDryRunLog creates a DryRunLog with the given capacity.
func NewDryRunLog(cap int) *DryRunLog {
	if cap <= 0 {
		cap = 100
	}
	return &DryRunLog{cap: cap}
}

// Record appends a simulated renewal action.
func (d *DryRunLog) Record(leaseID string, advice RenewalAdvice, reason string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.records) >= d.cap {
		d.records = d.records[1:]
	}
	d.records = append(d.records, DryRunRecord{
		LeaseID: leaseID,
		At:      time.Now().UTC(),
		Advice:  advice,
		Reason:  reason,
	})
}

// All returns a copy of all recorded dry-run entries.
func (d *DryRunLog) All() []DryRunRecord {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]DryRunRecord, len(d.records))
	copy(out, d.records)
	return out
}

// Len returns the number of recorded entries.
func (d *DryRunLog) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.records)
}

// Clear removes all entries from the log.
func (d *DryRunLog) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.records = d.records[:0]
}
