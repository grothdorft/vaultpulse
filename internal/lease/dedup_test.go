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
