package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type leaseEntry struct {
	LeaseID   string    `json:"lease_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Alerted   bool      `json:"alerted"`
	Expired   bool      `json:"expired"`
}

type leasesResponse struct {
	Count  int          `json:"count"`
	Leases []leaseEntry `json:"leases"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleLeases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	all := s.tracker.All()
	entries := make([]leaseEntry, 0, len(all))
	for _, e := range all {
		entries = append(entries, leaseEntry{
			LeaseID:   e.LeaseID,
			ExpiresAt: e.ExpiresAt,
			Alerted:   e.Alerted,
			Expired:   e.Expired,
		})
	}
	writeJSON(w, http.StatusOK, leasesResponse{
		Count:  len(entries),
		Leases: entries,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
