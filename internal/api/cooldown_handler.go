package api

import (
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

type cooldownStatusResponse struct {
	Policy    cooldownPolicyResponse    `json:"policy"`
	Cooldowns []leaseCooldownEntry      `json:"cooldowns"`
}

type cooldownPolicyResponse struct {
	MinIntervalSeconds float64 `json:"min_interval_seconds"`
}

type leaseCooldownEntry struct {
	LeaseID            string  `json:"lease_id"`
	RemainingSeconds   float64 `json:"remaining_seconds"`
	InCooldown         bool    `json:"in_cooldown"`
}

func handleCooldownStatus(cd *lease.RenewalCooldown, policy lease.CooldownPolicy, leaseIDs []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		entries := make([]leaseCooldownEntry, 0, len(leaseIDs))
		for _, id := range leaseIDs {
			remaining := cd.Remaining(id)
			entries = append(entries, leaseCooldownEntry{
				LeaseID:          id,
				RemainingSeconds: remaining.Seconds(),
				InCooldown:       remaining > 0,
			})
		}

		resp := cooldownStatusResponse{
			Policy: cooldownPolicyResponse{
				MinIntervalSeconds: policy.MinInterval.Seconds(),
			},
			Cooldowns: entries,
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

// buildCooldownLeaseIDs is a helper that extracts lease IDs from a slice of
// tracker entries for use with handleCooldownStatus.
func buildCooldownLeaseIDs(entries []*lease.Entry) []string {
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		if e != nil {
			ids = append(ids, e.LeaseID)
		}
	}
	return ids
}

// cooldownRemainingSeconds is a convenience wrapper used in tests.
func cooldownRemainingSeconds(cd *lease.RenewalCooldown, leaseID string) float64 {
	return cd.Remaining(leaseID).Seconds()
}

// ensure time is imported (used transitively via lease.RenewalCooldown).
var _ = time.Second
