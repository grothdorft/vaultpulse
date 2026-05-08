package lease

import (
	"testing"
	"time"
)

func TestRenewalStats_InitialState(t *testing.T) {
	var s RenewalStats
	snap := s.Snapshot()
	if snap.TotalAttempts != 0 || snap.Successes != 0 || snap.Failures != 0 {
		t.Fatal("expected zero initial state")
	}
}

func TestRenewalStats_RecordSuccess(t *testing.T) {
	var s RenewalStats
	before := time.Now()
	s.RecordSuccess()
	snap := s.Snapshot()

	if snap.TotalAttempts != 1 {
		t.Errorf("expected TotalAttempts=1, got %d", snap.TotalAttempts)
	}
	if snap.Successes != 1 {
		t.Errorf("expected Successes=1, got %d", snap.Successes)
	}
	if snap.Failures != 0 {
		t.Errorf("expected Failures=0, got %d", snap.Failures)
	}
	if snap.LastSuccessAt.Before(before) {
		t.Error("LastSuccessAt should be after test start")
	}
}

func TestRenewalStats_RecordFailure(t *testing.T) {
	var s RenewalStats
	before := time.Now()
	s.RecordFailure("permission denied")
	snap := s.Snapshot()

	if snap.TotalAttempts != 1 {
		t.Errorf("expected TotalAttempts=1, got %d", snap.TotalAttempts)
	}
	if snap.Failures != 1 {
		t.Errorf("expected Failures=1, got %d", snap.Failures)
	}
	if snap.Successes != 0 {
		t.Errorf("expected Successes=0, got %d", snap.Successes)
	}
	if snap.LastFailureMsg != "permission denied" {
		t.Errorf("unexpected failure msg: %s", snap.LastFailureMsg)
	}
	if snap.LastFailureAt.Before(before) {
		t.Error("LastFailureAt should be after test start")
	}
}

func TestRenewalStats_Mixed(t *testing.T) {
	var s RenewalStats
	s.RecordSuccess()
	s.RecordSuccess()
	s.RecordFailure("timeout")
	snap := s.Snapshot()

	if snap.TotalAttempts != 3 {
		t.Errorf("expected TotalAttempts=3, got %d", snap.TotalAttempts)
	}
	if snap.Successes != 2 {
		t.Errorf("expected Successes=2, got %d", snap.Successes)
	}
	if snap.Failures != 1 {
		t.Errorf("expected Failures=1, got %d", snap.Failures)
	}
}

func TestRenewalStats_Reset(t *testing.T) {
	var s RenewalStats
	s.RecordSuccess()
	s.RecordFailure("err")
	s.Reset()
	snap := s.Snapshot()

	if snap.TotalAttempts != 0 || snap.Successes != 0 || snap.Failures != 0 {
		t.Fatal("expected zero state after Reset")
	}
	if snap.LastFailureMsg != "" {
		t.Error("expected empty LastFailureMsg after Reset")
	}
}

func TestRenewalStats_SnapshotIsACopy(t *testing.T) {
	var s RenewalStats
	s.RecordSuccess()
	snap := s.Snapshot()
	s.RecordSuccess()

	// snap should still reflect the state at the time of the call
	if snap.TotalAttempts != 1 {
		t.Errorf("snapshot should not reflect later mutations, got %d", snap.TotalAttempts)
	}
}
