package lease_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

func TestTracker_UpsertAndGet(t *testing.T) {
	tr := lease.NewTracker()
	expiry := time.Now().Add(10 * time.Minute)

	tr.Upsert("lease-1", expiry)

	e := tr.Get("lease-1")
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.LeaseID != "lease-1" {
		t.Errorf("expected lease-1, got %s", e.LeaseID)
	}
	if e.State != lease.StateHealthy {
		t.Errorf("expected StateHealthy, got %v", e.State)
	}
}

func TestTracker_UpsertUpdatesExpiry(t *testing.T) {
	tr := lease.NewTracker()
	original := time.Now().Add(5 * time.Minute)
	updated := time.Now().Add(15 * time.Minute)

	tr.Upsert("lease-2", original)
	tr.Upsert("lease-2", updated)

	e := tr.Get("lease-2")
	if !e.ExpiresAt.Equal(updated) {
		t.Errorf("expected updated expiry, got %v", e.ExpiresAt)
	}
}

func TestTracker_MarkAlerted(t *testing.T) {
	tr := lease.NewTracker()
	tr.Upsert("lease-3", time.Now().Add(2*time.Minute))
	tr.MarkAlerted("lease-3")

	e := tr.Get("lease-3")
	if e.State != lease.StateAlerted {
		t.Errorf("expected StateAlerted, got %v", e.State)
	}
	if e.AlertedAt == nil {
		t.Error("expected AlertedAt to be set")
	}
}

func TestTracker_MarkExpired(t *testing.T) {
	tr := lease.NewTracker()
	tr.Upsert("lease-4", time.Now().Add(1*time.Minute))
	tr.MarkExpired("lease-4")

	e := tr.Get("lease-4")
	if e.State != lease.StateExpired {
		t.Errorf("expected StateExpired, got %v", e.State)
	}
}

func TestTracker_Remove(t *testing.T) {
	tr := lease.NewTracker()
	tr.Upsert("lease-5", time.Now().Add(5*time.Minute))
	tr.Remove("lease-5")

	if tr.Get("lease-5") != nil {
		t.Error("expected nil after removal")
	}
}

func TestTracker_All(t *testing.T) {
	tr := lease.NewTracker()
	tr.Upsert("lease-a", time.Now().Add(5*time.Minute))
	tr.Upsert("lease-b", time.Now().Add(10*time.Minute))

	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestTracker_GetMissing(t *testing.T) {
	tr := lease.NewTracker()
	if tr.Get("nonexistent") != nil {
		t.Error("expected nil for missing lease")
	}
}
