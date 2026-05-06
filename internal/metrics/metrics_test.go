package metrics

import (
	"testing"
)

func TestCounters_IncrementAndSnapshot(t *testing.T) {
	c := &Counters{}

	c.AlertsTotal.Add(3)
	c.AlertErrors.Add(1)
	c.LeasesTracked.Add(10)
	c.LeasesExpired.Add(2)
	c.ChecksRun.Add(5)

	s := c.Snapshot()

	if s.AlertsTotal != 3 {
		t.Errorf("expected AlertsTotal=3, got %d", s.AlertsTotal)
	}
	if s.AlertErrors != 1 {
		t.Errorf("expected AlertErrors=1, got %d", s.AlertErrors)
	}
	if s.LeasesTracked != 10 {
		t.Errorf("expected LeasesTracked=10, got %d", s.LeasesTracked)
	}
	if s.LeasesExpired != 2 {
		t.Errorf("expected LeasesExpired=2, got %d", s.LeasesExpired)
	}
	if s.ChecksRun != 5 {
		t.Errorf("expected ChecksRun=5, got %d", s.ChecksRun)
	}
}

func TestCounters_Reset(t *testing.T) {
	c := &Counters{}
	c.AlertsTotal.Add(7)
	c.ChecksRun.Add(4)

	c.Reset()
	s := c.Snapshot()

	if s.AlertsTotal != 0 {
		t.Errorf("expected AlertsTotal=0 after reset, got %d", s.AlertsTotal)
	}
	if s.ChecksRun != 0 {
		t.Errorf("expected ChecksRun=0 after reset, got %d", s.ChecksRun)
	}
}

func TestGlobal_IsNonNil(t *testing.T) {
	if Global == nil {
		t.Fatal("expected Global metrics to be non-nil")
	}
}
