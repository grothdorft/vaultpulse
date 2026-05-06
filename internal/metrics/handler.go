package metrics

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Handler returns an http.Handler that serves a JSON snapshot of c.
func Handler(c *Counters, logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s := c.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			logger.Error("failed to encode metrics", "error", err)
		}
	})
}
