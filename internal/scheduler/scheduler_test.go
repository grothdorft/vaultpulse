package scheduler_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/alert"
	"github.com/yourusername/vaultpulse/internal/lease"
	"github.com/yourusername/vaultpulse/internal/monitor"
	"github.com/yourusername/vaultpulse/internal/scheduler"
	"github.com/yourusername/vaultpulse/internal/vault"
)

func newTestScheduler(t *testing.T, interval time.Duration) *scheduler.Scheduler {
	t.Helper()
	mockVault := vault.NewMockClient(nil)
	mockAlert := alert.NewMockAlerter()
	tracker := lease.NewTracker()
	logger := log.New(os.Stderr, "[test] ", 0)
	m := monitor.New(mockVault, mockAlert, tracker, 5*time.Minute, logger)
	return scheduler.New(m, interval, logger)
}

func TestScheduler_RunsImmediateCheck(t *testing.T) {
	sched := newTestScheduler(t, 10*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Should return due to context cancellation, not an error.
	err := sched.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("expected context error, got: %v", err)
	}
}

func TestScheduler_TicksRepeatedly(t *testing.T) {
	sched := newTestScheduler(t, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_ = sched.Run(ctx)
	elapsed := time.Since(start)

	// We expect at least ~200ms to have passed (context timeout).
	if elapsed < 150*time.Millisecond {
		t.Errorf("expected scheduler to run for ~200ms, ran for %s", elapsed)
	}
}

func TestScheduler_NilLoggerUsesDefault(t *testing.T) {
	mockVault := vault.NewMockClient(nil)
	mockAlert := alert.NewMockAlerter()
	tracker := lease.NewTracker()
	m := monitor.New(mockVault, mockAlert, tracker, 5*time.Minute, nil)

	// Should not panic with nil logger.
	sched := scheduler.New(m, 100*time.Millisecond, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := sched.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("unexpected error: %v", err)
	}
}
