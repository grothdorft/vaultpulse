package lease

import (
	"testing"
	"time"
)

func makeRecord(leaseID string, success bool) RenewalRecord {
	now := time.Now()
	return RenewalRecord{
		LeaseID:   leaseID,
		RenewedAt: now,
		NewExpiry: now.Add(24 * time.Hour),
		Success:   success,
	}
}

func TestRenewalHistory_RecordAndGet(t *testing.T) {
	h := NewRenewalHistory(5)

	h.Record(makeRecord("lease-1", true))
	h.Record(makeRecord("lease-1", false))

	records := h.Get("lease-1")
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Success != true {
		t.Errorf("expected first record Success=true")
	}
	if records[1].Success != false {
		t.Errorf("expected second record Success=false")
	}
}

func TestRenewalHistory_CapEvictsOldest(t *testing.T) {
	h := NewRenewalHistory(3)

	for i := 0; i < 5; i++ {
		h.Record(makeRecord("lease-cap", i%2 == 0))
	}

	records := h.Get("lease-cap")
	if len(records) != 3 {
		t.Fatalf("expected cap of 3, got %d", len(records))
	}
}

func TestRenewalHistory_GetUnknownLeaseReturnsEmpty(t *testing.T) {
	h := NewRenewalHistory(5)
	records := h.Get("nonexistent")
	if len(records) != 0 {
		t.Errorf("expected empty slice for unknown lease, got %d", len(records))
	}
}

func TestRenewalHistory_Remove(t *testing.T) {
	h := NewRenewalHistory(5)
	h.Record(makeRecord("lease-del", true))
	h.Remove("lease-del")

	if records := h.Get("lease-del"); len(records) != 0 {
		t.Errorf("expected empty after remove, got %d record(s)", len(records))
	}
}

func TestRenewalHistory_Len(t *testing.T) {
	h := NewRenewalHistory(10)
	h.Record(makeRecord("a", true))
	h.Record(makeRecord("a", false))
	h.Record(makeRecord("b", true))

	if n := h.Len(); n != 3 {
		t.Errorf("expected Len=3, got %d", n)
	}
}

func TestRenewalHistory_DefaultCap(t *testing.T) {
	// maxPerLease <= 0 should default to 10
	h := NewRenewalHistory(0)
	for i := 0; i < 15; i++ {
		h.Record(makeRecord("lease-default", true))
	}
	if records := h.Get("lease-default"); len(records) != 10 {
		t.Errorf("expected default cap of 10, got %d", len(records))
	}
}

func TestRenewalHistory_GetReturnsCopy(t *testing.T) {
	h := NewRenewalHistory(5)
	h.Record(makeRecord("lease-copy", true))

	records := h.Get("lease-copy")
	records[0].LeaseID = "mutated"

	original := h.Get("lease-copy")
	if original[0].LeaseID == "mutated" {
		t.Error("Get should return a copy, not a reference to internal slice")
	}
}
