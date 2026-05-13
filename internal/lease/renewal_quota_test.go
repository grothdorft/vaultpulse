package lease

import (
	"testing"
	"time"
)

func newTestQuota(max int, window time.Duration) *RenewalQuota {
	return NewRenewalQuota(QuotaPolicy{
		MaxRenewalsPerWindow: max,
		Window:               window,
	})
}

func TestDefaultQuotaPolicy_Values(t *testing.T) {
	p := DefaultQuotaPolicy()
	if p.MaxRenewalsPerWindow != 50 {
		t.Errorf("expected 50, got %d", p.MaxRenewalsPerWindow)
	}
	if p.Window != time.Minute {
		t.Errorf("expected 1m, got %s", p.Window)
	}
}

func TestRenewalQuota_AllowsUpToLimit(t *testing.T) {
	q := newTestQuota(3, time.Minute)
	for i := 0; i < 3; i++ {
		ok, err := q.Allow("secret/myapp")
		if !ok || err != nil {
			t.Fatalf("expected allow on iteration %d", i)
		}
	}
}

func TestRenewalQuota_DeniedAfterLimit(t *testing.T) {
	q := newTestQuota(2, time.Minute)
	q.Allow("ns/db") //nolint
	q.Allow("ns/db") //nolint

	ok, err := q.Allow("ns/db")
	if ok {
		t.Error("expected denial after quota exceeded")
	}
	if err == nil {
		t.Error("expected non-nil error on denial")
	}
}

func TestRenewalQuota_IndependentKeys(t *testing.T) {
	q := newTestQuota(1, time.Minute)
	q.Allow("ns/a") //nolint

	ok, err := q.Allow("ns/b")
	if !ok || err != nil {
		t.Error("expected allow for independent key")
	}
}

func TestRenewalQuota_ResetsAfterWindow(t *testing.T) {
	q := newTestQuota(1, 10*time.Millisecond)
	q.Allow("ns/x") //nolint

	ok, _ := q.Allow("ns/x")
	if ok {
		t.Fatal("expected denial before window reset")
	}

	time.Sleep(20 * time.Millisecond)

	ok, err := q.Allow("ns/x")
	if !ok || err != nil {
		t.Error("expected allow after window reset")
	}
}

func TestRenewalQuota_Snapshot(t *testing.T) {
	q := newTestQuota(10, time.Minute)
	q.Allow("ns/a") //nolint
	q.Allow("ns/a") //nolint
	q.Allow("ns/b") //nolint

	snap := q.Snapshot()
	if snap["ns/a"].Count != 2 {
		t.Errorf("expected count 2 for ns/a, got %d", snap["ns/a"].Count)
	}
	if snap["ns/b"].Count != 1 {
		t.Errorf("expected count 1 for ns/b, got %d", snap["ns/b"].Count)
	}
}

func TestRenewalQuota_PolicyReturned(t *testing.T) {
	q := newTestQuota(7, 30*time.Second)
	p := q.Policy()
	if p.MaxRenewalsPerWindow != 7 {
		t.Errorf("expected 7, got %d", p.MaxRenewalsPerWindow)
	}
}
