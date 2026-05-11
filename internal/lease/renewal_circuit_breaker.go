package lease

import (
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	CircuitClosed   CircuitState = iota // normal operation
	CircuitOpen                         // blocking requests
	CircuitHalfOpen                     // testing recovery
)

// CircuitBreakerPolicy configures the circuit breaker behaviour.
type CircuitBreakerPolicy struct {
	FailureThreshold int
	SuccessThreshold int
	OpenTimeout      time.Duration
}

// DefaultCircuitBreakerPolicy returns sensible defaults.
func DefaultCircuitBreakerPolicy() CircuitBreakerPolicy {
	return CircuitBreakerPolicy{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		OpenTimeout:      30 * time.Second,
	}
}

// RenewalCircuitBreaker prevents renewal attempts when a lease endpoint
// is repeatedly failing.
type RenewalCircuitBreaker struct {
	mu             sync.Mutex
	policy         CircuitBreakerPolicy
	state          CircuitState
	failureCount   int
	successCount   int
	lastOpenedAt   time.Time
}

// NewRenewalCircuitBreaker creates a new breaker with the given policy.
func NewRenewalCircuitBreaker(p CircuitBreakerPolicy) *RenewalCircuitBreaker {
	return &RenewalCircuitBreaker{policy: p}
}

// Allow returns true if a renewal attempt should be permitted.
func (cb *RenewalCircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastOpenedAt) >= cb.policy.OpenTimeout {
			cb.state = CircuitHalfOpen
			cb.successCount = 0
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful renewal, potentially closing the circuit.
func (cb *RenewalCircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == CircuitHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.policy.SuccessThreshold {
			cb.state = CircuitClosed
			cb.failureCount = 0
		}
	} else {
		cb.failureCount = 0
	}
}

// RecordFailure records a failed renewal, potentially opening the circuit.
func (cb *RenewalCircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	if cb.state == CircuitHalfOpen || cb.failureCount >= cb.policy.FailureThreshold {
		cb.state = CircuitOpen
		cb.lastOpenedAt = time.Now()
		cb.failureCount = 0
	}
}

// State returns the current circuit state.
func (cb *RenewalCircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
