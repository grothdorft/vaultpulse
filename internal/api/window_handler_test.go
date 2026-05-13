package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

func newTestWindowHandler() (*lease.RenewalWindow, http.HandlerFunc) {
	policy := lease.DefaultWindowPolicy()
	rw := lease.NewRenewalWindow(policy)
	return rw, handleWindowStatus(rw, policy)
}

func TestWindowHandler_EmptySnapshot(t *testing.T) {
	_, h := newTestWindowHandler()
	req := httptest.NewRequest(http.MethodGet, "/window", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp windowPolicyResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp.ActiveWindows) != 0 {
		t.Errorf("expected no active windows, got %d", len(resp.ActiveWindows))
	}
}

func TestWindowHandler_PolicyFields(t *testing.T) {
	_, h := newTestWindowHandler()
	req := httptest.NewRequest(http.MethodGet, "/window", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	var resp windowPolicyResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.MinWindowSeconds != (30 * time.Second).Seconds() {
		t.Errorf("unexpected MinWindowSeconds: %v", resp.MinWindowSeconds)
	}
	if resp.MaxWindowSeconds != (24 * time.Hour).Seconds() {
		t.Errorf("unexpected MaxWindowSeconds: %v", resp.MaxWindowSeconds)
	}
	if resp.DefaultFraction != 0.25 {
		t.Errorf("unexpected DefaultFraction: %v", resp.DefaultFraction)
	}
}

func TestWindowHandler_ActiveWindows(t *testing.T) {
	rw, h := newTestWindowHandler()
	rw.Compute("lease-a", 4*time.Hour)
	rw.Compute("lease-b", 8*time.Hour)
	req := httptest.NewRequest(http.MethodGet, "/window", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	var resp windowPolicyResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.ActiveWindows) != 2 {
		t.Fatalf("expected 2 active windows, got %d", len(resp.ActiveWindows))
	}
	// sorted by lease_id
	if resp.ActiveWindows[0].LeaseID != "lease-a" {
		t.Errorf("expected lease-a first, got %s", resp.ActiveWindows[0].LeaseID)
	}
	if resp.ActiveWindows[0].WindowSeconds != time.Hour.Seconds() {
		t.Errorf("expected 3600s for lease-a, got %v", resp.ActiveWindows[0].WindowSeconds)
	}
	if resp.ActiveWindows[1].WindowSeconds != (2 * time.Hour).Seconds() {
		t.Errorf("expected 7200s for lease-b, got %v", resp.ActiveWindows[1].WindowSeconds)
	}
}

func TestWindowHandler_MethodNotAllowed(t *testing.T) {
	_, h := newTestWindowHandler()
	req := httptest.NewRequest(http.MethodPost, "/window", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestWindowHandler_ContentType(t *testing.T) {
	_, h := newTestWindowHandler()
	req := httptest.NewRequest(http.MethodGet, "/window", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}
