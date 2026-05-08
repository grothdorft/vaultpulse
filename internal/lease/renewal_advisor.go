package lease

import (
	"fmt"
	"time"
)

// RenewalAdvice summarises what action (if any) should be taken for a lease.
type RenewalAdvice struct {
	LeaseID       string
	ShouldRenew   bool
	Urgent        bool
	Reason        string
	TimeRemaining time.Duration
}

// String returns a human-readable representation of the advice.
func (a RenewalAdvice) String() string {
	if !a.ShouldRenew {
		return fmt.Sprintf("lease %s: no action needed (%.0fs remaining)", a.LeaseID, a.TimeRemaining.Seconds())
	}
	urgency := "normal"
	if a.Urgent {
		urgency = "URGENT"
	}
	return fmt.Sprintf("lease %s: renew [%s] — %s (%.0fs remaining)",
		a.LeaseID, urgency, a.Reason, a.TimeRemaining.Seconds())
}

// AdviseBatch evaluates a slice of tracker entries and returns renewal advice
// for each entry, using the supplied policy. Expired entries are skipped.
func AdviseBatch(entries []*Entry, policy RenewalPolicy) []RenewalAdvice {
	advices := make([]RenewalAdvice, 0, len(entries))
	for _, e := range entries {
		if e == nil {
			continue
		}
		adv := AdviseRenewal(e, policy)
		if adv == nil {
			continue
		}
		remaining := time.Until(e.ExpiresAt)
		advices = append(advices, RenewalAdvice{
			LeaseID:       e.LeaseID,
			ShouldRenew:   adv.ShouldRenew,
			Urgent:        adv.Urgent,
			Reason:        adv.Reason,
			TimeRemaining: remaining,
		})
	}
	return advices
}
