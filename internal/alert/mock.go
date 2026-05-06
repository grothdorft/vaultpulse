package alert

import (
	"context"
	"sync"
)

// Alerter is the interface implemented by all alert senders.
type Alerter interface {
	Send(ctx context.Context, p Payload) error
}

// MockAlerter records sent payloads and optionally returns a configured error.
type MockAlerter struct {
	mu       sync.Mutex
	Payloads []Payload
	Err      error
}

// Send records the payload and returns the configured error (if any).
func (m *MockAlerter) Send(_ context.Context, p Payload) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Payloads = append(m.Payloads, p)
	return m.Err
}

// Count returns the number of alerts sent.
func (m *MockAlerter) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Payloads)
}
