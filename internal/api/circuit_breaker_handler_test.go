package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/lease"
)

func newTestCircuitBreakerHandler() (http.HandlerFunc, *lease.RenewalCircuitBreaker, lease.CircuitBreakerPolicy) {
	policy := lease.CircuitBreakerPolicy{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		OpenTimeout:      30 * time.Second,
	}
	cb := lease.NewRenewalCircuitBreaker(policy)
	return handleCircuitBreakerStatus(cb, policy), cb, policy
}

func TestCircuitBreakerHandler_ClosedState(t *testing.T) {
	h, _, _ := newTestCircuitBreakerHandler()
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodGet, "/circuit-breaker", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp circuitBreakerStatusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.State != "closed" {
		t.Errorf("expected state=closed, got %s", resp.State)
	}
}

func TestCircuitBreakerHandler_OpenState(t *testing.T) {
	h, cb, _ := newTestCircuitBreakerHandler()
	for i := 0; i < 5; i++ {
		cb.RecordFailure()
	}
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodGet, "/circuit-breaker", nil))

	var resp circuitBreakerStatusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.State != "open" {
		t.Errorf("expected state=open, got %s", resp.State)
	}
}

func TestCircuitBreakerHandler_PolicyFields(t *testing.T) {
	h, _, _ := newTestCircuitBreakerHandler()
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodGet, "/circuit-breaker", nil))

	var resp circuitBreakerStatusResponse
	json.NewDecoder(rr.Body).Decode(&resp) //nolint:errcheck

	if resp.FailureThreshold != 5 {
		t.Errorf("expected failure_threshold=5, got %d", resp.FailureThreshold)
	}
	if resp.SuccessThreshold != 2 {
		t.Errorf("expected success_threshold=2, got %d", resp.SuccessThreshold)
	}
	if resp.OpenTimeoutSec != 30 {
		t.Errorf("expected open_timeout_seconds=30, got %d", resp.OpenTimeoutSec)
	}
}

func TestCircuitBreakerHandler_MethodNotAllowed(t *testing.T) {
	h, _, _ := newTestCircuitBreakerHandler()
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodPost, "/circuit-breaker", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestStateString_AllValues(t *testing.T) {
	cases := map[lease.CircuitState]string{
		lease.CircuitClosed:   "closed",
		lease.CircuitOpen:     "open",
		lease.CircuitHalfOpen: "half_open",
		lease.CircuitState(99): "unknown",
	}
	for state, want := range cases {
		if got := stateString(state); got != want {
			t.Errorf("stateString(%v) = %s, want %s", state, got, want)
		}
	}
}
