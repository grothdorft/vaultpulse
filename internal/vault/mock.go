package vault

import (
	"context"
)

// MockClient is a test double for the VaultClient interface.
type MockClient struct {
	Leases []Lease
	Err    error
}

// NewMockClient returns a MockClient pre-loaded with the given leases.
func NewMockClient(leases []Lease) *MockClient {
	if leases == nil {
		leases = []Lease{}
	}
	return &MockClient{Leases: leases}
}

// ListLeases returns the pre-configured leases or error.
func (m *MockClient) ListLeases(ctx context.Context) ([]Lease, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Leases, nil
}
