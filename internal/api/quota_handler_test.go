package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
)

func newTestQuotaHandler(max int, window time.Duration) http.HandlerFunc {
	q := lease.NewRenewalQuota(lease.QuotaPolicy{
		MaxRenewalsPerWindow: max,
		Window:               window,
	})
	return handleQuotaStatus(q)
}

func TestQuotaHandler_PolicyFields(t *testing.T) {
	h := newTestQuotaHandler(25, 30*time.Second)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/quota", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	pol := resp["policy"].(map[string]interface{})
	if int(pol["max_renewals_per_window"].(float64)) != 25 {
		t.Errorf("expected max 25")
	}
	if int(pol["window_seconds"].(float64)) != 30 {
		t.Errorf("expected window 30s")
	}
}

func TestQuotaHandler_EmptyUsage(t *testing.T) {
	h := newTestQuotaHandler(10, time.Minute)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/quota", nil))

	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp) //nolint

	usage := resp["usage"].([]interface{})
	if len(usage) != 0 {
		t.Errorf("expected empty usage, got %d entries", len(usage))
	}
}

func TestQuotaHandler_UsageAfterAllow(t *testing.T) {
	q := lease.NewRenewalQuota(lease.QuotaPolicy{MaxRenewalsPerWindow: 10, Window: time.Minute})
	q.Allow("ns/db")  //nolint
	q.Allow("ns/db")  //nolint
	q.Allow("ns/kv")  //nolint

	h := handleQuotaStatus(q)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/quota", nil))

	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp) //nolint

	usage := resp["usage"].([]interface{})
	if len(usage) != 2 {
		t.Errorf("expected 2 usage entries, got %d", len(usage))
	}
}

func TestQuotaHandler_MethodNotAllowed(t *testing.T) {
	h := newTestQuotaHandler(10, time.Minute)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/quota", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestQuotaHandler_ContentType(t *testing.T) {
	h := newTestQuotaHandler(10, time.Minute)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/quota", nil))

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}
