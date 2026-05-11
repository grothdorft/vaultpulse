package lease

import (
	"testing"
	"time"
)

func newTestBreaker() *RenewalCircuitBreaker {
	return NewRenewalCircuitBreaker(CircuitBreakerPolicy{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		OpenTimeout:      50 * time.Millisecond,
	})
}

func TestCircuitBreaker_InitiallyClosed(t *testing.T) {
	cb := newTestBreaker()
	if cb.State() != CircuitClosed {
		t.Fatalf("expected Closed, got %v", cb.State())
	}
	if !cb.Allow() {
		t.Fatal("expected Allow() == true when closed")
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := newTestBreaker()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != CircuitOpen {
		t.Fatalf("expected Open after threshold, got %v", cb.State())
	}
	if cb.Allow() {
		t.Fatal("expected Allow() == false when open")
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := newTestBreaker()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	if !cb.Allow() {
		t.Fatal("expected Allow() == true in half-open state")
	}
	if cb.State() != CircuitHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_ClosesAfterSuccessThreshold(t *testing.T) {
	cb := newTestBreaker()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	cb.Allow() // transition to half-open
	cb.RecordSuccess()
	if cb.State() != CircuitHalfOpen {
		t.Fatal("should still be half-open after one success")
	}
	cb.RecordSuccess()
	if cb.State() != CircuitClosed {
		t.Fatalf("expected Closed after success threshold, got %v", cb.State())
	}
}

func TestCircuitBreaker_ReopensOnFailureInHalfOpen(t *testing.T) {
	cb := newTestBreaker()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	cb.Allow()
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatalf("expected Open after failure in half-open, got %v", cb.State())
	}
}

func TestCircuitBreaker_SuccessResetsClosed(t *testing.T) {
	cb := newTestBreaker()
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess() // should reset failure count
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitClosed {
		t.Fatal("circuit should remain closed after success reset")
	}
}

func TestDefaultCircuitBreakerPolicy(t *testing.T) {
	p := DefaultCircuitBreakerPolicy()
	if p.FailureThreshold <= 0 {
		t.Error("FailureThreshold must be positive")
	}
	if p.SuccessThreshold <= 0 {
		t.Error("SuccessThreshold must be positive")
	}
	if p.OpenTimeout <= 0 {
		t.Error("OpenTimeout must be positive")
	}
}
