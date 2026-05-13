package api

import (
	"net/http"

	"github.com/yourusername/vaultpulse/internal/lease"
)

// RegisterWindowRoute attaches the renewal window status endpoint to the given mux.
// GET /api/v1/renewal/window returns the active per-lease computed windows
// alongside the policy configuration.
func RegisterWindowRoute(mux *http.ServeMux, rw *lease.RenewalWindow, policy lease.WindowPolicy) {
	mux.HandleFunc("/api/v1/renewal/window", handleWindowStatus(rw, policy))
}
