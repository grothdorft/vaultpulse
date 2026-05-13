package lease

import (
	"testing"
	"time"
)

func newTestCooldown(interval time.Duration) *RenewalCooldown {
	c := NewRenewalCooldown(CooldownPolicy{MinInterval: interval})
	return c
}

func TestRenewalCooldown_FirstCallAllowed(t *testing.T) {
	c := newTestCooldown(10 * time.Second)
	if !c.Allow("lease-1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestRenewalCooldown_SecondCallBlockedWithinWindow(t *testing.T) {
	now := time.Now()
	c := newTestCooldown(10 * time.Second)
	c.now = func() time.Time { return now }

	c.Allow("lease-1")
	if c.Allow("lease-1") {
		t.Fatal("expected second call within cooldown to be denied")
	}
}

func TestRenewalCooldown_AllowedAfterIntervalElapses(t *testing.T) {
	base := time.Now()
	c := newTestCooldown(10 * time.Second)
	c.now = func() time.Time { return base }

	c.Allow("lease-1")

	// Advance time beyond the cooldown window.
	c.now = func() time.Time { return base.Add(11 * time.Second) }
	if !c.Allow("lease-1") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestRenewalCooldown_IndependentLeases(t *testing.T) {
	c := newTestCooldown(10 * time.Second)
	c.Allow("lease-1")

	// A different lease should not be affected.
	if !c.Allow("lease-2") {
		t.Fatal("expected independent lease to be allowed")
	}
}

func TestRenewalCooldown_Reset(t *testing.T) {
	base := time.Now()
	c := newTestCooldown(10 * time.Second)
	c.now = func() time.Time { return base }

	c.Allow("lease-1")
	c.Reset("lease-1")

	if !c.Allow("lease-1") {
		t.Fatal("expected call after reset to be allowed")
	}
}

func TestRenewalCooldown_Remaining_WhenNotInCooldown(t *testing.T) {
	c := newTestCooldown(10 * time.Second)
	if r := c.Remaining("lease-unknown"); r != 0 {
		t.Fatalf("expected 0 remaining for unknown lease, got %v", r)
	}
}

func TestRenewalCooldown_Remaining_WithinWindow(t *testing.T) {
	base := time.Now()
	c := newTestCooldown(10 * time.Second)
	c.now = func() time.Time { return base }

	c.Allow("lease-1")

	// 3 seconds later — 7 seconds should remain.
	c.now = func() time.Time { return base.Add(3 * time.Second) }
	remaining := c.Remaining("lease-1")
	if remaining != 7*time.Second {
		t.Fatalf("expected 7s remaining, got %v", remaining)
	}
}

func TestDefaultCooldownPolicy_Values(t *testing.T) {
	p := DefaultCooldownPolicy()
	if p.MinInterval != 30*time.Second {
		t.Fatalf("expected 30s MinInterval, got %v", p.MinInterval)
	}
}
