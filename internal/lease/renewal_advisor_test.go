package lease

import (
	"testing"
	"time"
)

func makeEntry(id string, expiresIn time.Duration, alerted, expired bool) *Entry {
	return &Entry{
		LeaseID:   id,
		ExpiresAt: time.Now().Add(expiresIn),
		Alerted:   alerted,
		Expired:   expired,
	}
}

func TestAdviseBatch_EmptySlice(t *testing.T) {
	result := AdviseBatch(nil, DefaultRenewalPolicy())
	if len(result) != 0 {
		t.Fatalf("expected 0 advices, got %d", len(result))
	}
}

func TestAdviseBatch_NilEntrySkipped(t *testing.T) {
	entries := []*Entry{nil, nil}
	result := AdviseBatch(entries, DefaultRenewalPolicy())
	if len(result) != 0 {
		t.Fatalf("expected 0 advices for nil entries, got %d", len(result))
	}
}

func TestAdviseBatch_ExpiredEntrySkipped(t *testing.T) {
	entries := []*Entry{
		makeEntry("expired-lease", -10*time.Minute, false, true),
	}
	result := AdviseBatch(entries, DefaultRenewalPolicy())
	if len(result) != 0 {
		t.Fatalf("expected expired entry to be skipped, got %d advices", len(result))
	}
}

func TestAdviseBatch_HealthyEntryNoRenewal(t *testing.T) {
	entries := []*Entry{
		makeEntry("healthy-lease", 2*time.Hour, false, false),
	}
	result := AdviseBatch(entries, DefaultRenewalPolicy())
	if len(result) != 1 {
		t.Fatalf("expected 1 advice, got %d", len(result))
	}
	if result[0].ShouldRenew {
		t.Errorf("expected ShouldRenew=false for healthy lease")
	}
}

func TestAdviseBatch_ExpiringEntryShouldRenew(t *testing.T) {
	policy := DefaultRenewalPolicy()
	entries := []*Entry{
		makeEntry("expiring-lease", policy.RenewalWindow-time.Second, false, false),
	}
	result := AdviseBatch(entries, policy)
	if len(result) != 1 {
		t.Fatalf("expected 1 advice, got %d", len(result))
	}
	if !result[0].ShouldRenew {
		t.Errorf("expected ShouldRenew=true for expiring lease")
	}
	if result[0].LeaseID != "expiring-lease" {
		t.Errorf("unexpected lease ID: %s", result[0].LeaseID)
	}
}

func TestRenewalAdvice_String_NoRenew(t *testing.T) {
	adv := RenewalAdvice{LeaseID: "abc", ShouldRenew: false, TimeRemaining: 90 * time.Second}
	s := adv.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestRenewalAdvice_String_Urgent(t *testing.T) {
	adv := RenewalAdvice{LeaseID: "abc", ShouldRenew: true, Urgent: true, Reason: "critical", TimeRemaining: 30 * time.Second}
	s := adv.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
