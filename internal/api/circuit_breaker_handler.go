package api

import (
	"net/http"

	"github.com/your-org/vaultpulse/internal/lease"
)

// circuitBreakerStatusResponse is the JSON shape returned by the handler.
type circuitBreakerStatusResponse struct {
	State            string `json:"state"`
	FailureThreshold int    `json:"failure_threshold"`
	SuccessThreshold int    `json:"success_threshold"`
	OpenTimeoutSec   int    `json:"open_timeout_seconds"`
}

func stateString(s lease.CircuitState) string {
	switch s {
	case lease.CircuitClosed:
		return "closed"
	case lease.CircuitOpen:
		return "open"
	case lease.CircuitHalfOpen:
		return "half_open"
	default:
		return "unknown"
	}
}

func handleCircuitBreakerStatus(cb *lease.RenewalCircuitBreaker, policy lease.CircuitBreakerPolicy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := circuitBreakerStatusResponse{
			State:            stateString(cb.State()),
			FailureThreshold: policy.FailureThreshold,
			SuccessThreshold: policy.SuccessThreshold,
			OpenTimeoutSec:   int(policy.OpenTimeout.Seconds()),
		}
		writeJSON(w, http.StatusOK, resp)
	}
}
