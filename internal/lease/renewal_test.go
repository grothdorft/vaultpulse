package lease

import (
	"testing"
	"time"
)

func baseEntry(leaseID string, expiresIn time.Duration) *Entry {
	return &Entry{
		LeaseID:   leaseID,
		ExpiresAt: time.Now().Add(expiresIn),
		Expired:   false,
		Alerted:   false,
	}
}

func TestAdviseRenewal_NilEntry(t *testing.T) {
	advice := AdviseRenewal(nil, DefaultRenewalPolicy(), time.Now())
	if advice.ShouldRenew {
		t.Error("expected ShouldRenew=false for nil entry")
	}
	if advice.Reason != "nil entry" {
		t.Errorf("unexpected reason: %s", advice.Reason)
	}
}

func TestAdviseRenewal_AlreadyExpired(t *testing.T) {
	entry := baseEntry("lease/1", -1*time.Minute)
	entry.Expired = true
	advice := AdviseRenewal(entry, DefaultRenewalPolicy(), time.Now())
	if advice.ShouldRenew {
		t.Error("expected ShouldRenew=false for expired entry")
	}
	if advice.Urgency != StatusExpired {
		t.Errorf("expected StatusExpired urgency, got %v", advice.Urgency)
	}
}

func TestAdviseRenewal_InsufficientRemaining(t *testing.T) {
	policy := DefaultRenewalPolicy()
	// only 10 seconds left — below MinRemainingAfterRenewal (30s)
	entry := baseEntry("lease/2", 10*time.Second)
	advice := AdviseRenewal(entry, policy, time.Now())
	if advice.ShouldRenew {
		t.Error("expected ShouldRenew=false when remaining < MinRemainingAfterRenewal")
	}
}

func TestAdviseRenewal_WithinRenewalWindow(t *testing.T) {
	policy := DefaultRenewalPolicy()
	// 5 minutes left — within the 10-minute renewal window
	entry := baseEntry("lease/3", 5*time.Minute)
	advice := AdviseRenewal(entry, policy, time.Now())
	if !advice.ShouldRenew {
		t.Error("expected ShouldRenew=true when within renewal window")
	}
	if advice.LeaseID != "lease/3" {
		t.Errorf("unexpected LeaseID: %s", advice.LeaseID)
	}
}

func TestAdviseRenewal_SufficientTTL(t *testing.T) {
	policy := DefaultRenewalPolicy()
	// 1 hour left — well outside renewal window
	entry := baseEntry("lease/4", 1*time.Hour)
	advice := AdviseRenewal(entry, policy, time.Now())
	if advice.ShouldRenew {
		t.Error("expected ShouldRenew=false when TTL is sufficient")
	}
	if advice.Urgency != StatusOK {
		t.Errorf("expected StatusOK urgency, got %v", advice.Urgency)
	}
}

func TestAdviseRenewal_CustomPolicy(t *testing.T) {
	policy := RenewalPolicy{
		RenewBeforeTTL:           2 * time.Minute,
		MinRemainingAfterRenewal: 10 * time.Second,
	}
	// 90 seconds left — within custom 2-minute window
	entry := baseEntry("lease/5", 90*time.Second)
	advice := AdviseRenewal(entry, policy, time.Now())
	if !advice.ShouldRenew {
		t.Error("expected ShouldRenew=true with custom policy")
	}
}
