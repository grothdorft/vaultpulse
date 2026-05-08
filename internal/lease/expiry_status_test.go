package lease_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

const (
	warnThreshold = 24 * time.Hour
	critThreshold = 4 * time.Hour
)

func TestBuildStatusSummary_NilEntry(t *testing.T) {
	summary := lease.BuildStatusSummary(nil, warnThreshold, critThreshold)
	if summary.Status != lease.StatusExpired {
		t.Errorf("expected StatusExpired for nil entry, got %s", summary.Status)
	}
}

func TestBuildStatusSummary_OK(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "lease/ok",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	summary := lease.BuildStatusSummary(entry, warnThreshold, critThreshold)
	if summary.Status != lease.StatusOK {
		t.Errorf("expected StatusOK, got %s", summary.Status)
	}
	if summary.LeaseID != entry.LeaseID {
		t.Errorf("expected LeaseID %s, got %s", entry.LeaseID, summary.LeaseID)
	}
}

func TestBuildStatusSummary_Warning(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "lease/warn",
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}
	summary := lease.BuildStatusSummary(entry, warnThreshold, critThreshold)
	if summary.Status != lease.StatusWarning {
		t.Errorf("expected StatusWarning, got %s", summary.Status)
	}
}

func TestBuildStatusSummary_Critical(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "lease/crit",
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}
	summary := lease.BuildStatusSummary(entry, warnThreshold, critThreshold)
	if summary.Status != lease.StatusCritical {
		t.Errorf("expected StatusCritical, got %s", summary.Status)
	}
}

func TestBuildStatusSummary_Expired(t *testing.T) {
	entry := &lease.Entry{
		LeaseID:   "lease/expired",
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	summary := lease.BuildStatusSummary(entry, warnThreshold, critThreshold)
	if summary.Status != lease.StatusExpired {
		t.Errorf("expected StatusExpired, got %s", summary.Status)
	}
}

func TestBuildStatusSummary_TTLIsApproximate(t *testing.T) {
	duration := 30 * time.Hour
	entry := &lease.Entry{
		LeaseID:   "lease/ttl",
		ExpiresAt: time.Now().Add(duration),
	}
	summary := lease.BuildStatusSummary(entry, warnThreshold, critThreshold)
	if summary.TTL <= 0 {
		t.Errorf("expected positive TTL, got %v", summary.TTL)
	}
	if summary.TTL > duration {
		t.Errorf("TTL %v exceeds original duration %v", summary.TTL, duration)
	}
}
