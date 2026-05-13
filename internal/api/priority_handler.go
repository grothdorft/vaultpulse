package api

import (
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

type priorityEntry struct {
	LeaseID   string `json:"lease_id"`
	Priority  string `json:"priority"`
	ExpiresAt string `json:"expires_at"`
	Remaining string `json:"remaining"`
}

// handleLeasePriority returns all tracked leases ranked by renewal urgency.
func (s *Server) handleLeasePriority(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	warnDur := 30 * time.Minute
	criticalDur := 10 * time.Minute
	if s.cfg != nil {
		warnDur = s.cfg.WarnThreshold
		criticalDur = s.cfg.CriticalThreshold
	}

	all := s.tracker.All()
	entries := make([]*lease.Entry, 0, len(all))
	for _, e := range all {
		entries = append(entries, e)
	}

	ranked := lease.RankEntries(entries, warnDur, criticalDur)

	result := make([]priorityEntry, 0, len(ranked))
	for _, pe := range ranked {
		result = append(result, priorityEntry{
			LeaseID:   pe.Entry.LeaseID,
			Priority:  pe.Priority.String(),
			ExpiresAt: pe.Entry.ExpiresAt.UTC().Format(time.RFC3339),
			Remaining: pe.Remaining.Round(time.Second).String(),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(result),
		"leases":  result,
	})
}
