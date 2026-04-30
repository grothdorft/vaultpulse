package vault

import (
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Lease represents a Vault secret lease with its metadata.
type Lease struct {
	LeaseID   string
	Path      string
	TTL       time.Duration
	ExpiresAt time.Time
}

// Client wraps the Vault API client with lease monitoring capabilities.
type Client struct {
	api *vaultapi.Client
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	api.SetToken(token)

	return &Client{api: api}, nil
}

// LookupLease retrieves metadata for a specific lease ID.
func (c *Client) LookupLease(leaseID string) (*Lease, error) {
	secret, err := c.api.Sys().Lookup(leaseID)
	if err != nil {
		return nil, fmt.Errorf("looking up lease %q: %w", leaseID, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("lease %q not found", leaseID)
	}

	ttlRaw, ok := secret.Data["ttl"]
	if !ok {
		return nil, fmt.Errorf("lease %q missing ttl field", leaseID)
	}

	ttlFloat, ok := ttlRaw.(float64)
	if !ok {
		return nil, fmt.Errorf("lease %q ttl has unexpected type %T", leaseID, ttlRaw)
	}

	ttl := time.Duration(ttlFloat) * time.Second

	pathRaw, _ := secret.Data["path"]
	path, _ := pathRaw.(string)

	return &Lease{
		LeaseID:   leaseID,
		Path:      path,
		TTL:       ttl,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

// IsExpiringSoon returns true if the lease expires within the given threshold.
func (l *Lease) IsExpiringSoon(threshold time.Duration) bool {
	return time.Until(l.ExpiresAt) <= threshold
}
