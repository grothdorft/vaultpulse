package alert

import (
	"context"
	"sync"
)

// MockAlerter is a test double for the Alerter interface.
type MockAlerter struct {
	mu       sync.Mutex
	Calls    []string
	Err      error
	CallCount int
}

// NewMockAlerter returns a new MockAlerter with no pre-configured error.
func NewMockAlerter() *MockAlerter {
	return &MockAlerter{}
}

// Send records the leaseID and returns the pre-configured error.
func (m *MockAlerter) Send(ctx context.Context, leaseID string, ttl string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, leaseID)
	m.CallCount++
	return m.Err
}

// Reset clears recorded calls and resets the call count.
func (m *MockAlerter) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = nil
	m.CallCount = 0
	m.Err = nil
}
