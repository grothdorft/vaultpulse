package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

func TestHandleHealth_ResponseBody(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	s.handleHealth(rr, req)

	var resp healthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %s", resp.Status)
	}
	if resp.Time == "" {
		t.Error("expected non-empty time")
	}
}

func TestHandleLeases_CountMatchesTracker(t *testing.T) {
	tracker := lease.NewTracker()
	tracker.Upsert("lease-1", time.Now().Add(5*time.Minute))
	tracker.Upsert("lease-2", time.Now().Add(10*time.Minute))
	s := New(":0", tracker, nil)

	req := httptest.NewRequest(http.MethodGet, "/leases", nil)
	rr := httptest.NewRecorder()
	s.handleLeases(rr, req)

	var resp leasesResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Count != 2 {
		t.Errorf("expected count 2, got %d", resp.Count)
	}
	if len(resp.Leases) != 2 {
		t.Errorf("expected 2 leases, got %d", len(resp.Leases))
	}
}

func TestHandleLeases_MethodNotAllowed(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodDelete, "/leases", nil)
	rr := httptest.NewRecorder()
	s.handleLeases(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestHandleLeases_LeaseFields(t *testing.T) {
	tracker := lease.NewTracker()
	tracker.Upsert("lease-xyz", time.Now().Add(3*time.Minute))
	s := New(":0", tracker, nil)

	req := httptest.NewRequest(http.MethodGet, "/leases", nil)
	rr := httptest.NewRecorder()
	s.handleLeases(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "lease-xyz") {
		t.Errorf("expected lease-xyz in response body, got: %s", body)
	}
}
