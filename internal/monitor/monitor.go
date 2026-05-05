package monitor

import (
	"context"
	"log"
	"time"

	"github.com/user/vaultpulse/internal/vault"
)

// AlertFunc is called when a lease is expiring soon.
type AlertFunc func(lease vault.Lease) error

// Monitor periodically checks Vault leases and triggers alerts.
type Monitor struct {
	client   *vault.Client
	interval time.Duration
	threshold time.Duration
	alertFn  AlertFunc
}

// New creates a new Monitor.
func New(client *vault.Client, interval, threshold time.Duration, alertFn AlertFunc) *Monitor {
	return &Monitor{
		client:    client,
		interval:  interval,
		threshold: threshold,
		alertFn:   alertFn,
	}
}

// Run starts the monitoring loop and blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("monitor: shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := m.check(); err != nil {
				log.Printf("monitor: check error: %v", err)
			}
		}
	}
}

// check fetches leases and fires alerts for those expiring soon.
func (m *Monitor) check() error {
	leases, err := m.client.ListLeases()
	if err != nil {
		return err
	}

	for _, lease := range leases {
		if lease.IsExpiringSoon(m.threshold) {
			if err := m.alertFn(lease); err != nil {
				log.Printf("monitor: alert error for lease %s: %v", lease.ID, err)
			}
		}
	}
	return nil
}
