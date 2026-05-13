package api

import (
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

type quotaSnapshotEntry struct {
	Key       string `json:"key"`
	Count     int    `json:"count"`
	WindowEnd string `json:"window_end"`
}

type quotaResponse struct {
	Policy  quotaPolicyResponse  `json:"policy"`
	Usage   []quotaSnapshotEntry `json:"usage"`
}

type quotaPolicyResponse struct {
	MaxRenewalsPerWindow int    `json:"max_renewals_per_window"`
	WindowSeconds        int    `json:"window_seconds"`
}

func handleQuotaStatus(q *lease.RenewalQuota) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pol := q.Policy()
		snap := q.Snapshot()

		usage := make([]quotaSnapshotEntry, 0, len(snap))
		for key, entry := range snap {
			usage = append(usage, quotaSnapshotEntry{
				Key:       key,
				Count:     entry.Count,
				WindowEnd: entry.WindowEnd.Format(time.RFC3339),
			})
		}

		resp := quotaResponse{
			Policy: quotaPolicyResponse{
				MaxRenewalsPerWindow: pol.MaxRenewalsPerWindow,
				WindowSeconds:        int(pol.Window.Seconds()),
			},
			Usage: usage,
		}

		writeJSON(w, http.StatusOK, resp)
	}
}
