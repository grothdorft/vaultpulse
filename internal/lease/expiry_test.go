package lease

import (
	"testing"
	"time"
)

const (
	warn     = 24 * time.Hour
	critical = 1 * time.Hour
)

func TestClassifyExpiry_OK(t *testing.T) {
	expiry := time.Now().Add(48 * time.Hour)
	got := ClassifyExpiry(expiry, warn, critical)
	if got != ExpiryLevelOK {
		t.Errorf("expected OK, got %s", got)
	}
}

func TestClassifyExpiry_Warning(t *testing.T) {
	expiry := time.Now().Add(12 * time.Hour)
	got := ClassifyExpiry(expiry, warn, critical)
	if got != ExpiryLevelWarning {
		t.Errorf("expected warning, got %s", got)
	}
}

func TestClassifyExpiry_Critical(t *testing.T) {
	expiry := time.Now().Add(30 * time.Minute)
	got := ClassifyExpiry(expiry, warn, critical)
	if got != ExpiryLevelCritical {
		t.Errorf("expected critical, got %s", got)
	}
}

func TestClassifyExpiry_Expired(t *testing.T) {
	expiry := time.Now().Add(-5 * time.Minute)
	got := ClassifyExpiry(expiry, warn, critical)
	if got != ExpiryLevelExpired {
		t.Errorf("expected expired, got %s", got)
	}
}

func TestClassifyExpiry_ExactlyAtCriticalBoundary(t *testing.T) {
	// TTL == criticalThreshold should be classified as critical
	expiry := time.Now().Add(critical)
	got := ClassifyExpiry(expiry, warn, critical)
	if got != ExpiryLevelCritical {
		t.Errorf("expected critical at exact boundary, got %s", got)
	}
}

func TestClassifyExpiry_ExactlyAtWarnBoundary(t *testing.T) {
	// TTL == warnThreshold should be classified as warning
	expiry := time.Now().Add(warn)
	got := ClassifyExpiry(expiry, warn, critical)
	if got != ExpiryLevelWarning {
		t.Errorf("expected warning at exact boundary, got %s", got)
	}
}

func TestExpiryLevel_String(t *testing.T) {
	cases := []struct {
		level ExpiryLevel
		want  string
	}{
		{ExpiryLevelOK, "ok"},
		{ExpiryLevelWarning, "warning"},
		{ExpiryLevelCritical, "critical"},
		{ExpiryLevelExpired, "expired"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ExpiryLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}
