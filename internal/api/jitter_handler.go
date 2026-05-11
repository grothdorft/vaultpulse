package api

import (
	"net/http"

	"github.com/yourusername/vaultpulse/internal/lease"
)

// handleJitterPolicy serves the current jitter policy configuration as JSON.
func handleJitterPolicy(policy lease.JitterPolicy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type response struct {
			MaxFraction float64 `json:"max_fraction"`
			Description string  `json:"description"`
		}

		body := response{
			MaxFraction: policy.MaxFraction,
			Description: "Fraction of base renewal delay added as random jitter to spread load",
		}

		writeJSON(w, http.StatusOK, body)
	}
}
