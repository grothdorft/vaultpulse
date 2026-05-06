package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/vaultpulse/internal/monitor"
)

// Scheduler periodically triggers the monitor to check leases.
type Scheduler struct {
	monitor  *monitor.Monitor
	interval time.Duration
	logger   *log.Logger
}

// New creates a new Scheduler with the given monitor and polling interval.
func New(m *monitor.Monitor, interval time.Duration, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		monitor:  m,
		interval: interval,
		logger:   logger,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	s.logger.Printf("scheduler: starting with interval %s", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run an immediate check before waiting for the first tick.
	if err := s.runCheck(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := s.runCheck(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			s.logger.Println("scheduler: context cancelled, stopping")
			return ctx.Err()
		}
	}
}

func (s *Scheduler) runCheck(ctx context.Context) error {
	s.logger.Println("scheduler: running lease check")
	if err := s.monitor.Check(ctx); err != nil {
		s.logger.Printf("scheduler: check error: %v", err)
	}
	return nil
}
