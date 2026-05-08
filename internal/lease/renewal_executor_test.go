package lease

import (
	"context"
	"errors"
	"testing"
	"time"
)

func stubRenewOK(_ context.Context, _ string) (time.Duration, error) {
	return 3600 * time.Second, nil
}

func stubRenewFail(_ context.Context, _ string) (time.Duration, error) {
	return 0, errors.New("vault unavailable")
}

func newTestExecutor(renew RenewFunc) (*RenewalExecutor, *RenewalHistory) {
	h := NewRenewalHistory(10)
	exec := NewRenewalExecutor(h, renew, nil)
	return exec, h
}

func TestRenewalExecutor_SkipsNoRenewAdvice(t *testing.T) {
	exec, h := newTestExecutor(stubRenewOK)
	advice := []RenewalAdvice{{LeaseID: "lease/1", ShouldRenew: false}}
	n := exec.ExecuteAdvice(context.Background(), advice)
	if n != 0 {
		t.Fatalf("expected 0 renewals, got %d", n)
	}
	if records := h.Get("lease/1"); len(records) != 0 {
		t.Fatalf("expected no history, got %d records", len(records))
	}
}

func TestRenewalExecutor_SuccessfulRenewal(t *testing.T) {
	exec, h := newTestExecutor(stubRenewOK)
	advice := []RenewalAdvice{{LeaseID: "lease/2", ShouldRenew: true}}
	n := exec.ExecuteAdvice(context.Background(), advice)
	if n != 1 {
		t.Fatalf("expected 1 renewal, got %d", n)
	}
	records := h.Get("lease/2")
	if len(records) != 1 {
		t.Fatalf("expected 1 history record, got %d", len(records))
	}
	if records[0].Error != "" {
		t.Fatalf("expected no error in record, got %q", records[0].Error)
	}
	if records[0].NewTTL != 3600*time.Second {
		t.Fatalf("unexpected NewTTL: %v", records[0].NewTTL)
	}
}

func TestRenewalExecutor_FailedRenewalRecordsError(t *testing.T) {
	exec, h := newTestExecutor(stubRenewFail)
	advice := []RenewalAdvice{{LeaseID: "lease/3", ShouldRenew: true}}
	n := exec.ExecuteAdvice(context.Background(), advice)
	if n != 0 {
		t.Fatalf("expected 0 successes, got %d", n)
	}
	records := h.Get("lease/3")
	if len(records) != 1 {
		t.Fatalf("expected 1 history record, got %d", len(records))
	}
	if records[0].Error == "" {
		t.Fatal("expected error in record, got empty string")
	}
}

func TestRenewalExecutor_MixedBatch(t *testing.T) {
	exec, _ := newTestExecutor(stubRenewOK)
	advice := []RenewalAdvice{
		{LeaseID: "lease/a", ShouldRenew: true},
		{LeaseID: "lease/b", ShouldRenew: false},
		{LeaseID: "lease/c", ShouldRenew: true},
	}
	n := exec.ExecuteAdvice(context.Background(), advice)
	if n != 2 {
		t.Fatalf("expected 2 renewals, got %d", n)
	}
}
