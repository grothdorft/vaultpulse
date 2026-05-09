package lease

import (
	"testing"
	"time"
)

func newTestThrottle(max int, interval time.Duration) *RenewalThrottle {
	t := NewRenewalThrottle(ThrottlePolicy{MaxPerInterval: max, Interval: interval})
	return t
}

func TestRenewalThrottle_AllowsUpToLimit(t *testing.T) {
	th := newTestThrottle(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !th.Allow() {
			t.Fatalf("expected Allow()=true on attempt %d", i+1)
		}
	}
	if th.Allow() {
		t.Fatal("expected Allow()=false after limit reached")
	}
}

func TestRenewalThrottle_RemainingDecrementsCorrectly(t *testing.T) {
	th := newTestThrottle(5, time.Minute)
	if th.Remaining() != 5 {
		t.Fatalf("expected 5 remaining, got %d", th.Remaining())
	}
	th.Allow()
	th.Allow()
	if th.Remaining() != 3 {
		t.Fatalf("expected 3 remaining, got %d", th.Remaining())
	}
}

func TestRenewalThrottle_ResetsAfterInterval(t *testing.T) {
	now := time.Now()
	th := newTestThrottle(2, time.Minute)

	// Inject a custom clock that starts in the past.
	th.nowFn = func() time.Time { return now.Add(-90 * time.Second) }
	th.Allow()
	th.Allow()

	// Advance clock beyond the interval window.
	th.nowFn = func() time.Time { return now }

	if !th.Allow() {
		t.Fatal("expected Allow()=true after window reset")
	}
}

func TestRenewalThrottle_RemainingNeverNegative(t *testing.T) {
	th := newTestThrottle(1, time.Minute)
	th.Allow()
	th.Allow() // exceeds limit, should be blocked
	if th.Remaining() < 0 {
		t.Fatal("Remaining() should never be negative")
	}
}

func TestDefaultThrottlePolicy_Values(t *testing.T) {
	p := DefaultThrottlePolicy()
	if p.MaxPerInterval != 10 {
		t.Errorf("expected MaxPerInterval=10, got %d", p.MaxPerInterval)
	}
	if p.Interval != time.Minute {
		t.Errorf("expected Interval=1m, got %v", p.Interval)
	}
}

func TestRenewalThrottle_ConcurrentAccess(t *testing.T) {
	th := newTestThrottle(50, time.Minute)
	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func() {
			th.Allow()
			th.Remaining()
			done <- struct{}{}
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
