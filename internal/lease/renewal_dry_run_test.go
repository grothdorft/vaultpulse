package lease

import (
	"testing"
)

func TestDryRunLog_RecordAndAll(t *testing.T) {
	log := NewDryRunLog(10)
	log.Record("lease-1", RenewalAdvice{ShouldRenew: true}, "within window")
	log.Record("lease-2", RenewalAdvice{ShouldRenew: false}, "not due")

	all := log.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 records, got %d", len(all))
	}
	if all[0].LeaseID != "lease-1" {
		t.Errorf("expected lease-1, got %s", all[0].LeaseID)
	}
	if all[1].Advice.ShouldRenew {
		t.Errorf("expected ShouldRenew=false for lease-2")
	}
}

func TestDryRunLog_CapEvictsOldest(t *testing.T) {
	log := NewDryRunLog(3)
	for i := 0; i < 5; i++ {
		log.Record("lease", RenewalAdvice{}, "")
	}
	if log.Len() != 3 {
		t.Errorf("expected cap 3, got %d", log.Len())
	}
}

func TestDryRunLog_Clear(t *testing.T) {
	log := NewDryRunLog(10)
	log.Record("lease-1", RenewalAdvice{}, "")
	log.Record("lease-2", RenewalAdvice{}, "")
	log.Clear()
	if log.Len() != 0 {
		t.Errorf("expected 0 after clear, got %d", log.Len())
	}
}

func TestDryRunLog_DefaultCap(t *testing.T) {
	log := NewDryRunLog(0)
	if log.cap != 100 {
		t.Errorf("expected default cap 100, got %d", log.cap)
	}
}

func TestDryRunLog_AllReturnsCopy(t *testing.T) {
	log := NewDryRunLog(10)
	log.Record("lease-1", RenewalAdvice{ShouldRenew: true}, "ok")

	all := log.All()
	all[0].LeaseID = "mutated"

	original := log.All()
	if original[0].LeaseID == "mutated" {
		t.Error("All() should return a copy, not a reference")
	}
}

func TestDryRunLog_Len(t *testing.T) {
	log := NewDryRunLog(10)
	if log.Len() != 0 {
		t.Errorf("expected 0, got %d", log.Len())
	}
	log.Record("lease-1", RenewalAdvice{}, "")
	if log.Len() != 1 {
		t.Errorf("expected 1, got %d", log.Len())
	}
}
