package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

type windowEntry struct {
	LeaseID       string  `json:"lease_id"`
	WindowSeconds float64 `json:"window_seconds"`
}

type windowPolicyResponse struct {
	MinWindowSeconds  float64       `json:"min_window_seconds"`
	MaxWindowSeconds  float64       `json:"max_window_seconds"`
	DefaultFraction   float64       `json:"default_fraction"`
	ActiveWindows     []windowEntry `json:"active_windows"`
}

func handleWindowStatus(rw *lease.RenewalWindow, policy lease.WindowPolicy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := rw.Snapshot()
		entries := make([]windowEntry, 0, len(snap))
		for id, dur := range snap {
			entries = append(entries, windowEntry{
				LeaseID:       id,
				WindowSeconds: dur.Seconds(),
			})
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].LeaseID < entries[j].LeaseID
		})

		resp := windowPolicyResponse{
			MinWindowSeconds: policy.MinWindow.Seconds(),
			MaxWindowSeconds: policy.MaxWindow.Seconds(),
			DefaultFraction:  policy.DefaultFraction,
			ActiveWindows:    entries,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// windowSecondsFor is a helper used in tests.
func windowSecondsFor(d time.Duration) float64 {
	return d.Seconds()
}
