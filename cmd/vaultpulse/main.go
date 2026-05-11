// Package main is the entry point for the vaultpulse CLI tool.
// It wires together configuration, Vault client, lease tracker, monitor,
// scheduler, alerting, and the HTTP API server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/vaultpulse/internal/alert"
	"github.com/yourusername/vaultpulse/internal/api"
	"github.com/yourusername/vaultpulse/internal/config"
	"github.com/yourusername/vaultpulse/internal/lease"
	"github.com/yourusername/vaultpulse/internal/monitor"
	"github.com/yourusername/vaultpulse/internal/scheduler"
	"github.com/yourusername/vaultpulse/internal/vault"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if err := run(logger); err != nil {
		logger.Error("fatal error", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	// Load configuration from the default path or VAULTPULSE_CONFIG env var.
	cfgPath := os.Getenv("VAULTPULSE_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/vaultpulse.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger.Info("configuration loaded", "vault_address", cfg.VaultAddress, "interval", cfg.CheckInterval)

	// Build the Vault client.
	vaultClient, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}

	// Build the lease tracker.
	tracker := lease.NewTracker()

	// Build the alerter.
	alerter := alert.NewWebhook(cfg.WebhookURL, logger)

	// Build the monitor that ties Vault, tracker, and alerter together.
	mon := monitor.New(vaultClient, tracker, alerter, cfg, logger)

	// Build the HTTP API server.
	apiServer := api.New(tracker, logger)

	// Build the scheduler that drives periodic checks.
	sched := scheduler.New(mon, cfg.CheckInterval, logger)

	// Root context — cancelled on OS signal.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start the HTTP API server in a goroutine.
	httpAddr := cfg.ListenAddress
	if httpAddr == "" {
		httpAddr = ":8080"
	}
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: apiServer,
	}

	go func() {
		logger.Info("API server listening", "addr", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("API server error", "err", err)
		}
	}()

	// Start the scheduler (blocks until context is cancelled).
	sched.Run(ctx)

	// Graceful shutdown of the HTTP server.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Warn("HTTP server shutdown error", "err", err)
	}

	logger.Info("vaultpulse stopped")
	return nil
}
