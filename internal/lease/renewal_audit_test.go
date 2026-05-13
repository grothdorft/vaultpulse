package lease

import (
	"fmt"
	"testing"
)

func TestRenewalAuditLog_RecordAndAll(t *testing.T) {
	log := NewRenewalAuditLog(10)
	log.Record("lease-1", AuditKindSuccess, "renewed")
	log.Record("lease-2", AuditKindFailure, "timeout")

	events := log.All()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].LeaseID != "lease-1" || events[0].Kind != AuditKindSuccess {
		t.Errorf("unexpected first event: %+v", events[0])
	}
	if events[1].LeaseID != "lease-2" || events[1].Kind != AuditKindFailure {
		t.Errorf("unexpected second event: %+v", events[1])
	}
}

func TestRenewalAuditLog_CapEvictsOldest(t *testing.T) {
	log := NewRenewalAuditLog(3)
	for i := 0; i < 5; i++ {
		log.Record(fmt.Sprintf("lease-%d", i), AuditKindSuccess, "")
	}

	events := log.All()
	if len(events) != 3 {
		t.Fatalf("expected cap of 3, got %d", len(events))
	}
	// oldest two (lease-0, lease-1) should have been evicted
	if events[0].LeaseID != "lease-2" {
		t.Errorf("expected lease-2 as oldest remaining, got %s", events[0].LeaseID)
	}
}

func TestRenewalAuditLog_ForLease(t *testing.T) {
	log := NewRenewalAuditLog(20)
	log.Record("lease-A", AuditKindSuccess, "ok")
	log.Record("lease-B", AuditKindFailure, "err")
	log.Record("lease-A", AuditKindThrottled, "throttled")

	results := log.ForLease("lease-A")
	if len(results) != 2 {
		t.Fatalf("expected 2 events for lease-A, got %d", len(results))
	}
	for _, e := range results {
		if e.LeaseID != "lease-A" {
			t.Errorf("unexpected lease ID in filtered results: %s", e.LeaseID)
		}
	}
}

func TestRenewalAuditLog_ForLease_Unknown(t *testing.T) {
	log := NewRenewalAuditLog(10)
	log.Record("lease-X", AuditKindSuccess, "")

	results := log.ForLease("lease-unknown")
	if results != nil && len(results) != 0 {
		t.Errorf("expected empty result for unknown lease, got %v", results)
	}
}

func TestRenewalAuditLog_Len(t *testing.T) {
	log := NewRenewalAuditLog(10)
	if log.Len() != 0 {
		t.Fatalf("expected 0 initially, got %d", log.Len())
	}
	log.Record("l1", AuditKindSuccess, "")
	log.Record("l2", AuditKindCircuitOpen, "")
	if log.Len() != 2 {
		t.Fatalf("expected 2, got %d", log.Len())
	}
}

func TestRenewalAuditLog_DefaultCap(t *testing.T) {
	// cap <= 0 should default to 256
	log := NewRenewalAuditLog(0)
	for i := 0; i < 300; i++ {
		log.Record(fmt.Sprintf("lease-%d", i), AuditKindSuccess, "")
	}
	if log.Len() != 256 {
		t.Errorf("expected default cap of 256, got %d", log.Len())
	}
}

func TestRenewalAuditLog_OccurredAtSet(t *testing.T) {
	log := NewRenewalAuditLog(5)
	log.Record("lease-ts", AuditKindSkipped, "no advice")
	events := log.All()
	if events[0].OccurredAt.IsZero() {
		t.Error("expected OccurredAt to be set")
	}
}
