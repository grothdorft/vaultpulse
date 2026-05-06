package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the alert payload sent to a webhook.
type Payload struct {
	LeaseID   string    `json:"lease_id"`
	Path      string    `json:"path"`
	ExpiresAt time.Time `json:"expires_at"`
	TTL       int       `json:"ttl_seconds"`
	Message   string    `json:"message"`
}

// Webhook sends alert notifications to a configured HTTP endpoint.
type Webhook struct {
	URL    string
	client *http.Client
}

// NewWebhook creates a new Webhook alerter with the given URL.
func NewWebhook(url string) *Webhook {
	return &Webhook{
		URL: url,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send marshals the payload and posts it to the webhook URL.
func (w *Webhook) Send(ctx context.Context, p Payload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("alert: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("alert: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("alert: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("alert: webhook returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}
