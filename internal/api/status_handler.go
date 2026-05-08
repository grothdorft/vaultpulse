package api

import (
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

type leaseStatusResponse struct {
	LeaseID   string           `json:"lease_id"`
	Status    lease.ExpiryStatus `json:"status"`
	ExpiresAt time.Time        `json:"expires_at"`
	TTLSecs   int64            `json:"ttl_seconds"`
}

// handleLeaseStatus returns enriched expiry status for all tracked leases.
func (s *Server) handleLeaseStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entries := s.tracker.All()

	warnThreshold := time.Duration(s.cfg.Alert.WarnThresholdHours) * time.Hour
	critThreshold := time.Duration(s.cfg.Alert.CritThresholdHours) * time.Hour

	results := make([]leaseStatusResponse, 0, len(entries))
	for _, entry := range entries {
		e := entry // capture
		summary := lease.BuildStatusSummary(&e, warnThreshold, critThreshold)
		results = append(results, leaseStatusResponse{
			LeaseID:   summary.LeaseID,
			Status:    summary.Status,
			ExpiresAt: summary.ExpiresAt,
			TTLSecs:   int64(summary.TTL.Seconds()),
		})
	}

	writeJSON(w, http.StatusOK, results)
}
