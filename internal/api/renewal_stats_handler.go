package api

import (
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

// renewalStatsResponse is the JSON shape returned by the renewal stats endpoint.
type renewalStatsResponse struct {
	TotalAttempts  int    `json:"total_attempts"`
	Successes      int    `json:"successes"`
	Failures       int    `json:"failures"`
	LastSuccessAt  string `json:"last_success_at,omitempty"`
	LastFailureAt  string `json:"last_failure_at,omitempty"`
	LastFailureMsg string `json:"last_failure_msg,omitempty"`
}

const timeLayout = time.RFC3339

func formatOptionalTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

// handleRenewalStats serves a snapshot of renewal outcome counters.
func handleRenewalStats(stats *lease.RenewalStats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := stats.Snapshot()
		resp := renewalStatsResponse{
			TotalAttempts:  snap.TotalAttempts,
			Successes:      snap.Successes,
			Failures:       snap.Failures,
			LastSuccessAt:  formatOptionalTime(snap.LastSuccessAt),
			LastFailureAt:  formatOptionalTime(snap.LastFailureAt),
			LastFailureMsg: snap.LastFailureMsg,
		}

		writeJSON(w, http.StatusOK, resp)
	}
}
