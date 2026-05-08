package lease

import (
	"sync"
	"time"
)

// RenewalStats tracks aggregate renewal outcome statistics.
type RenewalStats struct {
	mu sync.Mutex

	TotalAttempts  int
	Successes      int
	Failures       int
	LastSuccessAt  time.Time
	LastFailureAt  time.Time
	LastFailureMsg string
}

// RecordSuccess increments the success counter and timestamps the event.
func (s *RenewalStats) RecordSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalAttempts++
	s.Successes++
	s.LastSuccessAt = time.Now()
}

// RecordFailure increments the failure counter and stores the error message.
func (s *RenewalStats) RecordFailure(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalAttempts++
	s.Failures++
	s.LastFailureAt = time.Now()
	s.LastFailureMsg = msg
}

// Snapshot returns a copy of the current stats without holding the lock.
func (s *RenewalStats) Snapshot() RenewalStats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return RenewalStats{
		TotalAttempts:  s.TotalAttempts,
		Successes:      s.Successes,
		Failures:       s.Failures,
		LastSuccessAt:  s.LastSuccessAt,
		LastFailureAt:  s.LastFailureAt,
		LastFailureMsg: s.LastFailureMsg,
	}
}

// Reset zeroes all counters and timestamps.
func (s *RenewalStats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	*s = RenewalStats{}
}
