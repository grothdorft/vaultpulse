package lease

import (
	"testing"
	"time"
)

const (
	testWarn     = 30 * time.Minute
	testCritical = 10 * time.Minute
)

func TestAssignPriority_Low(t *testing.T) {
	e := &Entry{ExpiresAt: time.Now().Add(2 * time.Hour)}
	if got := AssignPriority(e, testWarn, testCritical); got != PriorityLow {
		t.Errorf("expected low, got %s", got)
	}
}

func TestAssignPriority_Medium(t *testing.T) {
	e := &Entry{ExpiresAt: time.Now().Add(45 * time.Minute)}
	if got := AssignPriority(e, testWarn, testCritical); got != PriorityMedium {
		t.Errorf("expected medium, got %s", got)
	}
}

func TestAssignPriority_High(t *testing.T) {
	e := &Entry{ExpiresAt: time.Now().Add(20 * time.Minute)}
	if got := AssignPriority(e, testWarn, testCritical); got != PriorityHigh {
		t.Errorf("expected high, got %s", got)
	}
}

func TestAssignPriority_Critical(t *testing.T) {
	e := &Entry{ExpiresAt: time.Now().Add(5 * time.Minute)}
	if got := AssignPriority(e, testWarn, testCritical); got != PriorityCritical {
		t.Errorf("expected critical, got %s", got)
	}
}

func TestAssignPriority_Expired(t *testing.T) {
	e := &Entry{ExpiresAt: time.Now().Add(-1 * time.Minute)}
	if got := AssignPriority(e, testWarn, testCritical); got != PriorityCritical {
		t.Errorf("expected critical for expired, got %s", got)
	}
}

func TestAssignPriority_NilEntry(t *testing.T) {
	if got := AssignPriority(nil, testWarn, testCritical); got != PriorityLow {
		t.Errorf("expected low for nil entry, got %s", got)
	}
}

func TestRankEntries_SortedByPriority(t *testing.T) {
	now := time.Now()
	entries := []*Entry{
		{LeaseID: "low", ExpiresAt: now.Add(2 * time.Hour)},
		{LeaseID: "critical", ExpiresAt: now.Add(5 * time.Minute)},
		{LeaseID: "high", ExpiresAt: now.Add(20 * time.Minute)},
		{LeaseID: "medium", ExpiresAt: now.Add(45 * time.Minute)},
	}
	ranked := RankEntries(entries, testWarn, testCritical)
	if len(ranked) != 4 {
		t.Fatalf("expected 4 results, got %d", len(ranked))
	}
	if ranked[0].Entry.LeaseID != "critical" {
		t.Errorf("expected critical first, got %s", ranked[0].Entry.LeaseID)
	}
	if ranked[3].Entry.LeaseID != "low" {
		t.Errorf("expected low last, got %s", ranked[3].Entry.LeaseID)
	}
}

func TestRankEntries_EqualPrioritySortedByExpiry(t *testing.T) {
	now := time.Now()
	entries := []*Entry{
		{LeaseID: "later", ExpiresAt: now.Add(8 * time.Minute)},
		{LeaseID: "sooner", ExpiresAt: now.Add(3 * time.Minute)},
	}
	ranked := RankEntries(entries, testWarn, testCritical)
	if ranked[0].Entry.LeaseID != "sooner" {
		t.Errorf("expected sooner expiry first, got %s", ranked[0].Entry.LeaseID)
	}
}

func TestRankEntries_NilSkipped(t *testing.T) {
	entries := []*Entry{nil, {LeaseID: "a", ExpiresAt: time.Now().Add(time.Hour)}}
	ranked := RankEntries(entries, testWarn, testCritical)
	if len(ranked) != 1 {
		t.Errorf("expected 1 result after nil skip, got %d", len(ranked))
	}
}

func TestPriorityLevel_String(t *testing.T) {
	cases := map[PriorityLevel]string{
		PriorityLow:      "low",
		PriorityMedium:   "medium",
		PriorityHigh:     "high",
		PriorityCritical: "critical",
	}
	for level, want := range cases {
		if got := level.String(); got != want {
			t.Errorf("PriorityLevel(%d).String() = %q, want %q", level, got, want)
		}
	}
}
