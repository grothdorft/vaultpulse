package vault

// MockClient is a test double for vault.Client that returns predefined leases.
type MockClient struct {
	leases []Lease
}

// NewMockClient creates a MockClient with the provided leases.
func NewMockClient(leases []Lease) *Client {
	return &Client{mock: &MockClient{leases: leases}}
}
