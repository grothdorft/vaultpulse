package alert_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/alert"
)

func TestWebhook_Send_Success(t *testing.T) {
	var received alert.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	wh := alert.NewWebhook(server.URL)
	p := alert.Payload{
		LeaseID:   "lease/abc123",
		Path:      "secret/data/myapp",
		ExpiresAt: time.Now().Add(10 * time.Minute),
		TTL:       600,
		Message:   "Lease expiring soon",
	}

	if err := wh.Send(context.Background(), p); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received.LeaseID != p.LeaseID {
		t.Errorf("expected lease_id %q, got %q", p.LeaseID, received.LeaseID)
	}
}

func TestWebhook_Send_Non2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	wh := alert.NewWebhook(server.URL)
	err := wh.Send(context.Background(), alert.Payload{LeaseID: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}

func TestWebhook_Send_InvalidURL(t *testing.T) {
	wh := alert.NewWebhook("http://127.0.0.1:0/nonexistent")
	err := wh.Send(context.Background(), alert.Payload{LeaseID: "test"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
