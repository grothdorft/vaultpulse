package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/your-org/vaultpulse/internal/lease"
	"github.com/your-org/vaultpulse/internal/metrics"
)

// Server wraps an HTTP server with all VaultPulse API routes registered.
type Server struct {
	httpServer *http.Server
	logger     *log.Logger
}

// New creates and configures the API server.
func New(
	addr string,
	tracker *lease.Tracker,
	counters *metrics.Counters,
	stats *lease.RenewalStats,
	throttle *lease.RenewalThrottle,
	cb *lease.RenewalCircuitBreaker,
	cbPolicy lease.CircuitBreakerPolicy,
	logger *log.Logger,
) *Server {
	if logger == nil {
		logger = log.Default()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handleHealth())
	mux.HandleFunc("/leases", handleLeases(tracker))
	mux.HandleFunc("/status", handleStatus(tracker))
	mux.HandleFunc("/metrics", metrics.Handler(counters))
	mux.HandleFunc("/renewal/stats", handleRenewalStats(stats))
	mux.HandleFunc("/renewal/throttle", handleThrottleStatus(throttle))
	mux.HandleFunc("/renewal/circuit-breaker", handleCircuitBreakerStatus(cb, cbPolicy))

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// Start begins listening and blocks until the server stops.
func (s *Server) Start() error {
	s.logger.Printf("api: listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
