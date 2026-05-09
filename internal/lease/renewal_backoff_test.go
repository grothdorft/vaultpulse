package lease

import (
	"testing"
	"time"
)

func TestDefaultBackoffPolicy_Values(t *testing.T) {
	p := DefaultBackoffPolicy()
	if p.BaseDelay != 30*time.Second {
		t.Errorf("expected BaseDelay 30s, got %v", p.BaseDelay)
	}
	if p.MaxDelay != 10*time.Minute {
		t.Errorf("expected MaxDelay 10m, got %v", p.MaxDelay)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected Multiplier 2.0, got %v", p.Multiplier)
	}
}

func TestNextDelay_ZeroOrNegativeFailures(t *testing.T) {
	p := DefaultBackoffPolicy()
	for _, count := range []int{0, -1, -10} {
		if d := p.NextDelay(count); d != 0 {
			t.Errorf("failureCount=%d: expected 0, got %v", count, d)
		}
	}
}

func TestNextDelay_FirstFailure(t *testing.T) {
	p := DefaultBackoffPolicy()
	// exponent = 0 → delay = base * 1 = base
	got := p.NextDelay(1)
	if got != p.BaseDelay {
		t.Errorf("expected %v, got %v", p.BaseDelay, got)
	}
}

func TestNextDelay_SecondFailure(t *testing.T) {
	p := DefaultBackoffPolicy()
	// exponent = 1 → delay = 30s * 2 = 60s
	got := p.NextDelay(2)
	expected := 60 * time.Second
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestNextDelay_CappedAtMax(t *testing.T) {
	p := DefaultBackoffPolicy()
	// A very high failure count should be capped at MaxDelay.
	got := p.NextDelay(100)
	if got != p.MaxDelay {
		t.Errorf("expected MaxDelay %v, got %v", p.MaxDelay, got)
	}
}

func TestNextDelay_CustomPolicy(t *testing.T) {
	p := BackoffPolicy{
		BaseDelay:  10 * time.Second,
		MaxDelay:   1 * time.Minute,
		Multiplier: 3.0,
	}
	// failure 3 → 10s * 3^2 = 90s → capped at 60s
	got := p.NextDelay(3)
	if got != p.MaxDelay {
		t.Errorf("expected %v, got %v", p.MaxDelay, got)
	}
}

func TestShouldRetry_ZeroFailures(t *testing.T) {
	p := DefaultBackoffPolicy()
	// Zero failures → always retry regardless of last attempt.
	if !p.ShouldRetry(0, time.Now()) {
		t.Error("expected ShouldRetry=true for zero failures")
	}
}

func TestShouldRetry_DelayNotElapsed(t *testing.T) {
	p := DefaultBackoffPolicy()
	// Last attempt just now, first failure → delay not elapsed.
	if p.ShouldRetry(1, time.Now()) {
		t.Error("expected ShouldRetry=false when delay has not elapsed")
	}
}

func TestShouldRetry_DelayElapsed(t *testing.T) {
	p := DefaultBackoffPolicy()
	// Last attempt far in the past → delay definitely elapsed.
	past := time.Now().Add(-1 * time.Hour)
	if !p.ShouldRetry(1, past) {
		t.Error("expected ShouldRetry=true when delay has elapsed")
	}
}
