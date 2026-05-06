package lease_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

func TestShouldAlert_NotAlertedNotExpired(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "secret/data/myapp/db#abc123",
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Alerted:   false,
		Expired:   false,
	}

	if !lease.ShouldAlert(entry) {
		t.Error("expected ShouldAlert to return true for un-alerted, non-expired entry")
	}
}

func TestShouldAlert_AlreadyAlerted(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "secret/data/myapp/db#abc123",
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Alerted:   true,
		Expired:   false,
	}

	if lease.ShouldAlert(entry) {
		t.Error("expected ShouldAlert to return false for already-alerted entry")
	}
}

func TestShouldAlert_AlreadyExpired(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "secret/data/myapp/db#abc123",
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Alerted:   false,
		Expired:   true,
	}

	if lease.ShouldAlert(entry) {
		t.Error("expected ShouldAlert to return false for already-expired entry")
	}
}

func TestShouldAlert_BothAlertedAndExpired(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "secret/data/myapp/db#abc123",
		ExpiresAt: time.Now().Add(-5 * time.Minute),
		Alerted:   true,
		Expired:   true,
	}

	if lease.ShouldAlert(entry) {
		t.Error("expected ShouldAlert to return false for alerted and expired entry")
	}
}

func TestShouldAlert_NilEntry(t *testing.T) {
	if lease.ShouldAlert(nil) {
		t.Error("expected ShouldAlert to return false for nil entry")
	}
}

// newEntry is a test helper that constructs a lease.Entry with sensible defaults,
// reducing boilerplate in table-driven tests.
func newEntry(alerted, expired bool, expiresIn time.Duration) *lease.Entry {
	return &lease.Entry{
		LeaseID:   "secret/data/myapp/db#abc123",
		ExpiresAt: time.Now().Add(expiresIn),
		Alerted:   alerted,
		Expired:   expired,
	}
}

func TestShouldAlert_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		entry    *lease.Entry
		wantAlert bool
	}{
		{"not alerted, not expired", newEntry(false, false, 10*time.Minute), true},
		{"alerted, not expired", newEntry(true, false, 10*time.Minute), false},
		{"not alerted, expired", newEntry(false, true, -5*time.Minute), false},
		{"alerted and expired", newEntry(true, true, -5*time.Minute), false},
		{"nil entry", nil, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := lease.ShouldAlert(tc.entry)
			if got != tc.wantAlert {
				t.Errorf("ShouldAlert() = %v, want %v", got, tc.wantAlert)
			}
		})
	}
}
