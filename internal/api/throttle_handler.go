package api

import (
	"net/http"

	"github.com/yourusername/vaultpulse/internal/lease"
)

// throttleStatusResponse is the JSON shape returned by the throttle status endpoint.
type throttleStatusResponse struct {
	MaxPerInterval int    `json:"max_per_interval"`
	IntervalSec    int    `json:"interval_sec"`
	Remaining      int    `json:"remaining"`
	Policy         string `json:"policy"`
}

// handleThrottleStatus returns the current renewal throttle configuration and
// remaining capacity for the rolling window.
func handleThrottleStatus(throttle *lease.RenewalThrottle, policy lease.ThrottlePolicy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := throttleStatusResponse{
			MaxPerInterval: policy.MaxPerInterval,
			IntervalSec:    int(policy.Interval.Seconds()),
			Remaining:      throttle.Remaining(),
			Policy:         "sliding_window",
		}

		writeJSON(w, http.StatusOK, resp)
	}
}
