package api

import (
	"net/http"

	"github.com/yourusername/vaultpulse/internal/lease"
)

// handleAuditLog serves GET /audit — returns all renewal audit events.
// Optional query param ?lease_id=<id> filters by lease.
func handleAuditLog(log *lease.RenewalAuditLog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		leaseID := r.URL.Query().Get("lease_id")

		var events interface{}
		if leaseID != "" {
			events = log.ForLease(leaseID)
		} else {
			events = log.All()
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"total":  log.Len(),
			"events": events,
		})
	}
}
