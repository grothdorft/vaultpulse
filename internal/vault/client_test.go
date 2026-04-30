package vault

import (
	"testing"
	"time"
)

func TestLease_IsExpiringSoon(t *testing.T) {
	tests := []struct {
		name      string
		expiresIn time.Duration
		threshold time.Duration
		want      bool
	}{
		{
			name:      "expires within threshold",
			expiresIn: 10 * time.Minute,
			threshold: 30 * time.Minute,
			want:      true,
		},
		{
			name:      "expires exactly at threshold",
			expiresIn: 30 * time.Minute,
			threshold: 30 * time.Minute,
			want:      true,
		},
		{
			name:      "expires well beyond threshold",
			expiresIn: 2 * time.Hour,
			threshold: 30 * time.Minute,
			want:      false,
		},
		{
			name:      "already expired",
			expiresIn: -1 * time.Minute,
			threshold: 30 * time.Minute,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lease := &Lease{
				LeaseID:   "test/lease/id",
				Path:      "secret/data/myapp",
				TTL:       tt.expiresIn,
				ExpiresAt: time.Now().Add(tt.expiresIn),
			}
			got := lease.IsExpiringSoon(tt.threshold)
			if got != tt.want {
				t.Errorf("IsExpiringSoon(%v) = %v, want %v (expiresIn=%v)",
					tt.threshold, got, tt.want, tt.expiresIn)
			}
		})
	}
}

func TestNewClient_InvalidAddress(t *testing.T) {
	// NewClient itself does not dial; invalid address surfaces on first request.
	// We just verify it returns a non-nil client without error.
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("NewClient returned nil client")
	}
}
