package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	c := &Counters{}
	c.AlertsTotal.Add(2)
	c.ChecksRun.Add(8)

	h := Handler(c, nil)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var s Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if s.AlertsTotal != 2 {
		t.Errorf("expected AlertsTotal=2, got %d", s.AlertsTotal)
	}
	if s.ChecksRun != 8 {
		t.Errorf("expected ChecksRun=8, got %d", s.ChecksRun)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	c := &Counters{}
	h := Handler(c, nil)

	req := httptest.NewRequest(http.MethodPost, "/metrics", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_ContentTypeJSON(t *testing.T) {
	c := &Counters{}
	h := Handler(c, nil)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
