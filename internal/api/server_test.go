package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	tracker := lease.NewTracker()
	return New(":0", tracker, nil)
}

func TestServer_HealthEndpoint(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	s.httpSrv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestServer_HealthMethodNotAllowed(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rr := httptest.NewRecorder()
	s.httpSrv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestServer_LeasesEndpointEmpty(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/leases", nil)
	rr := httptest.NewRecorder()
	s.httpSrv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestServer_LeasesEndpointWithData(t *testing.T) {
	tracker := lease.NewTracker()
	tracker.Upsert("lease-abc", time.Now().Add(10*time.Minute))
	s := New(":0", tracker, nil)

	req := httptest.NewRequest(http.MethodGet, "/leases", nil)
	rr := httptest.NewRecorder()
	s.httpSrv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestServer_Shutdown(t *testing.T) {
	s := newTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}
