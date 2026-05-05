package monitor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/vaultpulse/internal/monitor"
	"github.com/user/vaultpulse/internal/vault"
)

func TestMonitor_CheckTriggersAlertForExpiringSoon(t *testing.T) {
	alerted := []vault.Lease{}
	alertFn := func(l vault.Lease) error {
		alerted = append(alerted, l)
		return nil
	}

	now := time.Now()
	leases := []vault.Lease{
		{ID: "lease-1", ExpiresAt: now.Add(5 * time.Minute)},  // expiring soon
		{ID: "lease-2", ExpiresAt: now.Add(60 * time.Minute)}, // not expiring soon
	}

	client := vault.NewMockClient(leases)
	m := monitor.New(client, time.Hour, 10*time.Minute, alertFn)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Trigger a single check
	if err := m.CheckOnce(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(alerted) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerted))
	}
	if alerted[0].ID != "lease-1" {
		t.Errorf("expected alert for lease-1, got %s", alerted[0].ID)
	}
	_ = ctx
}

func TestMonitor_AlertErrorDoesNotStopChecks(t *testing.T) {
	callCount := 0
	alertFn := func(l vault.Lease) error {
		callCount++
		return errors.New("webhook down")
	}

	now := time.Now()
	leases := []vault.Lease{
		{ID: "lease-a", ExpiresAt: now.Add(2 * time.Minute)},
		{ID: "lease-b", ExpiresAt: now.Add(3 * time.Minute)},
	}

	client := vault.NewMockClient(leases)
	m := monitor.New(client, time.Hour, 10*time.Minute, alertFn)

	if err := m.CheckOnce(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 alert calls, got %d", callCount)
	}
}
