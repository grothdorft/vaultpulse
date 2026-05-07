package api

import (
	"context"
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/lease"
	"github.com/yourusername/vaultpulse/internal/metrics"
	"go.uber.org/zap"
)

// Server exposes an HTTP API for VaultPulse status and lease info.
type Server struct {
	addr    string
	tracker *lease.Tracker
	logger  *zap.Logger
	httpSrv *http.Server
}

// New creates a new API Server.
func New(addr string, tracker *lease.Tracker, logger *zap.Logger) *Server {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	s := &Server{
		addr:    addr,
		tracker: tracker,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/leases", s.handleLeases)
	mux.Handle("/metrics", metrics.Handler())

	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return s
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	s.logger.Info("api server starting", zap.String("addr", s.addr))
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("api server shutting down")
	return s.httpSrv.Shutdown(ctx)
}
