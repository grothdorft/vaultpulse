package api

import (
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

type dryRunEntry struct {
	LeaseID     string    `json:"lease_id"`
	At          time.Time `json:"at"`
	ShouldRenew bool      `json:"should_renew"`
	Reason      string    `json:"reason"`
}

func handleDryRunLog(log *lease.DryRunLog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Method == http.MethodDelete {
			log.Clear()
			w.WriteHeader(http.StatusNoContent)
			return
		}

		records := log.All()
		entries := make([]dryRunEntry, 0, len(records))
		for _, rec := range records {
			entries = append(entries, dryRunEntry{
				LeaseID:     rec.LeaseID,
				At:          rec.At,
				ShouldRenew: rec.Advice.ShouldRenew,
				Reason:      rec.Reason,
			})
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"count":   len(entries),
			"entries": entries,
		})
	}
}
