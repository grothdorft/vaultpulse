package lease

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// RenewFunc is a function that renews a lease by ID and returns the new TTL.
type RenewFunc func(ctx context.Context, leaseID string) (time.Duration, error)

// RenewalExecutor attempts to renew leases advised for renewal and records
// the outcome in a RenewalHistory.
type RenewalExecutor struct {
	history *RenewalHistory
	renew  RenewFunc
	logger *slog.Logger
}

// NewRenewalExecutor creates a RenewalExecutor backed by the given history
// store and renewal function.
func NewRenewalExecutor(history *RenewalHistory, renew RenewFunc, logger *slog.Logger) *RenewalExecutor {
	if logger == nil {
		logger = slog.Default()
	}
	return &RenewalExecutor{
		history: history,
		renew:  renew,
		logger: logger,
	}
}

// ExecuteAdvice iterates over a batch of RenewalAdvice items, renews leases
// where ShouldRenew is true, and records each attempt in the history.
// It returns the number of successful renewals.
func (e *RenewalExecutor) ExecuteAdvice(ctx context.Context, advice []RenewalAdvice) int {
	successes := 0
	for _, a := range advice {
		if !a.ShouldRenew {
			continue
		}
		newTTL, err := e.renew(ctx, a.LeaseID)
		record := RenewalRecord{
			LeaseID:   a.LeaseID,
			RenewedAt: time.Now().UTC(),
		}
		if err != nil {
			record.Error = fmt.Sprintf("renewal failed: %v", err)
			e.logger.Warn("lease renewal failed",
				"lease_id", a.LeaseID,
				"error", err,
			)
		} else {
			record.NewTTL = newTTL
			successes++
			e.logger.Info("lease renewed",
				"lease_id", a.LeaseID,
				"new_ttl", newTTL,
			)
		}
		e.history.Record(record)
	}
	return successes
}
