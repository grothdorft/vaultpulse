package lease

import (
	"sort"
	"time"
)

// PriorityLevel ranks how urgently a lease needs renewal.
type PriorityLevel int

const (
	PriorityLow      PriorityLevel = 0
	PriorityMedium   PriorityLevel = 1
	PriorityHigh     PriorityLevel = 2
	PriorityCritical PriorityLevel = 3
)

// String returns a human-readable label for the priority level.
func (p PriorityLevel) String() string {
	switch p {
	case PriorityCritical:
		return "critical"
	case PriorityHigh:
		return "high"
	case PriorityMedium:
		return "medium"
	default:
		return "low"
	}
}

// PrioritizedEntry pairs a tracker entry with its computed priority.
type PrioritizedEntry struct {
	Entry    *Entry
	Priority PriorityLevel
	Remaining time.Duration
}

// AssignPriority returns the PriorityLevel for a lease entry based on
// how much time remains before expiry relative to the warning thresholds.
func AssignPriority(e *Entry, warn, critical time.Duration) PriorityLevel {
	if e == nil {
		return PriorityLow
	}
	remaining := time.Until(e.ExpiresAt)
	switch {
	case remaining <= 0:
		return PriorityCritical
	case remaining <= critical:
		return PriorityCritical
	case remaining <= warn:
		return PriorityHigh
	case remaining <= warn*2:
		return PriorityMedium
	default:
		return PriorityLow
	}
}

// RankEntries returns entries sorted from highest to lowest priority.
// Entries with equal priority are sorted by soonest expiry first.
func RankEntries(entries []*Entry, warn, critical time.Duration) []PrioritizedEntry {
	result := make([]PrioritizedEntry, 0, len(entries))
	for _, e := range entries {
		if e == nil {
			continue
		}
		result = append(result, PrioritizedEntry{
			Entry:     e,
			Priority:  AssignPriority(e, warn, critical),
			Remaining: time.Until(e.ExpiresAt),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Priority != result[j].Priority {
			return result[i].Priority > result[j].Priority
		}
		return result[i].Remaining < result[j].Remaining
	})
	return result
}
