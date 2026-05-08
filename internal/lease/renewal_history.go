package lease

import (
	"sync"
	"time"
)

// RenewalRecord captures a single renewal event for a lease.
type RenewalRecord struct {
	LeaseID   string
	RenewedAt time.Time
	NewExpiry time.Time
	Success   bool
	Error     string
}

// RenewalHistory tracks renewal attempts per lease, capped to a maximum
// number of records per lease to prevent unbounded growth.
type RenewalHistory struct {
	mu      sync.RWMutex
	records map[string][]RenewalRecord
	maxPerLease int
}

// NewRenewalHistory creates a RenewalHistory that retains up to maxPerLease
// records per lease ID. If maxPerLease is <= 0 it defaults to 10.
func NewRenewalHistory(maxPerLease int) *RenewalHistory {
	if maxPerLease <= 0 {
		maxPerLease = 10
	}
	return &RenewalHistory{
		records:     make(map[string][]RenewalRecord),
		maxPerLease: maxPerLease,
	}
}

// Record appends a RenewalRecord for the given lease, evicting the oldest
// entry when the cap is reached.
func (h *RenewalHistory) Record(r RenewalRecord) {
	h.mu.Lock()
	defer h.mu.Unlock()

	list := h.records[r.LeaseID]
	list = append(list, r)
	if len(list) > h.maxPerLease {
		list = list[len(list)-h.maxPerLease:]
	}
	h.records[r.LeaseID] = list
}

// Get returns a copy of all renewal records for a lease.
func (h *RenewalHistory) Get(leaseID string) []RenewalRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()

	src := h.records[leaseID]
	out := make([]RenewalRecord, len(src))
	copy(out, src)
	return out
}

// Remove deletes all renewal history for a lease.
func (h *RenewalHistory) Remove(leaseID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.records, leaseID)
}

// Len returns the total number of records stored across all leases.
func (h *RenewalHistory) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, v := range h.records {
		total += len(v)
	}
	return total
}
