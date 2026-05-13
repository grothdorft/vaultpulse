package lease

import (
	"testing"
	"time"
)

func newTestRateLimiter(p RateLimiterPolicy) (*RenewalRateLimiter, *time.Time) {
	now := time.Now()
	rl := NewRenewalRateLimiter(p)
	rl.lastAt = now
	rl.nowFn = func() time.Time { return now }
	return rl, &now
}

func TestDefaultRateLimiterPolicy_Values(t *testing.T) {
	p := DefaultRateLimiterPolicy()
	if p.TokensPerInterval != 10 {
		t.Errorf("expected TokensPerInterval=10, got %d", p.TokensPerInterval)
	}
	if p.Interval != time.Minute {
		t.Errorf("expected Interval=1m, got %v", p.Interval)
	}
	if p.BurstCapacity != 20 {
		t.Errorf("expected BurstCapacity=20, got %d", p.BurstCapacity)
	}
}

func TestRateLimiter_AllowConsumesToken(t *testing.T) {
	p := RateLimiterPolicy{TokensPerInterval: 5, Interval: time.Minute, BurstCapacity: 5}
	rl, _ := newTestRateLimiter(p)
	rl.tokens = 3

	if !rl.Allow() {
		t.Fatal("expected Allow() = true")
	}
	if rl.tokens != 2 {
		t.Errorf("expected tokens=2 after Allow, got %d", rl.tokens)
	}
}

func TestRateLimiter_DeniedWhenEmpty(t *testing.T) {
	p := RateLimiterPolicy{TokensPerInterval: 5, Interval: time.Minute, BurstCapacity: 5}
	rl, _ := newTestRateLimiter(p)
	rl.tokens = 0

	if rl.Allow() {
		t.Fatal("expected Allow() = false when no tokens")
	}
}

func TestRateLimiter_ReplenishesAfterInterval(t *testing.T) {
	p := RateLimiterPolicy{TokensPerInterval: 5, Interval: time.Minute, BurstCapacity: 10}
	now := time.Now()
	rl := NewRenewalRateLimiter(p)
	rl.tokens = 0
	rl.lastAt = now
	rl.nowFn = func() time.Time { return now.Add(2 * time.Minute) }

	if !rl.Allow() {
		t.Fatal("expected Allow() = true after replenishment")
	}
	// 2 intervals * 5 tokens = 10, capped at BurstCapacity=10, then -1 for Allow
	if rl.tokens != 9 {
		t.Errorf("expected tokens=9 after replenish+allow, got %d", rl.tokens)
	}
}

func TestRateLimiter_CappedAtBurst(t *testing.T) {
	p := RateLimiterPolicy{TokensPerInterval: 100, Interval: time.Minute, BurstCapacity: 15}
	now := time.Now()
	rl := NewRenewalRateLimiter(p)
	rl.tokens = 0
	rl.lastAt = now
	rl.nowFn = func() time.Time { return now.Add(5 * time.Minute) }

	avail := rl.Available()
	if avail != 15 {
		t.Errorf("expected Available()=15 (capped), got %d", avail)
	}
}

func TestRateLimiter_PolicyReturnsConfig(t *testing.T) {
	p := DefaultRateLimiterPolicy()
	rl := NewRenewalRateLimiter(p)
	got := rl.Policy()
	if got.BurstCapacity != p.BurstCapacity {
		t.Errorf("Policy() mismatch: got %+v", got)
	}
}
